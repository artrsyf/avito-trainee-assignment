package postgres

import (
	"database/sql"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
)

type UserPostgresRepository struct {
	DB *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		DB: db,
	}
}

func (repo *UserPostgresRepository) Create(user *entity.User) (*model.User, error) {
	err := repo.DB.
		QueryRow("SELECT 1 FROM users WHERE username = $1", user.Username).Scan(new(int))
	if err == nil {
		return nil, entity.ErrAlreadyCreated
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	createdUser := model.User{}
	err = repo.DB.QueryRow("INSERT INTO users (username, coins, password_hash) VALUES ($1, $2, $3) RETURNING id, username, coins, password_hash", user.Username, user.Coins, user.PasswordHash).
		Scan(&createdUser.ID, &createdUser.Username, &createdUser.Coins, &createdUser.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (repo *UserPostgresRepository) GetById(id uint) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRow("SELECT id, username, coins, password_hash FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, entity.ErrIsNotExist
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserPostgresRepository) GetByUsername(username string) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRow("SELECT id, username, coins, password_hash FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, entity.ErrIsNotExist
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
