package repository

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
)

type SessionRepositoryI interface {
	Create(ctx context.Context, sessionEntity *entity.Session) (*model.Session, error)
	Check(ctx context.Context, userID uint) (*model.Session, error)
}
