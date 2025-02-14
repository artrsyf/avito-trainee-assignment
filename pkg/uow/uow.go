package uow

import (
	"context"
	"database/sql"
)

//go:generate mockgen -source=uow.go -destination=mock_uow/uow_mock.go -package=mock_uow MockUnitOfWork
type UnitOfWorkI interface {
	Begin(ctx context.Context) error
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}
