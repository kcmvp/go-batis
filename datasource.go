package batis

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

// DataSource is an interface that defines methods for database adapters.
type DataSource interface {
	// Name returns the name of the database.
	Name() string

	// Ping returns an error if the DBMS could not be reached.
	Ping() error


	// Close terminates the currently active connection to the DBMS and clears
	// all caches.
	Close() error

	Driver() interface{}

	Tx(fn func(sess DataSource) error) error

	TxContext(ctx context.Context, fn func(sess DataSource) error, opts *sql.TxOptions) error

	Context() context.Context

	WithContext(ctx context.Context) DataSource

	//ExecInsert(dao Mapper, arg interface{}) (sql.Result, error)
	//ExecUpdate(dao Mapper, arg interface{}) (sql.Result, error)
	//Select(dest interface{}, dao Mapper, args ...interface{}) error
	//Get(dest interface{}, dao Mapper, args ...interface{}) error

	Exec(sql string, arg interface{}) (sql.Result, error)

	Settings
}
