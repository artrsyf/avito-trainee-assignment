package usecase

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"

	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

type PurchaseUsecaseI interface {
	Create(purchaseRequest *dto.PurchaseItemRequest) error
}

type PurchaseUsecase struct {
	purchaseRepo purchaseRepo.PurchaseRepositoryI
	userRepo     userRepo.UserRepositoryI
	uow          uow.UnitOfWorkI
}

func NewTransactionUsecase(purchaseRepository purchaseRepo.PurchaseRepositoryI, userRepository userRepo.UserRepositoryI, uow uow.UnitOfWorkI) *PurchaseUsecase {
	return &PurchaseUsecase{
		purchaseRepo: purchaseRepository,
		userRepo:     userRepository,
		uow:          uow,
	}
}

func (uc *PurchaseUsecase) Create(purchaseRequest *dto.PurchaseItemRequest) error {
	purchaseEntity := dto.PurchaseItemRequestToEntity(purchaseRequest)

	customerModel, err := uc.userRepo.GetById(purchaseEntity.PurchaserId)
	if err != nil {
		return err
	}

	purchaseType, err := uc.purchaseRepo.GetProductByType(purchaseEntity.PurchaseTypeName)
	if err != nil {
		return err
	}

	if customerModel.Coins < purchaseType.Cost {
		return entity.ErrNotEnoughBalance
	}

	customerModel.Coins -= purchaseType.Cost

	err = uc.uow.Begin(context.Background())
	if err != nil {
		return err
	}

	err = uc.userRepo.Update(uc.uow, customerModel)
	if err != nil {
		uc.uow.Rollback()

		return err
	}

	_, err = uc.purchaseRepo.Create(purchaseEntity)
	if err != nil {
		uc.uow.Rollback()

		return err
	}

	err = uc.uow.Commit()
	if err != nil {
		return err
	}

	return nil
}
