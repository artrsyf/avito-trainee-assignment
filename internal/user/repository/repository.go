package repository

import (
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
)

type UserRepositoryI interface {
	Create(user *entity.User) (*model.User, error)
	GetById(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
}
