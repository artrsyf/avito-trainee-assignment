package repository

import (
	"context"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
)

//go:generate mockgen -source=repository.go -destination=mock_repository/transaction_mock.go -package=mock_repository MockTransactionRepository
type TransactionRepositoryI interface {
	Create(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	GetReceivedByUserID(ctx context.Context, userID uint) (entity.ReceivedHistory, error)
	GetSentByUserID(ctx context.Context, userID uint) (entity.SentHistory, error)
}
