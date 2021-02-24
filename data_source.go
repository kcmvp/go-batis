package batis

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var dsMap sync.Map

const MaxOpenConns = 100
const MaxIdleConns = 50

type SqlHookFunc func(ctx context.Context, clause *Clause, ) (string, error)

type DbConfig struct {
	DriverName            string `mapstructure:"driverName"`
	Url                   string `mapstructure:"url"`
	MaxOpenConns          int    `mapstructure:"maxOpenConns"`
	MaxIdleConns          int    `mapstructure:"maxIdleConns"`
	MaxTransactionRetries int    `mapstructure:"maxTransactionRetries"`
	// Mapper dir name
	MapperDir             string `mapstructure:"mapperDir"`
}

type DataSource struct {
	*sql.DB
	cfg      *DbConfig
	cache    Cache
	sqlHooks []SqlHookFunc
}

func NewDsDefaultCache(cfg *DbConfig) (ds *DataSource, err error) {
	return NewDs(cfg, defaultCache())
}

func NewDs(cfg *DbConfig, cache Cache) (*DataSource, error) {
	var ds *DataSource
	var err error
	if _, err = os.Stat(cfg.MapperDir); err != nil {
		return nil,err
	} else {
		cfg.MapperDir,_ =filepath.Abs(cfg.MapperDir)
		fmt.Sprintf("**** Mapper dir is at %v", cfg.MapperDir)
	}
	if cachedDS, ok := dsMap.Load(cfg.Url); !ok {
		var db *sql.DB
		if db, err = sql.Open(cfg.DriverName, cfg.Url); err != nil {
			panic(fmt.Sprintf("failed to connect to database :%v", err.Error()))
			return nil, err
		} else {
			db.SetMaxOpenConns(MaxOpenConns)
			db.SetMaxOpenConns(MaxIdleConns)
			if cfg.MaxIdleConns > 0 {
				db.SetMaxIdleConns(cfg.MaxIdleConns)
			} else if cfg.MaxOpenConns > 0 {
				v := math.Max(float64(cfg.MaxOpenConns), float64(cfg.MaxIdleConns+1))
				db.SetMaxOpenConns(int(v))
			}
			ds = &DataSource{cfg: cfg, cache: cache}
			ds.DB = db
		} // init
		dsMap.Store(cfg.Url, ds)
	} else {
		 ds, _ = cachedDS.(*DataSource)
	}
	return ds, err
}

func (ds *DataSource) MapperDir() string {
	return ds.cfg.MapperDir
}

func (ds *DataSource) Set(key string, value interface{}) error {
	return ds.cache.Set(key, value)
}

func (ds *DataSource) Get(key string) (interface{}, error) {
	return ds.cache.Get(key)
}

func (ds *DataSource) Del(key string) error {
	return ds.cache.Del(key)
}

func (ds *DataSource) Ttl(key string, value interface{}, duration time.Duration) error {
	return ds.cache.Ttl(key, value, duration)
}

func (ds *DataSource) WithSqlHook(h ...SqlHookFunc) {
	ds.sqlHooks = append(ds.sqlHooks, h...)
}

// Select using this DB.
func (ds *DataSource) Query(dest interface{}, mapper SqlMapper, arg interface{}) error {
	return ds.QueryContext(context.Background(), dest, mapper, arg)
}

func (ds DataSource) QueryContext(ctx context.Context, dest interface{}, mapper SqlMapper, arg interface{}) error {
	if clause, err := ds.build(ctx, mapper, arg); err != nil {
		return err
	} else {
		cacheKey := ""
		if cacheKey, _ = clause.CacheKey(); len(cacheKey) > 0 {
			if rt, err := ds.Get(cacheKey); err == nil {
				//@todo unmarshall the json to object
				fmt.Println(fmt.Sprintf("@todo %v", rt))
				return nil
			}
		}
		if rows, err := ds.DB.Query(clause.statement, arg); err == nil {
			//@todo
			fmt.Println(rows)
			//@todo save to cache
			if len(cacheKey) > 0 {
				defer func() {
					ds.Set(cacheKey, nil)
				}()
			}
		}
	}
	return nil
}

// update, delete
func (ds *DataSource) Exec(mapper SqlMapper, args interface{}) error {
	return ds.ExecContext(context.Background(), mapper, args)
}

func (ds *DataSource) ExecContext(ctx context.Context, mapper SqlMapper, arg interface{}) error {
	if clause, err := ds.build(ctx, mapper, arg); err != nil {
		return err
	} else {
		if key, err := clause.CacheKey(); err == nil {
			defer func() {
				ds.Del(key)
			}()
		}
	}
	return nil
}

func (ds *DataSource) build(ctx context.Context, m SqlMapper, args interface{}) (clause *Clause, err error) {
	//@FixMe need to check args type, only support map & struct
	if clause, err = m.build(ds.MapperDir(), args); err == nil {
		//if isSelect != (clause.SqlType() == "select") {
		//	return nil, fmt.Errorf("incorrect statement type:[%v] : %v", clause.SqlType(),.  clause.Statement())
		//}
		for _, hook := range ds.sqlHooks {
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
