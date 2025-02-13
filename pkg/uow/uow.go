package uow

import (
	"context"
	"database/sql"
)

type UnitOfWorkI interface {
	Begin(ctx context.Context) error
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}
