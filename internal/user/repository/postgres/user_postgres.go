package postgres

import (
	"context"
	"database/sql"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

type UserPostgresRepository struct {
	DB *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		DB: db,
	}
}

func (repo *UserPostgresRepository) Create(ctx context.Context, user *entity.User) (*model.User, error) {
	err := repo.DB.
		QueryRowContext(ctx, "SELECT 1 FROM users WHERE username = $1", user.Username).Scan(new(int))
	if err == nil {
		return nil, entity.ErrAlreadyCreated
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	createdUser := model.User{}
	err = repo.DB.QueryRowContext(ctx, "INSERT INTO users (username, coins, password_hash) VALUES ($1, $2, $3) RETURNING id, username, coins, password_hash", user.Username, user.Coins, user.PasswordHash).
		Scan(&createdUser.ID, &createdUser.Username, &createdUser.Coins, &createdUser.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (repo *UserPostgresRepository) Update(ctx context.Context, uow uow.UnitOfWorkI, user *model.User) error {
	_, err := uow.ExecContext(ctx, "UPDATE users SET coins = $1 WHERE id = $2", user.Coins, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserPostgresRepository) GetById(ctx context.Context, id uint) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRowContext(ctx, "SELECT id, username, coins, password_hash FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, entity.ErrIsNotExist
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserPostgresRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRowContext(ctx, "SELECT id, username, coins, password_hash FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, entity.ErrIsNotExist
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
