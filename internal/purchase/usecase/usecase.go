package usecase

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/dto"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/sirupsen/logrus"

	purchaseRepo "github.com/artrsyf/avito-trainee-assignment/internal/purchase/repository"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
	uowI "github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

type PurchaseUsecaseI interface {
	Create(ctx context.Context, purchaseRequest *dto.PurchaseItemRequest) error
}

type PurchaseUsecase struct {
	purchaseRepo purchaseRepo.PurchaseRepositoryI
	userRepo     userRepo.UserRepositoryI
	uowFactory   uowI.Factory
	logger       *logrus.Logger
}

func NewPurchaseUsecase(
	purchaseRepository purchaseRepo.PurchaseRepositoryI,
	userRepository userRepo.UserRepositoryI,
	uowFactory uowI.Factory,
	logger *logrus.Logger,
) *PurchaseUsecase {
	return &PurchaseUsecase{
		purchaseRepo: purchaseRepository,
		userRepo:     userRepository,
		uowFactory:   uowFactory,
		logger:       logger,
	}
}

func (uc *PurchaseUsecase) Create(
	ctx context.Context,
	purchaseRequest *dto.PurchaseItemRequest,
) error {
	purchaseEntity := dto.PurchaseItemRequestToEntity(purchaseRequest)

	customerModel, err := uc.userRepo.GetByID(ctx, purchaseEntity.PurchaserID)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get customer user by id")
		return err
	}

	purchaseType, err := uc.purchaseRepo.GetProductByType(ctx, purchaseEntity.PurchaseTypeName)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get product by type name")
		return err
	}

	if customerModel.Coins < purchaseType.Cost {
		uc.logger.WithError(err).Error("Customer doesn't have enough balance")
		return entity.ErrNotEnoughBalance
	}

	customerModel.Coins -= purchaseType.Cost

	uow := uc.uowFactory.NewUnitOfWork()

	err = uow.Begin(ctx)
	if err != nil {
		uc.logger.WithError(err).Error("Transaction begin error")
		return err
	}

	err = uc.userRepo.Update(ctx, uow, customerModel)
	if err != nil {
		rbErr := uow.Rollback()
		if rbErr != nil {
			uc.logger.WithError(rbErr).Error("Rollback error encountered")
		}
		uc.logger.WithError(err).Error("Rollback money transfer due user updating")
		return err
	}

	purchaseModel, err := uc.purchaseRepo.Create(ctx, uow, purchaseEntity)
	if err != nil {
		rbErr := uow.Rollback()
		if rbErr != nil {
			uc.logger.WithError(rbErr).Error("Rollback error encountered")
		}
		uc.logger.WithError(err).Error("Rollback money transfer due purchase creating")
		return err
	}

	err = uow.Commit()
	if err != nil {
		uc.logger.WithError(err).Error("Transaction commit error")
		return err
	}

	uc.logger.WithFields(logrus.Fields{
		"customer_id":     purchaseModel.PurchaserID,
		"product_type_id": purchaseModel.PurchaseTypeID,
	}).Info("Successfully purchase product")

	return nil
}
