package postgres

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/entity"
	"github.com/artrsyf/avito-trainee-assignment/internal/user/domain/model"
	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

type UserPostgresRepository struct {
	DB     *sql.DB
	logger *logrus.Logger
}

func NewUserPostgresRepository(
	db *sql.DB,
	logger *logrus.Logger,
) *UserPostgresRepository {
	return &UserPostgresRepository{
		DB:     db,
		logger: logger,
	}
}

func (repo *UserPostgresRepository) Create(
	ctx context.Context,
	user *entity.User,
) (*model.User, error) {
	err := repo.DB.
		QueryRowContext(
			ctx,
			"SELECT 1 FROM users WHERE username = $1",
			user.Username,
		).Scan(new(int))
	if err == nil {
		repo.logger.WithError(err).Error("Trying to create existing user")
		return nil, entity.ErrAlreadyCreated
	}

	if err != sql.ErrNoRows {
		repo.logger.WithError(err).Error("SQL user select error")
		return nil, err
	}

	createdUser := model.User{}
	err = repo.DB.QueryRowContext(
		ctx,
		`INSERT INTO users (username, coins, password_hash) 
		VALUES ($1, $2, $3) 
		RETURNING id, username, coins, password_hash`,
		user.Username, user.Coins, user.PasswordHash,
	).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Coins,
		&createdUser.PasswordHash,
	)
	if err != nil {
		repo.logger.WithError(err).Error("Failed to create user")
		return nil, err
	}

	repo.logger.WithFields(logrus.Fields{
		"user_id": createdUser.ID,
	}).Debug("Created user in Postgres")

	return &createdUser, nil
}

func (repo *UserPostgresRepository) Update(
	ctx context.Context,
	uow uow.Executor,
	user *model.User,
) error {
	_, err := uow.ExecContext(
		ctx,
		"UPDATE users SET coins = $1 WHERE id = $2",
		user.Coins, user.ID,
	)
	if err != nil {
		repo.logger.WithError(err).Error("Failed to update user")
		return err
	}

	repo.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
	}).Debug("Updated user in Postgres")

	return nil
}

func (repo *UserPostgresRepository) GetByID(
	ctx context.Context,
	id uint,
) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRowContext(
			ctx,
			"SELECT id, username, coins, password_hash FROM users WHERE id = $1",
			id,
		).Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
	if err == sql.ErrNoRows {
		repo.logger.WithError(err).Error("Couldn't find such user by id")
		return nil, entity.ErrIsNotExist
	} else if err != nil {
		repo.logger.WithError(err).Error("SQL select user by id error")
		return nil, err
	}

	return &user, nil
}

func (repo *UserPostgresRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {
	user := model.User{}

	err := repo.DB.
		QueryRowContext(
			ctx,
			"SELECT id, username, coins, password_hash FROM users WHERE username = $1",
			username,
		).Scan(&user.ID, &user.Username, &user.Coins, &user.PasswordHash)
	if err == sql.ErrNoRows {
		repo.logger.WithError(err).Error("Couldn't find such user by username")
		return nil, entity.ErrIsNotExist
	} else if err != nil {
		repo.logger.WithError(err).Error("SQL select user by username error")
		return nil, err
	}

	return &user, nil
}
