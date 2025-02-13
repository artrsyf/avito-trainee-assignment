package usecase

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"

	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/uow"
)

type TransactionUsecaseI interface {
	Create(transactionEntity *entity.Transaction) error
}

type TransactionUsecase struct {
	transactionRepo transactionRepo.TransactionRepositoryI
	userRepo        userRepo.UserRepositoryI
	uow             uow.UnitOfWorkI
}

func NewTransactionUsecase(transactionRepository transactionRepo.TransactionRepositoryI, userRepository userRepo.UserRepositoryI, uow uow.UnitOfWorkI) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo: transactionRepository,
		userRepo:        userRepository,
		uow:             uow,
	}
}

func (uc *TransactionUsecase) Create(transactionEntity *entity.Transaction) error {
	senderUserModel, err := uc.userRepo.GetByUsername(transactionEntity.SenderUsername)
	if err != nil {
		return err
	}

	receiverUserModel, err := uc.userRepo.GetByUsername(transactionEntity.ReceiverUsername)
	if err != nil {
		return err
	}

	if senderUserModel.Coins < transactionEntity.Amount {
		return entity.ErrNotEnoughBalance
	}

	senderUserModel.Coins -= transactionEntity.Amount
	receiverUserModel.Coins += transactionEntity.Amount

	err = uc.uow.Begin(context.Background())
	if err != nil {
		return err
	}

	err = uc.userRepo.Update(uc.uow, senderUserModel)
	if err != nil {
		uc.uow.Rollback()

		return err
	}

	err = uc.userRepo.Update(uc.uow, receiverUserModel)
	if err != nil {
		uc.uow.Rollback()

		return err
	}

	_, err = uc.transactionRepo.Create(transactionEntity)
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
