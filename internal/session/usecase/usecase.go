package usecase

import (
	"time"

	sessionDTO "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	sessionEntity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	sessionRepo "github.com/artrsyf/avito-trainee-assignment/internal/session/repository"
	userDTO "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
	userEntity "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
)

type SessionUsecaseI interface {
	Signup(authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error)
	/*TODO Delete method*/
	// Create(authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error)
	/*TODO Implement method*/
	// Check(userID uint) (*sessionEntity.Session, error)
}

type SessionUsecase struct {
	sessionRepo sessionRepo.SessionRepositoryI
	userRepo    userRepo.UserRepositoryI
}

func NewSessionUsecase(sessionRepository sessionRepo.SessionRepositoryI, userRepository userRepo.UserRepositoryI) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepository,
		userRepo:    userRepository,
	}
}

func (uc *SessionUsecase) Signup(authRequest *sessionDTO.AuthRequest) (*sessionEntity.Session, error) {
	_, err := uc.userRepo.GetByUsername(authRequest.Username)
	if err == nil {
		return nil, userEntity.ErrAlreadyCreated
	}

	if err != userEntity.ErrIsNotExist {
		return nil, err
	}

	userEntity, err := userDTO.AuthRequestToEntity(authRequest)
	if err != nil {
		return nil, err
	}

	createdUser, err := uc.userRepo.Create(userEntity)
	if err != nil {
		return nil, err
	}

	session, err := sessionDTO.AuthRequestToEntity(
		authRequest,
		createdUser.ID,
		time.Now().Add(15*time.Minute), /*TODO magic dates*/
		time.Now().Add(24*time.Hour),   /*TODO magic dates*/
	)
	if err != nil {
		return nil, err
	}

	_, err = uc.sessionRepo.Check(session.UserID)
	if err == nil {
		return nil, sessionEntity.ErrAlreadyCreated
	}

	if err != sessionEntity.ErrNoSession {
		return nil, err
	}

	createdSessionModel, err := uc.sessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	createdSessionEntity := sessionDTO.SessionModelToEntity(createdSessionModel)

	return createdSessionEntity, nil
}
