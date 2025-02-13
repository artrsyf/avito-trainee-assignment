package usecase

import (
	"context"

	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository"
	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
)

type UserUsecaseI interface {
	GetInfoById(ctx context.Context, userID uint) (*dto.GetInfoResponse, error)
}

type UserUsecase struct {
	purchaseRepo    purchaseRepo.PurchaseRepositoryI
	transactionRepo transactionRepo.TransactionRepositoryI
	userRepo        userRepo.UserRepositoryI
}

func NewUserUsecase(purchaseRepository purchaseRepo.PurchaseRepositoryI, transactionRepository transactionRepo.TransactionRepositoryI, userRepository userRepo.UserRepositoryI) *UserUsecase {
	return &UserUsecase{
		purchaseRepo:    purchaseRepository,
		transactionRepo: transactionRepository,
		userRepo:        userRepository,
	}
}

func (uc *UserUsecase) GetInfoById(ctx context.Context, userID uint) (*dto.GetInfoResponse, error) {
	userInfo, err := uc.userRepo.GetById(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInventory, err := uc.purchaseRepo.GetPurchasesByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	userSentTransactions, err := uc.transactionRepo.GetSentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userReceivedTransactions, err := uc.transactionRepo.GetReceivedByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	getInfoResponse := dto.CreateGetInfoResponse(userInfo.Coins, &userInventory, &userSentTransactions, &userReceivedTransactions)

	return getInfoResponse, nil
}
