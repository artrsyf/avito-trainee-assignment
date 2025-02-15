package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestUserPostgresRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db, logrus.New())

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT 1 FROM users WHERE username = \\$1").
			WithArgs("testuser").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery("INSERT INTO users .* RETURNING .*").
			WithArgs("testuser", 1000, "hash").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "coins", "password_hash"}).
				AddRow(1, "testuser", 1000, "hash"))

		user, err := repo.Create(context.Background(), &entity.User{
			Username:     "testuser",
			Coins:        1000,
			PasswordHash: "hash",
		})

		assert.NoError(t, err)
		assert.Equal(t, &model.User{
			ID:           1,
			Username:     "testuser",
			Coins:        1000,
			PasswordHash: "hash",
		}, user)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		mock.ExpectQuery("SELECT 1 FROM users WHERE username = \\$1").
			WithArgs("existinguser").
			WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

		_, err := repo.Create(context.Background(), &entity.User{
			Username: "existinguser",
		})

		assert.Equal(t, entity.ErrAlreadyCreated, err)
	})

	t.Run("CheckUserError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery("SELECT 1 FROM users WHERE username = \\$1").
			WithArgs("testuser").
			WillReturnError(expectedErr)

		_, err := repo.Create(context.Background(), &entity.User{
			Username: "testuser",
		})

		assert.Equal(t, expectedErr, err)
	})
}

func TestUserPostgresRepository_Update(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db, logrus.New())
	mockUOW := &MockUnitOfWork{}

	t.Run("Success", func(t *testing.T) {
		mockUOW.ExecContextFn = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			assert.Equal(t, "UPDATE users SET coins = $1 WHERE id = $2", query)
			assert.Equal(t, uint(500), args[0].(uint))
			assert.Equal(t, uint(1), args[1].(uint))
			return sqlmock.NewResult(0, 1), nil
		}

		err := repo.Update(context.Background(), mockUOW, &model.User{
			ID:    1,
			Coins: 500,
		})

		assert.NoError(t, err)
	})

	t.Run("UpdateError", func(t *testing.T) {
		expectedErr := sql.ErrTxDone
		mockUOW.ExecContextFn = func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return nil, expectedErr
		}

		err := repo.Update(context.Background(), mockUOW, &model.User{
			ID: 1,
		})

		assert.Equal(t, expectedErr, err)
	})
}

func TestUserPostgresRepository_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db, logrus.New())

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM users WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "coins", "password_hash"}).
				AddRow(1, "testuser", 1000, "hash"))

		user, err := repo.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, &model.User{
			ID:           1,
			Username:     "testuser",
			Coins:        1000,
			PasswordHash: "hash",
		}, user)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM users WHERE id = \\$1").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetByID(context.Background(), 2)

		assert.Equal(t, entity.ErrIsNotExist, err)
	})

	t.Run("QueryError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery("SELECT .* FROM users WHERE id = \\$1").
			WithArgs(3).
			WillReturnError(expectedErr)

		_, err := repo.GetByID(context.Background(), 3)

		assert.Equal(t, expectedErr, err)
	})
}

func TestUserPostgresRepository_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserPostgresRepository(db, logrus.New())

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM users WHERE username = \\$1").
			WithArgs("testuser").
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "coins", "password_hash"}).
				AddRow(1, "testuser", 1000, "hash"))

		user, err := repo.GetByUsername(context.Background(), "testuser")

		assert.NoError(t, err)
		assert.Equal(t, &model.User{
			ID:           1,
			Username:     "testuser",
			Coins:        1000,
			PasswordHash: "hash",
		}, user)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM users WHERE username = \\$1").
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetByUsername(context.Background(), "nonexistent")

		assert.Equal(t, entity.ErrIsNotExist, err)
	})

	t.Run("QueryError", func(t *testing.T) {
		expectedErr := sql.ErrConnDone
		mock.ExpectQuery("SELECT .* FROM users WHERE username = \\$1").
			WithArgs("erroruser").
			WillReturnError(expectedErr)

		_, err := repo.GetByUsername(context.Background(), "erroruser")

		assert.Equal(t, expectedErr, err)
	})
}

type MockUnitOfWork struct {
	uow.UnitOfWorkI
	ExecContextFn func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (m *MockUnitOfWork) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.ExecContextFn(ctx, query, args...)
}
