package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/config"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"

	purchaseRepoI "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository"
	sessionRepoI "github.com/artrsyf/avito-trainee-assignment/internal/session/repository"
	transactionRepoI "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository"
	userRepoI "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"

	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository/postgres"
	sessionRepo "github.com/artrsyf/avito-trainee-assignment/internal/session/repository/redis"
	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository/postgres"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/postgres"

	uow "github.com/artrsyf/avito-trainee-assignment/pkg/uow/postgres"

	purchaseUsecase "github.com/artrsyf/avito-trainee-assignment/internal/purchase/usecase"
	sessionUsecase "github.com/artrsyf/avito-trainee-assignment/internal/session/usecase"
	transactionUsecase "github.com/artrsyf/avito-trainee-assignment/internal/transaction/usecase"
	userUsecase "github.com/artrsyf/avito-trainee-assignment/internal/user/usecase"

	purchaseDelivery "github.com/artrsyf/avito-trainee-assignment/internal/purchase/delivery/http"
	sessionDelivery "github.com/artrsyf/avito-trainee-assignment/internal/session/delivery/http"
	transactionDelivery "github.com/artrsyf/avito-trainee-assignment/internal/transaction/delivery/http"
	userDelivery "github.com/artrsyf/avito-trainee-assignment/internal/user/delivery/http"

	"github.com/artrsyf/avito-trainee-assignment/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	DB          *sql.DB
	RedisClient *redis.Client
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	pgDSN := "postgres://user:password@localhost:5432/reward_service_postgres_integration?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", pgDSN)
	if err != nil {
		panic(err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = DB.PingContext(ctx); err != nil {
		panic(err)
	}

	if _, err = RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}
}

func teardown() {
	ctx := context.Background()

	_, _ = DB.ExecContext(ctx, "DELETE FROM users")
	_ = RedisClient.FlushDB(ctx).Err()

	_ = DB.Close()
	_ = RedisClient.Close()
}

func TestAuthScenarios(t *testing.T) {
	cfg := setupTestEnvironment()

	t.Run("Success registration new user", func(t *testing.T) {
		testSuccessfulRegistration(t, cfg)
	})

	t.Run("Success login existing user", func(t *testing.T) {
		testSuccessfulLogin(t, cfg)
	})

	t.Run("Login with invalid password", func(t *testing.T) {
		testInvalidPassword(t, cfg)
	})

	t.Run("Register existing username", func(t *testing.T) {
		testDuplicateRegistration(t, cfg)
	})

	t.Run("Invalid request format", func(t *testing.T) {
		testInvalidRequestFormat(t, cfg)
	})

	t.Run("Token validation", func(t *testing.T) {
		testTokenValidation(t, cfg)
	})
}

func testSuccessfulRegistration(t *testing.T, cfg *TestConfig) {
	payload := dto.AuthRequest{
		Username: "new_user",
		Password: "securePassword123",
	}

	rr := sendAuthRequest(cfg, payload)

	assert.Equal(t, http.StatusOK, rr.Code)
	assertValidAuthResponse(t, rr)
	assertSessionCookies(t, rr)

	user, err := cfg.UserRepo.GetByUsername(cfg.Ctx, payload.Username)
	require.NoError(t, err)
	assert.Equal(t, payload.Username, user.Username)
	assert.Equal(t, cfg.UserConfig.InitCoinsBalance, user.Coins)

	session, err := cfg.SessionRepo.Check(cfg.Ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, session.UserID)
}

func testSuccessfulLogin(t *testing.T, cfg *TestConfig) {
	payload := dto.AuthRequest{
		Username: "existing_user",
		Password: "anotherSecurePass",
	}
	rr := sendAuthRequest(cfg, payload)
	require.Equal(t, http.StatusOK, rr.Code)

	rr = sendAuthRequest(cfg, payload)

	assert.Equal(t, http.StatusOK, rr.Code)
	assertValidAuthResponse(t, rr)
	assertSessionCookies(t, rr)

	user, err := cfg.UserRepo.GetByUsername(cfg.Ctx, payload.Username)
	require.NoError(t, err)

	newSession, err := cfg.SessionRepo.Check(cfg.Ctx, user.ID)
	require.NoError(t, err)
	assert.True(t, newSession.AccessExpiresAt.After(time.Now()))
}

func testInvalidPassword(t *testing.T, cfg *TestConfig) {
	payload := dto.AuthRequest{
		Username: "test_user",
		Password: "correctPassword",
	}
	rr := sendAuthRequest(cfg, payload)
	require.Equal(t, http.StatusOK, rr.Code)

	invalidPayload := dto.AuthRequest{
		Username: "test_user",
		Password: "wrongPassword",
	}
	rr = sendAuthRequest(cfg, invalidPayload)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "wrong credentials", response["errors"])
}

func testDuplicateRegistration(t *testing.T, cfg *TestConfig) {
	payload := dto.AuthRequest{
		Username: "duplicate_user",
		Password: "password123",
	}

	rr := sendAuthRequest(cfg, payload)
	require.Equal(t, http.StatusOK, rr.Code)

	var jwtResp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &jwtResp)
	firstToken := jwtResp["token"]

	rr = sendAuthRequest(cfg, payload)
	assert.Equal(t, http.StatusOK, rr.Code)

	json.Unmarshal(rr.Body.Bytes(), &jwtResp)
	secondToken := jwtResp["token"]

	assert.Equal(t, firstToken, secondToken)
}

func testInvalidRequestFormat(t *testing.T, cfg *TestConfig) {
	tests := []struct {
		name    string
		payload interface{}
		error   string
	}{
		{
			name:    "Empty username",
			payload: dto.AuthRequest{Password: "pass"},
			error:   "Username is required",
		},
		{
			name:    "Short password",
			payload: dto.AuthRequest{Username: "user", Password: "123"},
			error:   "Password is too short",
		},
		{
			name:    "Invalid JSON",
			payload: `{"username": "user", "password": 123}`,
			error:   "bad request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			switch v := tc.payload.(type) {
			case dto.AuthRequest:
				body, _ = json.Marshal(v)
			case string:
				body = []byte(v)
			}

			req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			cfg.Router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			var response map[string]string
			json.Unmarshal(rr.Body.Bytes(), &response)
			assert.Contains(t, response["errors"], tc.error)
		})
	}
}

func testTokenValidation(t *testing.T, cfg *TestConfig) {
	payload := dto.AuthRequest{
		Username: "token_user",
		Password: "tokenPass123",
	}
	rr := sendAuthRequest(cfg, payload)
	require.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)
	accessToken := response["token"]

	req, _ := http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr = httptest.NewRecorder()
	cfg.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expiredToken := "expired_token"
	req, _ = http.NewRequest("GET", "/api/info", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	rr = httptest.NewRecorder()
	cfg.Router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func sendAuthRequest(cfg *TestConfig, payload dto.AuthRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	cfg.Router.ServeHTTP(rr, req)
	return rr
}

func assertValidAuthResponse(t *testing.T, rr *httptest.ResponseRecorder) {
	var response dto.AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.Token)
}

func assertSessionCookies(t *testing.T, rr *httptest.ResponseRecorder) {
	cookies := rr.Result().Cookies()
	require.Len(t, cookies, 3)

	cookieMap := make(map[string]*http.Cookie)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie
	}

	assert.NotEmpty(t, cookieMap["access_token"].Value)
	assert.True(t, cookieMap["access_token"].Expires.After(time.Now()))

	assert.NotEmpty(t, cookieMap["refresh_token"].Value)
	assert.True(t, cookieMap["refresh_token"].Expires.After(time.Now()))

	userID, err := strconv.Atoi(cookieMap["user_id"].Value)
	require.NoError(t, err)
	assert.True(t, userID > 0)
}

type TestConfig struct {
	Router          *mux.Router
	UserRepo        userRepoI.UserRepositoryI
	SessionRepo     sessionRepoI.SessionRepositoryI
	TransactionRepo transactionRepoI.TransactionRepositoryI
	PurchaseRepo    purchaseRepoI.PurchaseRepositoryI
	UserConfig      config.UserConfig
	Ctx             context.Context
}

func setupTestEnvironment() *TestConfig {
	userRepo := userRepo.NewUserPostgresRepository(DB, logrus.New())
	sessionRepo := sessionRepo.NewSessionRedisRepository(RedisClient, logrus.New())
	purchaseRepo := purchaseRepo.NewPurchasePostgresRepository(DB, logrus.New())
	transactionRepo := transactionRepo.NewTransactionPostgresRepository(DB, logrus.New())

	uowFactory := uow.NewFactory(DB)

	cfg := config.UserConfig{
		InitCoinsBalance: 1000,
		Auth: config.AuthConfig{
			AccessTokenExpiration:  "1h",
			RefreshTokenExpiration: "24h",
		},
	}

	validator := validator.New()

	sessionUC := sessionUsecase.NewSessionUsecase(
		sessionRepo,
		userRepo,
		cfg,
		logrus.New(),
	)
	transactionUC := transactionUsecase.NewTransactionUsecase(
		transactionRepo,
		userRepo,
		uowFactory,
		logrus.New(),
	)
	purchaseUC := purchaseUsecase.NewPurchaseUsecase(
		purchaseRepo,
		userRepo,
		uowFactory,
		logrus.New(),
	)
	userUC := userUsecase.NewUserUsecase(
		purchaseRepo,
		transactionRepo,
		userRepo,
		logrus.New(),
	)

	router := mux.NewRouter()

	authHandler := sessionDelivery.NewSessionHandler(sessionUC, validator, logrus.New())
	transactionHandler := transactionDelivery.NewTransactionHandler(transactionUC, validator, logrus.New())
	purchaseHandler := purchaseDelivery.NewPurchaseHandler(purchaseUC, validator, logrus.New())
	userHandler := userDelivery.NewUserHandler(userUC, logrus.New())

	router.Handle("/api/auth",
		http.HandlerFunc(authHandler.Auth)).Methods("POST")

	router.Handle("/api/sendCoin",
		middleware.ValidateJWTToken(
			http.HandlerFunc(transactionHandler.SendCoins), logrus.New())).Methods("POST")

	router.Handle("/api/buy/{item}",
		middleware.ValidateJWTToken(
			http.HandlerFunc(purchaseHandler.BuyItem), logrus.New())).Methods("GET")

	router.Handle("/api/info",
		middleware.ValidateJWTToken(
			http.HandlerFunc(userHandler.GetInfo), logrus.New())).Methods("GET")

	return &TestConfig{
		Router:          router,
		UserRepo:        userRepo,
		SessionRepo:     sessionRepo,
		PurchaseRepo:    purchaseRepo,
		TransactionRepo: transactionRepo,
		UserConfig:      cfg,
		Ctx:             context.Background(),
	}
}
