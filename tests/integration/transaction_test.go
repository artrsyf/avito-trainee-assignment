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
	uowI "github.com/artrsyf/avito-trainee-assignment/pkg/uow"
	uow "github.com/artrsyf/avito-trainee-assignment/pkg/uow/postgres"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestTransactionUsecase_Integration(t *testing.T) {
	userRepo := userRepo.NewUserPostgresRepository(DB, logrus.New())
	transactionRepo := transactionRepo.NewTransactionPostgresRepository(DB, logrus.New())
	uowFactory := uow.NewFactory(DB)

	uc := usecase.NewTransactionUsecase(transactionRepo, userRepo, uowFactory, logrus.New())
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

		senderUser, _ := userRepo.GetByUsername(ctx, "sender")
		require.Equal(t, uint(700), senderUser.Coins)

		receiverUser, _ := userRepo.GetByUsername(ctx, "receiver")
		require.Equal(t, uint(800), receiverUser.Coins)

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

		transaction := &entity.Transaction{
			SenderUsername:   "sender",
			ReceiverUsername: "receiver",
			Amount:           300,
		}

		faultyUowFactory := NewFaultyUOWFactory(DB, 2)
		uc := usecase.NewTransactionUsecase(transactionRepo, userRepo, faultyUowFactory, logrus.New())

		err := uc.Create(ctx, transaction)
		require.Error(t, err)

		senderUser, _ := userRepo.GetByUsername(ctx, "sender")
		require.Equal(t, uint(1000), senderUser.Coins)

		receiverUser, _ := userRepo.GetByUsername(ctx, "receiver")
		require.Equal(t, uint(500), receiverUser.Coins)
	})
}

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

func (u *FaultyUOW) Exec(query string, args ...any) (sql.Result, error) {
	return u.ExecContext(context.Background(), query, args...)
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

func (u *FaultyUOW) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if u.tx == nil {
		return nil
	}
	return u.tx.QueryRowContext(ctx, query, args...)
}

// Реализация фабрики для FaultyUOW
type FaultyUOWFactory struct {
	db     *sql.DB
	failOn int
}

func NewFaultyUOWFactory(db *sql.DB, failOn int) uowI.Factory {
	return &FaultyUOWFactory{db: db, failOn: failOn}
}

func (f *FaultyUOWFactory) NewUnitOfWork() uowI.UnitOfWork {
	return NewFaultyUOW(f.db, f.failOn)
}
