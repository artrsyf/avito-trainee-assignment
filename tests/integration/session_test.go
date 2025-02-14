package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/config"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	sessionRepo "github.com/artrsyf/avito-trainee-assignment/internal/session/repository/redis"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/usecase"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/repository/postgres"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// Инициализация PostgreSQL
	pgDSN := "postgres://user:password@localhost:5432/reward_service_postgres_integration?sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", pgDSN)
	if err != nil {
		panic(err)
	}

	// Инициализация Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Проверка подключений
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
	// Очистка данных
	ctx := context.Background()

	// Очистка PostgreSQL
	_, _ = DB.ExecContext(ctx, "DELETE FROM users")

	// Очистка Redis
	_ = RedisClient.FlushDB(ctx).Err()

	// Закрытие подключений
	_ = DB.Close()
	_ = RedisClient.Close()
}

func TestSessionUsecase_Integration(t *testing.T) {
	userRepo := postgres.NewUserPostgresRepository(DB)
	sessionRepo := sessionRepo.NewSessionRedisRepository(RedisClient)

	cfg := config.UserConfig{
		InitCoinsBalance: 100,
		Auth: config.AuthConfig{
			AccessTokenExpiration:  "5s",
			RefreshTokenExpiration: "24h",
		},
	}

	uc := usecase.NewSessionUsecase(sessionRepo, userRepo, cfg)
	ctx := context.Background()

	t.Run("successful signup and session creation", func(t *testing.T) {
		req := &dto.AuthRequest{
			Username: "newuser",
			Password: "password123",
		}

		// Первая попытка - регистрация
		session, err := uc.LoginOrSignup(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, session.JWTAccess)
		require.NotEmpty(t, session.JWTRefresh)
		require.Equal(t, "newuser", session.Username)

		// Проверка создания пользователя
		user, err := userRepo.GetByUsername(ctx, "newuser")
		require.NoError(t, err)
		require.Equal(t, cfg.InitCoinsBalance, user.Coins)

		// Проверка сессии в Redis
		redisSession, err := sessionRepo.Check(ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, session.JWTAccess, redisSession.JWTAccess)
	})

	t.Run("successful login with existing user", func(t *testing.T) {
		req := &dto.AuthRequest{
			Username: "existinguser",
			Password: "password123",
		}

		// Создаем пользователя
		_, err := uc.LoginOrSignup(ctx, req)
		require.NoError(t, err)

		// Повторный логин
		session, err := uc.LoginOrSignup(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, session.JWTAccess)

		// Проверка обновления сессии
		user, err := userRepo.GetByUsername(ctx, "existinguser")
		require.NoError(t, err)

		redisSession, err := sessionRepo.Check(ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, session.JWTAccess, redisSession.JWTAccess)
	})

	t.Run("login with wrong password", func(t *testing.T) {
		req := &dto.AuthRequest{
			Username: "wrongpassuser",
			Password: "correctpass",
		}

		// Создаем пользователя
		_, err := uc.LoginOrSignup(ctx, req)
		require.NoError(t, err)

		// Попытка входа с неверным паролем
		req.Password = "wrongpass"
		_, err = uc.LoginOrSignup(ctx, req)
		require.ErrorIs(t, err, entity.ErrWrongCredentials)
	})

	t.Run("duplicate username", func(t *testing.T) {
		req := &dto.AuthRequest{
			Username: "duplicateuser",
			Password: "password123",
		}

		// Первая регистрация
		sessionFirstInstance, err := uc.LoginOrSignup(ctx, req)
		require.NoError(t, err)

		// Вторая попытка регистрации
		sessionSecondInstance, err := uc.LoginOrSignup(ctx, req)
		require.NoError(t, err)
		require.Equal(t, sessionSecondInstance.JWTAccess, sessionFirstInstance.JWTAccess)
		require.Equal(t, sessionSecondInstance.JWTRefresh, sessionFirstInstance.JWTRefresh)
	})

	t.Run("session expiration", func(t *testing.T) {
		req := &dto.AuthRequest{
			Username: "expireduser",
			Password: "password123",
		}

		// Создаем пользователя
		_, err := uc.LoginOrSignup(ctx, req)
		require.NoError(t, err)

		// Проверяем наличие сессии
		user, err := userRepo.GetByUsername(ctx, "expireduser")
		require.NoError(t, err)

		// Ждем истечения срока действия
		time.Sleep(6 * time.Second) // Для теста уменьшаем TTL в конфиге

		// Проверяем отсутствие сессии
		_, err = sessionRepo.Check(ctx, user.ID)
		require.ErrorIs(t, err, entity.ErrNoSession)
	})
}
