package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/artrsyf/avito-trainee-assignment/config"
	"github.com/artrsyf/avito-trainee-assignment/middleware"

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err) /*TODO*/
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
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
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}

	redisAddr := fmt.Sprintf("%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DATABASE"))

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})

	defer func() {
		if err = postgresConnect.Close(); err != nil {
			panic(err)
		}

		if err := redisClient.Close(); err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()

	userRepo := userRepository.NewUserPostgresRepository(postgresConnect)
	sessionRepo := sessionRepository.NewSessionRedisRepository(redisClient)
	transactionRepo := transactionRepository.NewTransactionPostgresRepository(postgresConnect)
	purchaseRepo := purchaseRepository.NewPurchasePostgresRepository(postgresConnect)

	transactionUOW := uow.NewPostgresUnitOfWork(postgresConnect)
	purchaseUOW := uow.NewPostgresUnitOfWork(postgresConnect)

	sessionUC := sessionUsecase.NewSessionUsecase(
		sessionRepo,
		userRepo,
		cfg.User,
	)
	transactionUC := transactionUsecase.NewTransactionUsecase(
		transactionRepo,
		userRepo,
		transactionUOW,
	)
	purchaseUC := purchaseUsecase.NewPurchaseUsecase(
		purchaseRepo,
		userRepo,
		purchaseUOW,
	)
	userUC := userUsecase.NewUserUsecase(
		purchaseRepo,
		transactionRepo,
		userRepo,
	)

	authHandler := sessionDelivery.NewSessionHandler(sessionUC)
	transactionHandler := transactionDelivery.NewTransactionHandler(transactionUC)
	purchaseHandler := purchaseDelivery.NewPurchaseHandler(purchaseUC)
	userHandler := userDelivery.NewTransactionHandler(userUC)

	router.Handle("/api/auth",
		http.HandlerFunc(authHandler.Auth)).Methods("POST")

	router.Handle("/api/sendCoin",
		middleware.ValidateJWTToken(
			http.HandlerFunc(transactionHandler.SendCoins))).Methods("POST")

	router.Handle("/api/buy/{item}",
		middleware.ValidateJWTToken(
			http.HandlerFunc(purchaseHandler.BuyItem))).Methods("GET")

	router.Handle("/api/info",
		middleware.ValidateJWTToken(
			http.HandlerFunc(userHandler.GetInfo))).Methods("GET")

	fmt.Println("server starts on :8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), router))
}
