package dto

import (
	sessionDTO "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	"golang.org/x/crypto/bcrypt"
)

func AuthRequestToEntity(authRequest *sessionDTO.AuthRequest) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Username:     authRequest.Username,
		Coins:        1000, /*TODO magic number*/
		PasswordHash: string(hashedPassword),
	}, nil
}

func ModelToEntity(user *model.User) *entity.User {
	return &entity.User{
		Username:     user.Username,
		Coins:        user.Coins,
		PasswordHash: user.PasswordHash,
	}
}
