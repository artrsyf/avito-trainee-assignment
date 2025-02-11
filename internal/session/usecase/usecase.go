package usecase

import (
	"errors"
	"time"

	sessionDTO "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/dto"
	sessionEntity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	sessionRepo "github.com/artrsyf/avito-trainee-assignment/internal/session/repository"
	userDTO "github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
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
	ok, err := uc.userRepo.GetByUsername(authRequest.Username)
	if err != nil {
		return nil, err
	}

	if ok != nil {
		return nil, errors.New("user is already signed up")
	}

	userEntity, err := userDTO.AuthRequestToEntity(authRequest)
	if err != nil {
		return nil, err
	}

	createdUser, err := uc.userRepo.Create(userEntity)
	if err != nil {
		return nil, err
	}

	sessionEntity, err := sessionDTO.AuthRequestToEntity(
		authRequest,
		createdUser.ID,
		time.Now().Add(15*time.Minute), /*TODO magic dates*/
		time.Now().Add(24*time.Hour),   /*TODO magic dates*/
	)
	if err != nil {
		return nil, err
	}

	_, err = uc.sessionRepo.Check(sessionEntity.UserID)
	if err != nil {
		return nil, err
	}

	createdSessionModel, err := uc.sessionRepo.Create(sessionEntity)
	if err != nil {
		return nil, err
	}

	createdSessionEntity := sessionDTO.SessionModelToEntity(createdSessionModel)

	return createdSessionEntity, nil
}

/*TODO Delete method*/
// func (uc *SessionUsecase) Create(authRequest *sessionDTO.SignupRequest) (*sessionEntity.Session, error) {
// 	sessionEntity, err := sessionDTO.SignupRequestToEntity(
// 		authRequest,

// 		time.Now().Add(15*time.Minute), /*TODO*/
// 		time.Now().Add(24*time.Hour),   /*TODO*/
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = uc.sessionRepo.Check(uint(sessionEntity.UserID))
// 	if err != nil {
// 		return nil, err
// 	}

// 	createdSessionModel, err := uc.sessionRepo.Create(sessionEntity)
// 	if err != nil {
// 		return nil, err
// 	}

// 	createdSessionEntity := sessionDTO.SessionModelToEntity(createdSessionModel)

// 	return createdSessionEntity, nil
// }

// func (uc *SessionUsecase) Check(userID uint) (*sessionEntity.Session, error) {
// 	sessionModel, err := uc.sessionRepo.Check(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sessionEntity := sessionDTO.SessionModelToEntity(sessionModel)

// 	return sessionEntity, nil
// }
