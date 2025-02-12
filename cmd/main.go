package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/artrsyf/avito-trainee-assignment/config"

	sessionRepository "github.com/artrsyf/avito-trainee-assignment/internal/session/repository/redis"
	userRepository "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/postgres"

	sessionUsecase "github.com/artrsyf/avito-trainee-assignment/internal/session/usecase"

	sessionDelivery "github.com/artrsyf/avito-trainee-assignment/internal/session/delivery/http"
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
	fmt.Println(cfg)

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

	redisURL := fmt.Sprintf("redis://user:@%s:%s/%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_DATABASE"),
	)
	redisConn, err := redis.DialURL(redisURL)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = postgresConnect.Close(); err != nil {
			panic(err)
		}

		if err = redisConn.Close(); err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()

	userRepo := userRepository.NewUserPostgresRepository(postgresConnect)
	sessionRepo := sessionRepository.NewSessionRedisRepository(redisConn)

	sessionUC := sessionUsecase.NewSessionUsecase(
		sessionRepo,
		userRepo,
		cfg.User,
	)

	authHandler := sessionDelivery.NewSessionHandler(sessionUC)

	router.Handle("/api/auth",
		http.HandlerFunc(authHandler.Auth)).Methods("POST")

	fmt.Println("server starts on :8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), router))
}
