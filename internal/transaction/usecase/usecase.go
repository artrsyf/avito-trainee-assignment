package usecase

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
	"github.com/sirupsen/logrus"

	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
	uowI "github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

type TransactionUsecaseI interface {
	Create(ctx context.Context, transactionEntity *entity.Transaction) error
}

type TransactionUsecase struct {
	transactionRepo transactionRepo.TransactionRepositoryI
	userRepo        userRepo.UserRepositoryI
	uow             uowI.UnitOfWorkI
	logger          *logrus.Logger
}

func NewTransactionUsecase(
	transactionRepository transactionRepo.TransactionRepositoryI,
	userRepository userRepo.UserRepositoryI,
	uow uowI.UnitOfWorkI,
	logger *logrus.Logger,
) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo: transactionRepository,
		userRepo:        userRepository,
		uow:             uow,
		logger:          logger,
	}
}

func (uc *TransactionUsecase) Create(
	ctx context.Context,
	transactionEntity *entity.Transaction,
) error {
	senderUserModel, err := uc.userRepo.GetByUsername(
		ctx,
		transactionEntity.SenderUsername,
	)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get sender user by username")
		return err
	}

	receiverUserModel, err := uc.userRepo.GetByUsername(
		ctx,
		transactionEntity.ReceiverUsername,
	)
	if err != nil {
		uc.logger.WithError(err).Error("Failed to get receiver user by username")
		return err
	}

	if senderUserModel.Coins < transactionEntity.Amount {
		uc.logger.WithError(err).Error("Sender user doesn't have enough balance")
		return entity.ErrNotEnoughBalance
	}

	senderUserModel.Coins -= transactionEntity.Amount
	receiverUserModel.Coins += transactionEntity.Amount

	err = uc.uow.Begin(ctx)
	if err != nil {
		uc.logger.WithError(err).Error("Transaction begin error")
		return err
	}

	err = uc.userRepo.Update(ctx, uc.uow, senderUserModel)
	if err != nil {
		rbErr := uc.uow.Rollback()
		if rbErr != nil {
			uc.logger.WithError(rbErr).Error("Rollback error encountered")
		}
		uc.logger.WithError(err).Error("Rollback money transfer due user updating")
		return err
	}

	err = uc.userRepo.Update(ctx, uc.uow, receiverUserModel)
	if err != nil {
		rbErr := uc.uow.Rollback()
		if rbErr != nil {
			uc.logger.WithError(rbErr).Error("Rollback error encountered")
		}
		uc.logger.WithError(err).Error("Rollback money transfer due user updating")
		return err
	}

	transactionModel := &model.Transaction{
		SenderUserID:   senderUserModel.ID,
		ReceiverUserID: receiverUserModel.ID,
		Amount:         transactionEntity.Amount,
	}
	_, err = uc.transactionRepo.Create(ctx, transactionModel)
	if err != nil {
		rbErr := uc.uow.Rollback()
		if rbErr != nil {
			uc.logger.WithError(rbErr).Error("Rollback error encountered")
		}
		uc.logger.WithError(err).Error("Rollback money transfer due transaction creating")
		return err
	}

	err = uc.uow.Commit()
	if err != nil {
		uc.logger.WithError(err).Error("Transaction commit error")
		return err
	}

	uc.logger.WithFields(logrus.Fields{
		"sender_username":   transactionEntity.SenderUsername,
		"receiver_username": transactionEntity.ReceiverUsername,
		"amount":            transactionEntity.Amount,
	}).Info("Successfully create transaction")

	return nil
}
