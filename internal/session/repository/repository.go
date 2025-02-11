package repository

import (
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
)

type SessionRepositoryI interface {
	Create(sessionEntity *entity.Session) (*model.Session, error)
	Check(userID uint) (*model.Session, error)
	Delete(userID uint) error
}
