package usecase

import (
	"context"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/config"
	"golang.org/x/crypto/bcrypt"

	sessionDTO "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	sessionEntity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	sessionRepo "github.com/artrsyf/avito-trainee-assignment/internal/session/repository"
	userDTO "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
)

type SessionUsecaseI interface {
	LoginOrSignup(ctx context.Context, authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error)
	/*TODO Implement method*/
	// Check(userID uint) (*sessionEntity.Session, error)
}

type SessionUsecase struct {
	sessionRepo sessionRepo.SessionRepositoryI
	userRepo    userRepo.UserRepositoryI
	userConfig  config.UserConfig
}

func NewSessionUsecase(sessionRepository sessionRepo.SessionRepositoryI, userRepository userRepo.UserRepositoryI, config config.UserConfig) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
		userConfig:  config,
	}
}

func (uc *SessionUsecase) LoginOrSignup(ctx context.Context, authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error) {
	userModel, err := uc.userRepo.GetByUsername(ctx, authRequest.Username)
	if err != nil && err != userEntity.ErrIsNotExist {
		return nil, err
	}

	if userModel != nil {
		if !checkPassword(authRequest.Password, userModel.PasswordHash) {
			return nil, sessionEntity.ErrWrongCredentials
		}

		sessionModel, err := uc.sessionRepo.Check(ctx, userModel.ID)
		if err == sessionEntity.ErrNoSession {
			return uc.grantSession(ctx, userModel.ID, authRequest)
		}
		if err != nil {
			return nil, err
		}

		session := sessionDTO.SessionModelToEntity(sessionModel)
		return session, nil
	}

	user, err := userDTO.AuthRequestToEntity(authRequest, uc.userConfig.InitCoinsBalance)
	if err != nil {
		return nil, err
	}

	createdUserModel, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return uc.grantSession(ctx, createdUserModel.ID, authRequest)
}

func checkPassword(inputPassword, storedPasswordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(inputPassword)) == nil
}

func (uc *SessionUsecase) grantSession(ctx context.Context, userID uint, authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error) {
	accessTokenExpiration, err := uc.userConfig.Auth.GetAccessTokenExpiration()
	if err != nil {
		return nil, err
	}

	refreshTokenExpiration, err := uc.userConfig.Auth.GetRefreshTokenExpiration()
	if err != nil {
		return nil, err
	}

	session, err := sessionDTO.AuthRequestToEntity(
		authRequest,
		userID,
		time.Now().Add(accessTokenExpiration), /*TODO*/
		time.Now().Add(refreshTokenExpiration),
	)
	if err != nil {
		return nil, err
	}

	createdSessionModel, err := uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return sessionDTO.SessionModelToEntity(createdSessionModel), nil
}
