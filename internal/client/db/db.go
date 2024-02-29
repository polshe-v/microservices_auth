package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Handler is a function executed in transaction.
type Handler func(ctx context.Context) error

// Client interface is a client for DB.
type Client interface {
	DB() DB
	Close() error
}

// Query is the wrapper for raw DB query and its name.
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor interface is used for DB transactions.
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// TxManager interface contains handler for transaction.
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// SQLExecutor interface gathers all DB query executors.
type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

// NamedExecutor interface for named queries which use tags.
type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecutor interface for common DB queries.
type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
}

// Pinger interface for checking DB connectivity.
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB interface for communication with DB.
type DB interface {
	SQLExecutor
	Transactor
	Pinger
	Close()
}
