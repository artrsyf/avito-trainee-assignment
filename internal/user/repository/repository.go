package repository

import (
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/uow"
)

type UserRepositoryI interface {
	Create(user *entity.User) (*model.User, error)
	Update(uow uow.UnitOfWorkI, user *model.User) error
	GetById(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}
