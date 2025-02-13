package usecase

import (
	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository"
	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/dto"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
)

type UserUsecaseI interface {
	GetInfoById(userID uint) (*dto.GetInfoResponse, error)
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

func (uc *UserUsecase) GetInfoById(userID uint) (*dto.GetInfoResponse, error) {
	userInfo, err := uc.userRepo.GetById(userID)
	if err != nil {
		return nil, err
	}

	userInventory, err := uc.purchaseRepo.GetPurchasesByUserId(userID)
	if err != nil {
		return nil, err
	}

	userSentTransactions, err := uc.transactionRepo.GetSentByUserID(userID)
	if err != nil {
		return nil, err
	}

	userReceivedTransactions, err := uc.transactionRepo.GetReceivedByUserID(userID)
	if err != nil {
		return nil, err
	}

	getInfoResponse := dto.CreateGetInfoResponse(userInfo.Coins, &userInventory, &userSentTransactions, &userReceivedTransactions)

	return getInfoResponse, nil
}
