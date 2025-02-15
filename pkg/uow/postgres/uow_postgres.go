package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type SQLUnitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSQLUnitOfWork(db *sql.DB) *SQLUnitOfWork {
	return &SQLUnitOfWork{db: db}
}

func (u *SQLUnitOfWork) Begin(ctx context.Context) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.tx = tx
	return nil
}

func (u *SQLUnitOfWork) Exec(query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, fmt.Errorf("transaction has not been started")
	}
	return u.tx.Exec(query, args...)
}

func (u *SQLUnitOfWork) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, fmt.Errorf("transaction has not been started")
	}
	return u.tx.ExecContext(ctx, query, args...)
}

func (u *SQLUnitOfWork) Commit() error {
	if u.tx == nil {
		return fmt.Errorf("transaction has not been started")
	}

	return u.tx.Commit()
}

func (u *SQLUnitOfWork) Rollback() error {
	if u.tx == nil {
		return fmt.Errorf("transaction has not been started")
	}

	return u.tx.Rollback()
}
