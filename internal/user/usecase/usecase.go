package usecase

import (
	"context"

	"github.com/sirupsen/logrus"

	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository"
	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
)

type UserUsecaseI interface {
	GetInfoByID(ctx context.Context, userID uint) (*dto.GetInfoResponse, error)
}

type UserUsecase struct {
	purchaseRepo    purchaseRepo.PurchaseRepositoryI
	transactionRepo transactionRepo.TransactionRepositoryI
	userRepo        userRepo.UserRepositoryI
	logger          *logrus.Logger
}

func NewUserUsecase(
	purchaseRepository purchaseRepo.PurchaseRepositoryI,
	transactionRepository transactionRepo.TransactionRepositoryI,
	userRepository userRepo.UserRepositoryI,
	logger *logrus.Logger,
) *UserUsecase {
	return &UserUsecase{
		purchaseRepo:    purchaseRepository,
		transactionRepo: transactionRepository,
		userRepo:        userRepository,
		logger:          logger,
	}
}

func (uc *UserUsecase) GetInfoByID(
	ctx context.Context,
	userID uint,
) (*dto.GetInfoResponse, error) {
	userInfo, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get user by id")
		return nil, err
	}

	userInventory, err := uc.purchaseRepo.GetPurchasesByUserID(ctx, userID)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get user purchases by user id")
		return nil, err
	}

	userSentTransactions, err := uc.transactionRepo.GetSentByUserID(ctx, userID)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get user sent transactions by user id")
		return nil, err
	}

	userReceivedTransactions, err := uc.transactionRepo.GetReceivedByUserID(ctx, userID)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get user received transactions by user id")
		return nil, err
	}

	getInfoResponse := dto.CreateGetInfoResponse(
		userInfo.Coins,
		&userInventory,
		&userSentTransactions,
		&userReceivedTransactions,
	)

	uc.logger.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("Successfully got user info")

	return getInfoResponse, nil
}
