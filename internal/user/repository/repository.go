package repository

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

type UserRepositoryI interface {
	Create(ctx context.Context, user *entity.User) (*model.User, error)
	Update(ctx context.Context, uow uow.UnitOfWorkI, user *model.User) error
	GetById(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}
