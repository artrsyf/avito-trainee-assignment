package repository

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

//go:generate mockgen -source=repository.go -destination=mock_repository/user_mock.go -package=mock_repository MockUserRepository
type UserRepositoryI interface {
	Create(ctx context.Context, user *entity.User) (*model.User, error)
	Update(ctx context.Context, uow uow.UnitOfWorkI, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}
