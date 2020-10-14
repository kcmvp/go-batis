package udb

import (
	"context"
	"database/sql"
)

// Session is an interface that defines methods for database adapters.
type Session interface {
	// Name returns the name of the database.
	Name() string

	// Ping returns an error if the DBMS could not be reached.
	Ping() error


	// Close terminates the currently active connection to the DBMS and clears
	// all caches.
	Close() error

	Driver() interface{}

	Tx(fn func(sess Session) error) error

	TxContext(ctx context.Context, fn func(sess Session) error, opts *sql.TxOptions) error

	Context() context.Context

	WithContext(ctx context.Context) Session

	//ExecInsert(dao Mapper, arg interface{}) (sql.Result, error)
	//ExecUpdate(dao Mapper, arg interface{}) (sql.Result, error)
	//Select(dest interface{}, dao Mapper, args ...interface{}) error
	//Get(dest interface{}, dao Mapper, args ...interface{}) error

	Exec(arg interface{}) (sql.Result, error)

	Settings
}
