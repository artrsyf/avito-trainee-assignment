package repository

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
)

//go:generate mockgen -source=repository.go -destination=mock_repository/session_mock.go -package=mock_repository MockSessionRepository
type SessionRepositoryI interface {
	Create(ctx context.Context, sessionEntity *entity.Session) (*model.Session, error)
	Check(ctx context.Context, userID uint) (*model.Session, error)
}
