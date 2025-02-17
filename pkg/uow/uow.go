// uow/uow.go
package uow

import (
	"context"
	"database/sql"
)

//go:generate mockgen -source=uow.go -destination=mock_uow/uow_mock.go -package=mock_uow UnitOfWork,Factory
type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
	Executor
}

type Executor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Factory interface {
	NewUnitOfWork() UnitOfWork
}
