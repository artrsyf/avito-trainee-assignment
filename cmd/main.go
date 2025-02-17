package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/artrsyf/avito-trainee-assignment/config"
	"github.com/artrsyf/avito-trainee-assignment/internal/middleware"

	purchaseRepository "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository/postgres"
	sessionRepository "github.com/artrsyf/avito-trainee-assignment/internal/session/repository/redis"
	transactionRepository "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository/postgres"
	userRepository "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/postgres"

	uow "github.com/artrsyf/avito-trainee-assignment/pkg/uow/postgres"

	purchaseUsecase "github.com/artrsyf/avito-trainee-assignment/internal/purchase/usecase"
	sessionUsecase "github.com/artrsyf/avito-trainee-assignment/internal/session/usecase"
	transactionUsecase "github.com/artrsyf/avito-trainee-assignment/internal/transaction/usecase"
	userUsecase "github.com/artrsyf/avito-trainee-assignment/internal/user/usecase"

	purchaseDelivery "github.com/artrsyf/avito-trainee-assignment/internal/purchase/delivery/http"
	sessionDelivery "github.com/artrsyf/avito-trainee-assignment/internal/session/delivery/http"
	transactionDelivery "github.com/artrsyf/avito-trainee-assignment/internal/transaction/delivery/http"
	userDelivery "github.com/artrsyf/avito-trainee-assignment/internal/user/delivery/http"
)

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return logger
}

func main() {
	logger := initLogger()
	validate := validator.New()

	err := godotenv.Load()
	if err != nil {
		logger.WithError(err).Fatal("Отсутствует .env файл")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Ошибка при загрузке конфиг файла")
	}

	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	postgresConnect, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		logger.WithError(err).Fatal("Ошибка при подключении к БД")
	}

	postgresConnect.SetMaxOpenConns(50)
	postgresConnect.SetMaxIdleConns(15)
	postgresConnect.SetConnMaxLifetime(10 * time.Minute)

	redisAddr := fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		logger.WithError(err).Fatal("Ошибка при касте индекса БД Redis")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	router := mux.NewRouter()

	userRepo := userRepository.NewUserPostgresRepository(postgresConnect, logger)
	sessionRepo := sessionRepository.NewSessionRedisRepository(redisClient, logger)
	transactionRepo := transactionRepository.NewTransactionPostgresRepository(postgresConnect, logger)
	purchaseRepo := purchaseRepository.NewPurchasePostgresRepository(postgresConnect, logger)

	uowFactory := uow.NewFactory(postgresConnect)

	sessionUC := sessionUsecase.NewSessionUsecase(
		sessionRepo,
		userRepo,
		cfg.User,
		logger,
	)
	transactionUC := transactionUsecase.NewTransactionUsecase(
		transactionRepo,
		userRepo,
		uowFactory,
		logger,
	)
	purchaseUC := purchaseUsecase.NewPurchaseUsecase(
		purchaseRepo,
		userRepo,
		uowFactory,
		logger,
	)
	userUC := userUsecase.NewUserUsecase(
		purchaseRepo,
		transactionRepo,
		userRepo,
		logger,
	)

	authHandler := sessionDelivery.NewSessionHandler(sessionUC, validate, logger)
	transactionHandler := transactionDelivery.NewTransactionHandler(transactionUC, validate, logger)
	purchaseHandler := purchaseDelivery.NewPurchaseHandler(purchaseUC, validate, logger)
	userHandler := userDelivery.NewUserHandler(userUC, logger)

	router.Handle("/api/auth",
		http.HandlerFunc(authHandler.Auth)).Methods("POST")

	router.Handle("/api/sendCoin",
		middleware.ValidateJWTToken(
			http.HandlerFunc(transactionHandler.SendCoins), logger)).Methods("POST")

	router.Handle("/api/buy/{item}",
		middleware.ValidateJWTToken(
			http.HandlerFunc(purchaseHandler.BuyItem), logger)).Methods("GET")

	router.Handle("/api/info",
		middleware.ValidateJWTToken(
			http.HandlerFunc(userHandler.GetInfo), logger)).Methods("GET")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: router,
	}

	go func() {
		logger.WithField("port", 8080).Info("Сервер запущен")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Ошибка запуска сервера")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Ошибка при остановке сервера")
	}

	if err := postgresConnect.Close(); err != nil {
		logger.WithError(err).Error("Ошибка при закрытии подключения к PostgreSQL")
	}

	if err := redisClient.Close(); err != nil {
		logger.WithError(err).Error("Ошибка при закрытии подключения к Redis")
	}

	logger.Info("Сервер успешно остановлен")
}
