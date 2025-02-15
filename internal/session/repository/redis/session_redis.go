package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type SessionRedisRepository struct {
	client *redis.Client
	logger *logrus.Logger
}

func NewSessionRedisRepository(
	client *redis.Client,
	logger *logrus.Logger,
) *SessionRedisRepository {
	return &SessionRedisRepository{
		client: client,
		logger: logger,
	}
}

func (repo *SessionRedisRepository) Create(
	ctx context.Context,
	sessionEntity *entity.Session,
) (*model.Session, error) {
	mkey := "sessions:" + strconv.FormatUint(uint64(sessionEntity.UserID), 10)

	sessionModel := dto.SessionEntityToModel(sessionEntity)
	sessionSerialized, err := json.Marshal(sessionModel)
	if err != nil {
		repo.logger.WithError(err).Error("Failed to serialize session model")
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	ttl := time.Until(sessionModel.AccessExpiresAt)

	err = repo.client.SetEx(
		ctx,
		mkey,
		sessionSerialized,
		ttl,
	).Err()

	if err != nil {
		repo.logger.WithError(err).Error("Failed to set session in Redis")
		return nil, fmt.Errorf("redis error: %w", err)
	}

	repo.logger.WithFields(logrus.Fields{
		"user_id": sessionEntity.UserID,
	}).Debug("Created session in Redis")

	return sessionModel, nil
}

func (repo *SessionRedisRepository) Check(
	ctx context.Context,
	userID uint,
) (*model.Session, error) {
	mkey := "sessions:" + strconv.FormatUint(uint64(userID), 10)

	data, err := repo.client.Get(ctx, mkey).Bytes()
	if err == redis.Nil {
		repo.logger.WithFields(logrus.Fields{
			"user_id": userID,
		}).Debug("Couldn't find user session")
		return nil, entity.ErrNoSession
	}
	if err != nil {
		repo.logger.WithError(err).Error("Failed to get session from Redis")
		return nil, fmt.Errorf("redis error: %w", err)
	}

	session := &model.Session{}
	if err := json.Unmarshal(data, session); err != nil {
		repo.logger.WithError(err).Error("Failed to deserialize session in model")
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return session, nil
}
