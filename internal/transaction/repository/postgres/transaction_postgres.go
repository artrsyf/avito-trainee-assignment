package postgres

import (
	"context"
	"database/sql"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
	"github.com/sirupsen/logrus"
)

type TransactionPostgresRepository struct {
	DB     *sql.DB
	logger *logrus.Logger
}

func NewTransactionPostgresRepository(db *sql.DB, logger *logrus.Logger) *TransactionPostgresRepository {
	return &TransactionPostgresRepository{
		DB:     db,
		logger: logger,
	}
}

func (repo *TransactionPostgresRepository) Create(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	createdTransaction := model.Transaction{}
	err := repo.DB.QueryRowContext(ctx, "INSERT INTO transactions (sender_user_id, receiver_user_id, amount) VALUES ($1, $2, $3) RETURNING id, sender_user_id, receiver_user_id, amount", transaction.SenderUserID, transaction.ReceiverUserID, transaction.Amount).
		Scan(&createdTransaction.ID, &createdTransaction.SenderUserID, &createdTransaction.ReceiverUserID, &createdTransaction.Amount)
	if err != nil {
		repo.logger.WithError(err).Error("Failed to create transaction")
		return nil, err
	}

	repo.logger.WithFields(logrus.Fields{
		"transaction_id": createdTransaction.ID,
	}).Debug("Created transaction in Postgres")

	return &createdTransaction, nil
}

func (repo *TransactionPostgresRepository) GetReceivedByUserID(ctx context.Context, userID uint) (entity.ReceivedHistory, error) {
	rows, err := repo.DB.QueryContext(ctx, `
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
		repo.logger.WithError(err).Error("Failed to select received transaction group")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.WithError(err).Warn("Failed to close rows selecting received transactions")
		}
	}()

	receivedHistory := entity.ReceivedHistory{}
	for rows.Next() {
		currentReceivedTransactionGroup := entity.ReceivedTransactionGroup{}
		err := rows.Scan(&currentReceivedTransactionGroup.SenderUsername, &currentReceivedTransactionGroup.Amount)
		if err != nil {
			repo.logger.WithError(err).Error("Failed to select received transactions")
			return nil, err
		}

		receivedHistory = append(receivedHistory, currentReceivedTransactionGroup)
	}

	return receivedHistory, nil
}

func (repo *TransactionPostgresRepository) GetSentByUserID(ctx context.Context, userID uint) (entity.SentHistory, error) {
	rows, err := repo.DB.QueryContext(ctx, `
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
		repo.logger.WithError(err).Error("Failed to select sent transaction group")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.WithError(err).Warn("Failed to close rows selecting sent transactions")
		}
	}()

	sentHistory := entity.SentHistory{}
	for rows.Next() {
		currentSentTransactionGroup := entity.SentTransactionGroup{}
		err := rows.Scan(&currentSentTransactionGroup.ReceiverUsername, &currentSentTransactionGroup.Amount)
		if err != nil {
			repo.logger.WithError(err).Error("Failed to select sent transactions")
			return nil, err
		}

		sentHistory = append(sentHistory, currentSentTransactionGroup)
	}

	return sentHistory, nil
}
