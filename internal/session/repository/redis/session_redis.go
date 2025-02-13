package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
	"github.com/gomodule/redigo/redis"
)

type SessionRedisRepository struct {
	redisConn redis.Conn
	mu        *sync.Mutex
}

func NewSessionRedisRepository(conn redis.Conn) *SessionRedisRepository {
	return &SessionRedisRepository{
		redisConn: conn,
		mu:        &sync.Mutex{},
	}
}

func (repo *SessionRedisRepository) Create(ctx context.Context, sessionEntity *entity.Session) (*model.Session, error) {
	mkey := "sessions:" + strconv.FormatUint(uint64(sessionEntity.UserID), 10)

	sessionModel := dto.SessionEntityToModel(sessionEntity)
	sessionSerialized, err := json.Marshal(sessionModel)
	if err != nil {
		return nil, err
	}

	ttlSeconds := int64(time.Until(sessionModel.AccessExpiresAt).Seconds())

	repo.mu.Lock()
	result, err := redis.String(redis.DoContext(repo.redisConn, ctx, "SET", mkey, sessionSerialized, "EX", ttlSeconds))
	repo.mu.Unlock()

	if err != nil {
		return nil, err
	}

	if result != "OK" {
		return nil, fmt.Errorf("unexpected Redis response: %v", result)
	}

	return sessionModel, nil
}

func (repo *SessionRedisRepository) Check(ctx context.Context, userID uint) (*model.Session, error) {
	mkey := "sessions:" + strconv.FormatUint(uint64(userID), 10)

	repo.mu.Lock()
	bytes, err := redis.Bytes(redis.DoContext(repo.redisConn, ctx, "GET", mkey))
	repo.mu.Unlock()

	if err == redis.ErrNil {
		return nil, entity.ErrNoSession
	}

	if err != nil {
		return nil, err
	}

	session := &model.Session{}
	err = json.Unmarshal(bytes, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
