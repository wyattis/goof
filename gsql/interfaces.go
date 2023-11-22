package gsql

import (
	"context"
	"database/sql"
)

// SQL Database interfaces
type ISelect interface {
	Select(dest any, query string, params ...interface{}) error
	Get(dest any, query string, params ...interface{}) error
}

type ISelectContext interface {
	SelectContext(ctx context.Context, dest any, query string, params ...interface{}) error
	GetContext(ctx context.Context, dest any, query string, params ...interface{}) error
}

type IQuery interface {
	Query(sql string, params ...interface{}) (*sql.Rows, error)
}

type IQueryContext interface {
	QueryContext(ctx context.Context, sql string, params ...interface{}) (*sql.Rows, error)
}

type IQueryRow interface {
	QueryRow(sql string, params ...interface{}) *sql.Row
}

type IQueryRowContext interface {
	QueryRowContext(ctx context.Context, sql string, params ...interface{}) *sql.Row
}

type IExec interface {
	Exec(sql string, params ...interface{}) (sql.Result, error)
}

type IExecContext interface {
	ExecContext(ctx context.Context, sql string, params ...interface{}) (sql.Result, error)
}

type INamedExec interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type INamedExecContext interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type IBegin interface {
	Begin() (*sql.Tx, error)
}

type IBeginTx interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type IDB interface {
	IExec
	IQuery
	IQueryRow
	IBegin
}

type IDBContext interface {
	IExecContext
	IQueryContext
	IQueryRowContext
	IBeginTx
}
