package usecase

import (
	"time"

	"github.com/artrsyf/avito-trainee-assignment/config"

	sessionDTO "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	sessionEntity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	sessionRepo "github.com/artrsyf/avito-trainee-assignment/internal/session/repository"
	userDTO "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
)

type SessionUsecaseI interface {
	LoginOrSignup(authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error)
	/*TODO Implement method*/
	// Check(userID uint) (*sessionEntity.Session, error)
}

type SessionUsecase struct {
	sessionRepo sessionRepo.SessionRepositoryI
	userRepo    userRepo.UserRepositoryI
	authConfig  config.AuthConfig
}

func NewSessionUsecase(sessionRepository sessionRepo.SessionRepositoryI, userRepository userRepo.UserRepositoryI, config config.AuthConfig) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
		authConfig:  config,
	}
}

func (uc *SessionUsecase) LoginOrSignup(authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error) {
	userModel, err := uc.userRepo.GetByUsername(authRequest.Username)
	if err != nil && err != userEntity.ErrIsNotExist {
		return nil, err
	}

	if userModel != nil {
		sessionModel, err := uc.sessionRepo.Check(userModel.ID)
		if err == sessionEntity.ErrNoSession {
			return uc.grantSession(userModel.ID, authRequest)
		}
		if err != nil {
			return nil, err
		}

		session := sessionDTO.SessionModelToEntity(sessionModel)
		return session, nil
	}

	user, err := userDTO.AuthRequestToEntity(authRequest)
	if err != nil {
		return nil, err
	}

	createdUserModel, err := uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return uc.grantSession(createdUserModel.ID, authRequest)
}

func (uc *SessionUsecase) grantSession(userID uint, authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error) {
	accessTokenExpiration, err := uc.authConfig.GetAccessTokenExpiration()
	if err != nil {
		return nil, err
	}

	refreshTokenExpiration, err := uc.authConfig.GetRefreshTokenExpiration()
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

	createdSessionModel, err := uc.sessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	return sessionDTO.SessionModelToEntity(createdSessionModel), nil
}
