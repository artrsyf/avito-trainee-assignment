package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	transactionRepo "github.com/artrsyf/avito-trainee-assignment/internal/transaction/repository/postgres"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/usecase"
	userRepo "github.com/artrsyf/avito-trainee-assignment/internal/user/repository/postgres"
	uow "github.com/artrsyf/avito-trainee-assignment/pkg/uow/postgres"
	"github.com/stretchr/testify/require"
)

func TestTransactionUsecase_Integration(t *testing.T) {
	// Инициализация репозиториев
	userRepo := userRepo.NewUserPostgresRepository(DB)
	transactionRepo := transactionRepo.NewTransactionPostgresRepository(DB)
	uow := uow.NewPostgresUnitOfWork(DB)

	uc := usecase.NewTransactionUsecase(transactionRepo, userRepo, uow)
	ctx := context.Background()

	t.Run("successful transaction", func(t *testing.T) {
		SetupTestData(t, DB)
		_ = CreateTestUser(t, "sender", 1000)
		_ = CreateTestUser(t, "receiver", 500)

		transaction := &entity.Transaction{
			SenderUsername:   "sender",
			ReceiverUsername: "receiver",
			Amount:           300,
		}

		err := uc.Create(ctx, transaction)
		require.NoError(t, err)

		// Проверка балансов
		senderUser, _ := userRepo.GetByUsername(ctx, "sender")
		require.Equal(t, uint(700), senderUser.Coins)

		receiverUser, _ := userRepo.GetByUsername(ctx, "receiver")
		require.Equal(t, uint(800), receiverUser.Coins)

		// Проверка записи транзакции
		var txCount int
		DB.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&txCount)
		require.Equal(t, 1, txCount)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		SetupTestData(t, DB)
		_ = CreateTestUser(t, "poor_sender", 200)
		_ = CreateTestUser(t, "receiver", 500)

		transaction := &entity.Transaction{
			SenderUsername:   "poor_sender",
			ReceiverUsername: "receiver",
			Amount:           300,
		}

		err := uc.Create(ctx, transaction)
		require.ErrorIs(t, err, entity.ErrNotEnoughBalance)

		// Проверка неизменности балансов
		senderUser, _ := userRepo.GetByUsername(ctx, "poor_sender")
		require.Equal(t, uint(200), senderUser.Coins)

		receiverUser, _ := userRepo.GetByUsername(ctx, "receiver")
		require.Equal(t, uint(500), receiverUser.Coins)
	})

	t.Run("sender not found", func(t *testing.T) {
		SetupTestData(t, DB)
		_ = CreateTestUser(t, "receiver", 500)

		transaction := &entity.Transaction{
			SenderUsername:   "unknown",
			ReceiverUsername: "receiver",
			Amount:           100,
		}

		err := uc.Create(ctx, transaction)
		require.Error(t, err)
	})

	t.Run("receiver not found", func(t *testing.T) {
		SetupTestData(t, DB)
		_ = CreateTestUser(t, "sender", 1000)

		transaction := &entity.Transaction{
			SenderUsername:   "sender",
			ReceiverUsername: "unknown",
			Amount:           100,
		}

		err := uc.Create(ctx, transaction)
		require.Error(t, err)
	})

	t.Run("transaction rollback on update error", func(t *testing.T) {
		SetupTestData(t, DB)
		_ = CreateTestUser(t, "sender", 1000)
		_ = CreateTestUser(t, "receiver", 500)

		// Сломаем таблицу пользователей после первого обновления
		transaction := &entity.Transaction{
			SenderUsername:   "sender",
			ReceiverUsername: "receiver",
			Amount:           300,
		}

		// Модифицируем UOW для эмуляции ошибки
		uow := &FaultyUOW{db: DB, failOn: 2}
		uc := usecase.NewTransactionUsecase(transactionRepo, userRepo, uow)

		err := uc.Create(ctx, transaction)
		require.Error(t, err)

		// Проверка отката
		senderUser, _ := userRepo.GetByUsername(ctx, "sender")
		require.Equal(t, uint(1000), senderUser.Coins)

		receiverUser, _ := userRepo.GetByUsername(ctx, "receiver")
		require.Equal(t, uint(500), receiverUser.Coins)
	})
}

// Вспомогательные функции и структуры

// func setupTestData(t *testing.T) {
// 	_, err := db.Exec(`
// 		DELETE FROM users;
// 		DELETE FROM transactions;
// 	`)
// 	require.NoError(t, err)
// }

// func createTestUser(t *testing.T, username string, coins uint) uint {
// 	var id uint
// 	err := db.QueryRow(
// 		"INSERT INTO users (username, coins, password_hash) VALUES ($1, $2, 'hash') RETURNING id",
// 		username, coins,
// 	).Scan(&id)
// 	require.NoError(t, err)
// 	return id
// }

// FaultyUOW — реализация UnitOfWorkI для тестирования откатов
type FaultyUOW struct {
	db      *sql.DB
	tx      *sql.Tx
	counter int
	failOn  int
}

func NewFaultyUOW(db *sql.DB, failOn int) *FaultyUOW {
	return &FaultyUOW{db: db, failOn: failOn}
}

func (u *FaultyUOW) Begin(ctx context.Context) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.tx = tx
	return nil
}

func (u *FaultyUOW) Exec(query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, errors.New("transaction has not been started")
	}

	u.counter++
	if u.counter == u.failOn {
		return nil, errors.New("simulated execution failure")
	}

	return u.tx.Exec(query, args...)
}

func (u *FaultyUOW) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, errors.New("transaction has not been started")
	}

	u.counter++
	if u.counter == u.failOn {
		return nil, errors.New("simulated execution failure")
	}

	return u.tx.ExecContext(ctx, query, args...)
}

func (u *FaultyUOW) Commit() error {
	if u.tx == nil {
		return errors.New("transaction has not been started")
	}

	u.counter++
	if u.counter == u.failOn {
		return errors.New("simulated commit failure")
	}

	return u.tx.Commit()
}

func (u *FaultyUOW) Rollback() error {
	if u.tx == nil {
		return errors.New("transaction has not been started")
	}

	fmt.Println("Rolling back transaction...")
	return u.tx.Rollback()
}
