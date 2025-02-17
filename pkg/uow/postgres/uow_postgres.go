package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/artrsyf/avito-trainee-assignment/pkg/uow"
)

type UnitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

func NewFactory(db *sql.DB) uow.Factory {
	return &factory{db: db}
}

type factory struct {
	db *sql.DB
}

func (f *factory) NewUnitOfWork() uow.UnitOfWork {
	return &UnitOfWork{db: f.db}
}

func (u *UnitOfWork) Begin(ctx context.Context) error {
	if u.tx != nil {
		return fmt.Errorf("transaction already started")
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.tx = tx
	return nil
}

func (u *UnitOfWork) Commit() error {
	defer func() { u.tx = nil }()
	if u.tx == nil {
		return fmt.Errorf("transaction not started")
	}
	return u.tx.Commit()
}

func (u *UnitOfWork) Rollback() error {
	defer func() { u.tx = nil }()
	if u.tx == nil {
		return fmt.Errorf("transaction not started")
	}
	return u.tx.Rollback()
}

// Реализация методов Executor
func (u *UnitOfWork) Exec(query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, fmt.Errorf("transaction not started")
	}
	return u.tx.Exec(query, args...)
}

func (u *UnitOfWork) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if u.tx == nil {
		return nil, fmt.Errorf("transaction not started")
	}
	return u.tx.ExecContext(ctx, query, args...)
}

func (u *UnitOfWork) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if u.tx == nil {
		return nil
	}
	return u.tx.QueryRowContext(ctx, query, args...)
}
