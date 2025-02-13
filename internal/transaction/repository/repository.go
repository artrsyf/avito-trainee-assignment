package repository

import (
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
)

type TransactionRepositoryI interface {
	Create(transaction *entity.Transaction) (*model.Transaction, error)
}
