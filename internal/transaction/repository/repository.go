package repository

import (
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
)

type TransactionRepositoryI interface {
	Create(transaction *model.Transaction) (*model.Transaction, error)
	GetReceivedByUserID(userID uint) (entity.ReceivedHistory, error)
	GetSentByUserID(userID uint) (entity.SentHistory, error)
}
