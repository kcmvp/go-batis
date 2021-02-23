package batis

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

var (
	once      sync.Once
	dbSession *Session
)

type SqlHookFunc func(ctx context.Context, clause *Clause, ) (string, error)

type Session struct {
	*sql.DB
	driverName string
	unsafe     bool
	mapperDir  string
	sqlHooks   []SqlHookFunc
	cache      Cache
}

func (s *Session) Set(key string, value interface{}) error {
	return s.cache.Set(key, value)
}

func (s *Session) Get(key string) (interface{}, error) {
	return s.cache.Get(key)
}

func (s *Session) Del(key string) error {
	return s.cache.Del(key)
}

func (s *Session) Ttl(key string, value interface{}, duration time.Duration) error {
	return s.cache.Ttl(key, value, duration)
}

func NewSession(cfg *DbConfig) (*Session, error) {
	once.Do(func() {
		mapperDir := cfg.validate()
		if db, err := sql.Open(cfg.DriverName, cfg.Url); err != nil {
			panic(fmt.Sprintf("failed to connect to database :%v", err.Error()))
		} else {
			cfg.merge()
			db.SetMaxIdleConns(cfg.MaxIdleConns)
			db.SetMaxOpenConns(cfg.MaxOpenConns)
			dbSession = &Session{DB: db, driverName: cfg.DriverName, mapperDir: mapperDir, cache: cfg.CacheStore}
		}
	})
	return dbSession, nil
}

func (cfg *DbConfig) merge() {
	if cfg.MaxIdleConns < 1 {
		cfg.MaxIdleConns = defaultDbConfig.MaxIdleConns
	}
	if cfg.MaxOpenConns < 1 {
		cfg.MaxOpenConns = defaultDbConfig.MaxOpenConns
	}
	if cfg.MaxOpenConns < cfg.MaxIdleConns {
		cfg.MaxOpenConns = cfg.MaxIdleConns + 1
	}
	if cfg.MaxTransactionRetries < 1 {
		cfg.MaxTransactionRetries = defaultDbConfig.MaxTransactionRetries
	}
	if cfg.CacheStore == nil {
		cfg.CacheStore = defaultCache()
	}
}

func (s *Session) WithSqlHook(h ...SqlHookFunc) {
	s.sqlHooks = append(s.sqlHooks, h...)
}

// Select using this DB.
func (s *Session) Query(dest interface{}, mapper SqlMapper, arg interface{}) error {
	return s.QueryContext(context.Background(), dest, mapper, arg)
}

func (s Session) QueryContext(ctx context.Context, dest interface{}, mapper SqlMapper, arg interface{}) error {
	if clause, err := s.build(ctx, mapper, arg); err != nil {
		return err
	} else {
		cacheKey :=""
		if cacheKey, _ = clause.CacheKey(); len(cacheKey) > 0 {
			if rt, err := s.Get(cacheKey); err == nil {
				//@todo unmarshall the json to object
				fmt.Println(fmt.Sprintf("@todo %v", rt))
				return nil
			}
		}
		if rows, err := s.DB.Query(clause.statement, arg); err == nil {
			//@todo
			fmt.Println(rows)
			//@todo save to cache
			if len(cacheKey) > 0 {
				defer func() {
					s.Set(cacheKey, nil)
				}()
			}
		}
	}
	return nil
}

// update, delete
func (s *Session) Exec(mapper SqlMapper, args interface{}) error {
	return s.ExecContext(context.Background(), mapper, args)
}

func (s Session) ExecContext(ctx context.Context, mapper SqlMapper, arg interface{}) error {
	if clause, err := s.build(ctx, mapper, arg); err != nil {
		return err
	} else {
		if key, err := clause.CacheKey(); err == nil {
			defer func() {
				s.Del(key)
			}()
		}
	}
	return nil
}

func (s Session) build(ctx context.Context, m SqlMapper, args interface{}) (clause *Clause, err error) {
	//@FixMe need to check args type, only support map & struct
	if clause, err = m.build(s.mapperDir, args); err == nil {
		//if isSelect != (clause.SqlType() == "select") {
		//	return nil, fmt.Errorf("incorrect statement type:[%v] : %v", clause.SqlType(),.  clause.Statement())
		//}
		for _, hook := range s.sqlHooks {
			// todo
			if st, err := hook(ctx, clause); err != nil {
				return nil, err
			} else {
				clause.statement = st
			}
		}
	}
	return
}
