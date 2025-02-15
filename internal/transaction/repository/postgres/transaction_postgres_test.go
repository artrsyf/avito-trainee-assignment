package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/transaction/domain/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTransactionPostgresRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionPostgresRepository(db, logrus.New())

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO transactions .* RETURNING .*").
			WithArgs(1, 2, 100).
			WillReturnRows(sqlmock.NewRows([]string{"id", "sender_user_id", "receiver_user_id", "amount"}).
				AddRow(1, 1, 2, 100))

		tx, err := repo.Create(context.Background(), &model.Transaction{
			SenderUserID:   1,
			ReceiverUserID: 2,
			Amount:         100,
		})

		assert.NoError(t, err)
		assert.Equal(t, &model.Transaction{
			ID:             1,
			SenderUserID:   1,
			ReceiverUserID: 2,
			Amount:         100,
		}, tx)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery("INSERT INTO transactions .* RETURNING .*").
			WithArgs(1, 2, 100).
			WillReturnError(expectedErr)

		_, err := repo.Create(context.Background(), &model.Transaction{
			SenderUserID:   1,
			ReceiverUserID: 2,
			Amount:         100,
		})

		assert.Equal(t, expectedErr, err)
	})
}

func TestTransactionPostgresRepository_GetReceivedByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionPostgresRepository(db, logrus.New())
	userID := uint(1)

	t.Run("SuccessWithData", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username", "sum"}).
			AddRow("user1", 200).
			AddRow("user2", 300)

		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnRows(rows)

		result, err := repo.GetReceivedByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, entity.ReceivedHistory{
			{SenderUsername: "user1", Amount: 200},
			{SenderUsername: "user2", Amount: 300},
		}, result)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username", "sum"})

		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnRows(rows)

		result, err := repo.GetReceivedByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("QueryError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnError(expectedErr)

		_, err := repo.GetReceivedByUserID(context.Background(), userID)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username", "sum"}).
			AddRow("user1", "invalid_amount")

		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnRows(rows)

		_, err := repo.GetReceivedByUserID(context.Background(), userID)

		assert.ErrorContains(t, err, "converting driver.Value type string")
	})
}

func TestTransactionPostgresRepository_GetSentByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTransactionPostgresRepository(db, logrus.New())
	userID := uint(1)

	t.Run("SuccessWithData", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username", "sum"}).
			AddRow("user3", 150).
			AddRow("user4", 250)

		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnRows(rows)

		result, err := repo.GetSentByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, entity.SentHistory{
			{ReceiverUsername: "user3", Amount: 150},
			{ReceiverUsername: "user4", Amount: 250},
		}, result)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username", "sum"})

		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnRows(rows)

		result, err := repo.GetSentByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("QueryError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnError(expectedErr)

		_, err := repo.GetSentByUserID(context.Background(), userID)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username", "sum"}).
			AddRow(nil, 100)

		mock.ExpectQuery(`SELECT u1.username, SUM\(t.amount\).*`).
			WithArgs(userID).
			WillReturnRows(rows)

		_, err := repo.GetSentByUserID(context.Background(), userID)

		assert.ErrorContains(t, err, "converting NULL to string")
	})
}
