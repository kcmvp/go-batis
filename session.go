package batis

import (
	"database/sql"
	"fmt"
	"sync"
)


var (
	once     sync.Once
	dbSession       *Session
)


type Session struct {
	*sql.DB
	driverName string
	unsafe     bool
	mapperDir  string
	sqlHookFunc SqlHookFunc
}

func NewSession(cfg Config) (*Session, error) {
	once.Do(func() {
		mapperDir := cfg.validate()
		if db, err := sql.Open(cfg.DriverName, cfg.Url); err != nil {
			panic(fmt.Sprintf("failed to connect to database :%v", err.Error()))
		} else {
			db.SetMaxIdleConns(cfg.MaxIdleConns)
			db.SetMaxOpenConns(cfg.MaxOpenConns)
			dbSession = &Session{DB: db, driverName: cfg.DriverName, mapperDir: mapperDir}
		}
	})
	return dbSession,nil
}

func (s *Session) WithHook(h SqlHookFunc)  {
	s.sqlHookFunc = h
}

// Select using this DB.
func (db *Session) Select(dest interface{}, mapper SqlMapper, args ...interface{}) error {
	panic("")
}

// update, delete
func (db *Session) Exec(mapper SqlMapper, args ...interface{}) error {
	panic("")
}
