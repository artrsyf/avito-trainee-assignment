package dto

import (
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
)

func ModelToEntity(user *model.User) *entity.User {
	return &entity.User{
		Username:     user.Username,
		Coins:        user.Coins,
		PasswordHash: user.PasswordHash,
	}
}
