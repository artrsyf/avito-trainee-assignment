package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
)

func TestSessionRedisRepository_Create(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := NewSessionRedisRepository(db, logrus.New())

	now := time.Now()
	sessionEntity := &entity.Session{
		JWTAccess:        "access_token",
		JWTRefresh:       "refresh_token",
		UserID:           1,
		Username:         "testuser",
		AccessExpiresAt:  now.Add(1 * time.Hour),
		RefreshExpiresAt: now.Add(24 * time.Hour),
	}

	t.Run("Success", func(t *testing.T) {
		sessionModel := model.Session{
			JWTAccess:        "access_token",
			JWTRefresh:       "refresh_token",
			UserID:           1,
			Username:         "testuser",
			AccessExpiresAt:  now.Add(1 * time.Hour),
			RefreshExpiresAt: now.Add(24 * time.Hour),
		}

		serialized, _ := json.Marshal(sessionModel)
		ttl := time.Until(sessionModel.AccessExpiresAt)

		mock.ExpectSetEx("sessions:1", serialized, ttl).SetVal("OK")

		result, err := repo.Create(ctx, sessionEntity)

		assert.NoError(t, err)
		assert.Equal(t, &sessionModel, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("MarshalError", func(t *testing.T) {
		invalidEntity := &entity.Session{
			AccessExpiresAt: time.Now().Add(1 * time.Hour),
		}

		_, err := repo.Create(ctx, invalidEntity)
		assert.Error(t, err)
	})

	t.Run("RedisError", func(t *testing.T) {
		mock.ExpectSetEx("sessions:1", nil, 0).SetErr(errors.New("redis error"))

		_, err := repo.Create(ctx, sessionEntity)
		assert.ErrorContains(t, err, "redis error")
	})
}

func TestSessionRedisRepository_Check(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := NewSessionRedisRepository(db, logrus.New())

	fixedTime := time.Date(2025, time.February, 14, 12, 0, 0, 0, time.UTC)

	validSession := model.Session{
		JWTAccess:        "access_token",
		JWTRefresh:       "refresh_token",
		UserID:           1,
		Username:         "testuser",
		AccessExpiresAt:  fixedTime.Add(1 * time.Hour),
		RefreshExpiresAt: fixedTime.Add(24 * time.Hour),
	}
	serialized, _ := json.Marshal(validSession)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectGet("sessions:1").SetVal(string(serialized))

		session, err := repo.Check(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, &validSession, session)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectGet("sessions:2").RedisNil()

		_, err := repo.Check(ctx, 2)
		assert.ErrorIs(t, err, entity.ErrNoSession)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("RedisError", func(t *testing.T) {
		mock.ExpectGet("sessions:3").SetErr(errors.New("connection error"))

		_, err := repo.Check(ctx, 3)
		assert.ErrorContains(t, err, "connection error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InvalidData", func(t *testing.T) {
		mock.ExpectGet("sessions:4").SetVal("{invalid json}")

		_, err := repo.Check(ctx, 4)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestNewSessionRedisRepository(t *testing.T) {
	db := redis.NewClient(&redis.Options{})
	repo := NewSessionRedisRepository(db, logrus.New())
	assert.NotNil(t, repo)
}
