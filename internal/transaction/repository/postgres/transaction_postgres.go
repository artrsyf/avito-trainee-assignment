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

func (repo *TransactionPostgresRepository) Create(transaction *entity.Transaction) (*model.Transaction, error) {
	createdTransaction := model.Transaction{}
	err := repo.DB.QueryRow("INSERT INTO transactions (sender_username, receiver_username, amount) VALUES ($1, $2, $3) RETURNING id, sender_username, receiver_username, amount", transaction.SenderUsername, transaction.ReceiverUsername, transaction.Amount).
		Scan(&createdTransaction.ID, &createdTransaction.SenderUsername, &createdTransaction.ReceiverUsername, &createdTransaction.Amount)
	if err != nil {
		return nil, err
	}

	return &createdTransaction, nil
}

// func (repo *UserPostgresRepository) GetById(id uint) (*model.User, error) {
// 	user := model.User{}

// 	err := repo.DB.
// 		QueryRow("SELECT id, username, coins, password_hash FROM users WHERE id = $1", id).
// 		Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
// 	if err == sql.ErrNoRows {
// 		return nil, entity.ErrIsNotExist
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (repo *UserPostgresRepository) GetByUsername(username string) (*model.User, error) {
// 	user := model.User{}

// 	err := repo.DB.
// 		QueryRow("SELECT id, username, coins, password_hash FROM users WHERE username = $1", username).
// 		Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
// 	if err == sql.ErrNoRows {
// 		return nil, entity.ErrIsNotExist
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }
