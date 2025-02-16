package usecase

import (
	"context"
	"time"

	"github.com/artrsyf/avito-trainee-assignment/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	sessionDTO "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	sessionEntity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	sessionRepo "github.com/artrsyf/avito-trainee-assignment/internal/session/repository"
	userDTO "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
)

type SessionUsecaseI interface {
	LoginOrSignup(
		ctx context.Context,
		authRequest *sessionDTO.AuthRequest,
	) (*sessionEntity.Session, error)
}

type SessionUsecase struct {
	sessionRepo sessionRepo.SessionRepositoryI
	userRepo    userRepo.UserRepositoryI
	userConfig  config.UserConfig
	logger      *logrus.Logger
}

func NewSessionUsecase(
	sessionRepository sessionRepo.SessionRepositoryI,
	userRepository userRepo.UserRepositoryI,
	cfg config.UserConfig,
	logger *logrus.Logger,
) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
		userConfig:  cfg,
		logger:      logger,
	}
}

func (uc *SessionUsecase) LoginOrSignup(
	ctx context.Context,
	authRequest *sessionDTO.AuthRequest,
) (*sessionEntity.Session, error) {
	userModel, err := uc.userRepo.GetByUsername(ctx, authRequest.Username)
	if err != nil && err != userEntity.ErrIsNotExist {
		uc.logger.WithError(err).Error("Failed to get user by username")
		return nil, err
	}

	if userModel != nil {
		if !checkPassword(authRequest.Password, userModel.PasswordHash) {
			uc.logger.Info("Failed to authenticate user")
			return nil, sessionEntity.ErrWrongCredentials
		}

		sessionModel, checkErr := uc.sessionRepo.Check(ctx, userModel.ID)
		if checkErr == sessionEntity.ErrNoSession {
			return uc.grantSession(ctx, userModel.ID, authRequest)
		}
		if checkErr != nil {
			uc.logger.WithError(checkErr).Error("An error occured due checking user session")
			return nil, checkErr
		}

		session := sessionDTO.SessionModelToEntity(sessionModel)
		return session, nil
	}

	user, err := userDTO.AuthRequestToEntity(
		authRequest,
		uc.userConfig.InitCoinsBalance,
	)
	if err != nil {
		uc.logger.WithError(err).WithField(
			"wrong request", authRequest,
		).Error("Failed cast request to entity")
		return nil, err
	}

	createdUserModel, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		uc.logger.WithError(err).WithField(
			"broken user", user,
		).Error("Failed create new user")
		return nil, err
	}

	return uc.grantSession(ctx, createdUserModel.ID, authRequest)
}

func checkPassword(inputPassword, storedPasswordHash string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(storedPasswordHash),
		[]byte(inputPassword),
	) == nil
}

func (uc *SessionUsecase) grantSession(
	ctx context.Context,
	userID uint,
	authRequest *sessionDTO.AuthRequest,
) (*sessionEntity.Session, error) {
	accessTokenExpiration, err := uc.userConfig.Auth.GetAccessTokenExpiration()
	if err != nil {
		uc.logger.WithError(err).Error("Failed to parse access token expiration")
		return nil, err
	}

	refreshTokenExpiration, err := uc.userConfig.Auth.GetRefreshTokenExpiration()
	if err != nil {
		uc.logger.WithError(err).Error("Failed to parse refresh token expiration")
		return nil, err
	}

	session, err := sessionDTO.AuthRequestToEntity(
		authRequest,
		userID,
		time.Now().Add(accessTokenExpiration),
		time.Now().Add(refreshTokenExpiration),
	)
	if err != nil {
		uc.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":                 userID,
			"auth_request":            authRequest,
			"access_token_duratuion":  accessTokenExpiration,
			"refresh_token_duratuion": refreshTokenExpiration,
		}).Error("Failed cast session request to entity")
		return nil, err
	}

	createdSessionModel, err := uc.sessionRepo.Create(ctx, session)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to create user session")
		return nil, err
	}

	uc.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"username": authRequest.Username,
	}).Info("Granted new session")

	return sessionDTO.SessionModelToEntity(createdSessionModel), nil
}
