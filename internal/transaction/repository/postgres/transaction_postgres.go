package postgres

import (
	"database/sql"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
)

type TransactionPostgresRepository struct {
	DB *sql.DB
}

func NewTransactionPostgresRepository(db *sql.DB) *TransactionPostgresRepository {
	return &TransactionPostgresRepository{
		DB: db,
	}
}

func (repo *TransactionPostgresRepository) Create(transaction *model.Transaction) (*model.Transaction, error) {
	createdTransaction := model.Transaction{}
	err := repo.DB.QueryRow("INSERT INTO transactions (sender_user_id, receiver_user_id, amount) VALUES ($1, $2, $3) RETURNING id, sender_user_id, receiver_user_id, amount", transaction.SenderUserID, transaction.ReceiverUserID, transaction.Amount).
		Scan(&createdTransaction.ID, &createdTransaction.SenderUserID, &createdTransaction.ReceiverUserID, &createdTransaction.Amount)
	if err != nil {
		return nil, err
	}

	return &createdTransaction, nil
}
