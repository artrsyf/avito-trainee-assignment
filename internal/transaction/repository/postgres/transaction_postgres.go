package postgres

import (
	"database/sql"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
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

func (repo *TransactionPostgresRepository) GetReceivedByUserID(userID uint) (entity.ReceivedHistory, error) {
	rows, err := repo.DB.Query(`
		SELECT
			u1.username, 
			SUM(t.amount)
		FROM 
			transactions t
		JOIN 
			users u1 ON t.sender_user_id = u1.id
		WHERE 
			t.receiver_user_id = $1
		GROUP BY 
			u1.username`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receivedHistory := entity.ReceivedHistory{}
	for rows.Next() {
		currentReceivedTransactionGroup := entity.ReceivedTransactionGroup{}
		err := rows.Scan(&currentReceivedTransactionGroup.SenderUsername, &currentReceivedTransactionGroup.Amount)
		if err != nil {
			return nil, err
		}

		receivedHistory = append(receivedHistory, currentReceivedTransactionGroup)
	}

	return receivedHistory, nil
}

func (repo *TransactionPostgresRepository) GetSentByUserID(userID uint) (entity.SentHistory, error) {
	rows, err := repo.DB.Query(`
		SELECT 
			u1.username,
			SUM(t.amount)
		FROM 
			transactions t
		JOIN 
			users u1 ON t.receiver_user_id = u1.id
		WHERE 
			t.sender_user_id = $1
		GROUP BY 
			u1.username`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sentHistory := entity.SentHistory{}
	for rows.Next() {
		currentSentTransactionGroup := entity.SentTransactionGroup{}
		err := rows.Scan(&currentSentTransactionGroup.ReceiverUsername, &currentSentTransactionGroup.Amount)
		if err != nil {
			return nil, err
		}

		sentHistory = append(sentHistory, currentSentTransactionGroup)
	}

	return sentHistory, nil
}
