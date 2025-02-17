package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/purchase/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

func TestPurchasePostgresRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPurchasePostgresRepository(db, logrus.New())
	mockUOW := &MockUnitOfWork{db: db}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM purchase_types WHERE name = \\$1").
			WithArgs("t-shirt").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("INSERT INTO purchases .* RETURNING .*").
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "purchaser_id", "purchase_type_id"}).
				AddRow(1, 1, 1))

		purchase, err := repo.Create(context.Background(), mockUOW, &entity.Purchase{
			PurchaserID:      1,
			PurchaseTypeName: "t-shirt",
		})

		assert.NoError(t, err)
		assert.Equal(t, &model.Purchase{
			ID:             1,
			PurchaserID:    1,
			PurchaseTypeID: 1,
		}, purchase)
	})

	t.Run("PurchaseTypeNotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM purchase_types WHERE name = \\$1").
			WithArgs("invalid-type").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.Create(context.Background(), mockUOW, &entity.Purchase{
			PurchaserID:      1,
			PurchaseTypeName: "invalid-type",
		})

		assert.ErrorIs(t, err, entity.ErrNotExistedProduct)
	})

	t.Run("InsertError", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM purchase_types WHERE name = \\$1").
			WithArgs("t-shirt").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		expectedErr := sql.ErrConnDone
		mock.ExpectQuery("INSERT INTO purchases .* RETURNING .*").
			WithArgs(1, 1).
			WillReturnError(expectedErr)

		_, err := repo.Create(context.Background(), mockUOW, &entity.Purchase{
			PurchaserID:      1,
			PurchaseTypeName: "t-shirt",
		})

		assert.Equal(t, expectedErr, err)
	})
}

func TestPurchasePostgresRepository_GetProductByType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPurchasePostgresRepository(db, logrus.New())

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, cost FROM purchase_types WHERE name = \\$1").
			WithArgs("t-shirt").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "cost"}).
				AddRow(1, "t-shirt", 80))

		pt, err := repo.GetProductByType(context.Background(), "t-shirt")

		assert.NoError(t, err)
		assert.Equal(t, &model.PurchaseType{
			ID:   1,
			Name: "t-shirt",
			Cost: 80,
		}, pt)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, cost FROM purchase_types WHERE name = \\$1").
			WithArgs("invalid-type").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetProductByType(context.Background(), "invalid-type")

		assert.ErrorIs(t, err, entity.ErrNotExistedProduct)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery("SELECT id, name, cost FROM purchase_types WHERE name = \\$1").
			WithArgs("t-shirt").
			WillReturnError(expectedErr)

		_, err := repo.GetProductByType(context.Background(), "t-shirt")

		assert.Equal(t, expectedErr, err)
	})
}

func TestPurchasePostgresRepository_GetPurchasesByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPurchasePostgresRepository(db, logrus.New())
	userID := uint(1)

	t.Run("SuccessWithData", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "quantity"}).
			AddRow("t-shirt", 3).
			AddRow("cup", 5)

		mock.ExpectQuery(`SELECT pt.name, COUNT\(p.id\) as quantity .*`).
			WithArgs(userID).
			WillReturnRows(rows)

		result, err := repo.GetPurchasesByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, entity.Inventory{
			{PurchaseTypeName: "t-shirt", Quantity: 3},
			{PurchaseTypeName: "cup", Quantity: 5},
		}, result)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "quantity"})

		mock.ExpectQuery(`SELECT pt.name, COUNT\(p.id\) as quantity .*`).
			WithArgs(userID).
			WillReturnRows(rows)

		result, err := repo.GetPurchasesByUserID(context.Background(), userID)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("QueryError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery(`SELECT pt.name, COUNT\(p.id\) as quantity .*`).
			WithArgs(userID).
			WillReturnError(expectedErr)

		_, err := repo.GetPurchasesByUserID(context.Background(), userID)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("ScanError", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name", "quantity"}).
			AddRow(nil, 5)

		mock.ExpectQuery(`SELECT pt.name, COUNT\(p.id\) as quantity .*`).
			WithArgs(userID).
			WillReturnRows(rows)

		_, err := repo.GetPurchasesByUserID(context.Background(), userID)

		assert.ErrorContains(t, err, "converting NULL to string")
	})
}

type MockUnitOfWork struct {
	uow.Executor
	db *sql.DB
}

func (m *MockUnitOfWork) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *MockUnitOfWork) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *MockUnitOfWork) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return m.db.QueryRowContext(ctx, query, args...)
}
