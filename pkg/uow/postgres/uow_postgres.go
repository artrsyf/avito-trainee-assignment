package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresUnitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

func NewPostgresUnitOfWork(db *sql.DB) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{db: db}
}

func (u *PostgresUnitOfWork) Begin(ctx context.Context) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.tx = tx
	return nil
}

func (u *PostgresUnitOfWork) Exec(query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, fmt.Errorf("transaction has not been started")
	}
	return u.tx.Exec(query, args...)
}

func (u *PostgresUnitOfWork) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, fmt.Errorf("transaction has not been started")
	}
	return u.tx.ExecContext(ctx, query, args...)
}

func (u *PostgresUnitOfWork) Commit() error {
	if u.tx == nil {
		return fmt.Errorf("transaction has not been started")
	}

	return u.tx.Commit()
}

func (u *PostgresUnitOfWork) Rollback() error {
	if u.tx == nil {
		return fmt.Errorf("transaction has not been started")
	}

	return u.tx.Rollback()
}
