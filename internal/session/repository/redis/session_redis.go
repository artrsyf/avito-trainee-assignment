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
)

type SessionRedisRepository struct {
	client *redis.Client
}

func NewSessionRedisRepository(client *redis.Client) *SessionRedisRepository {
	return &SessionRedisRepository{
		client: client,
	}
}

func (repo *SessionRedisRepository) Create(ctx context.Context, sessionEntity *entity.Session) (*model.Session, error) {
	mkey := "sessions:" + strconv.FormatUint(uint64(sessionEntity.UserID), 10)

	sessionModel := dto.SessionEntityToModel(sessionEntity)
	sessionSerialized, err := json.Marshal(sessionModel)
	if err != nil {
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
		return nil, fmt.Errorf("redis error: %w", err)
	}

	return sessionModel, nil
}

func (repo *SessionRedisRepository) Check(ctx context.Context, userID uint) (*model.Session, error) {
	mkey := "sessions:" + strconv.FormatUint(uint64(userID), 10)

	data, err := repo.client.Get(ctx, mkey).Bytes()
	if err == redis.Nil {
		return nil, entity.ErrNoSession
	}
	if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	session := &model.Session{}
	if err := json.Unmarshal(data, session); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return session, nil
}
