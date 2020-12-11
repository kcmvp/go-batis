package batis

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"os"
	"path/filepath"
	"time"
)

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Ttl(key string, value interface{}, duration time.Duration) error
}

type DbConfig struct {
	DriverName            string `mapstructure:"driverName"`
	Url                   string `mapstructure:"url"`
	MaxOpenConns          int    `mapstructure:"maxOpenConns"`
	MaxIdleConns          int    `mapstructure:"maxIdleConns"`
	MaxTransactionRetries int    `mapstructure:"maxTransactionRetries"`
	MapperDir             string `mapstructure:"mapperDir"`
	CacheStore            Cache
}

var defaultDbConfig = DbConfig{
	MaxOpenConns:          200,
	MaxIdleConns:          50,
	MaxTransactionRetries: 3,
}

func (cfg DbConfig) validate() string {
	if ex, err := os.Executable(); err != nil {
		panic(err)
	} else {
		exPath := filepath.Dir(ex)
		path := filepath.Join(exPath, cfg.MapperDir)
		if _, err := os.Stat(path); err != nil {
			panic(fmt.Sprintf("Can not find the mapper dir %v", path))
		} else {
			fmt.Sprintf("**** Mapper dir is at %v", path)
		}
		return path
	}
}

type ristrettoCache struct {
	ristCache *ristretto.Cache
}

func (r ristrettoCache) Set(key string, value interface{}) error {
	rs := r.ristCache.Set(key, value, 0)
	if rs {
		return nil
	} else {
		return errors.New("failed to save the value")
	}
}

func (r ristrettoCache) Get(key string) (interface{}, error) {
	v, b := r.ristCache.Get(key)
	if b {
		return v, nil
	} else {
		return nil, errors.New(fmt.Sprintf("failed to get the value of the key %v", key))
	}
}

func (r ristrettoCache) Del(key string) error {
	r.ristCache.Del(key)
	return nil
}

func (r ristrettoCache) Ttl(key string, value interface{}, duration time.Duration) error {

	if r.ristCache.SetWithTTL(key, value, 0, duration) {
		return nil
	} else {
		return errors.New(fmt.Sprintf("failed to set the key with ttl: %v", key))
	}
}

func defaultCache() Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Metrics:     true,
	})
	if err != nil {
		panic(err)
	}

	return &ristrettoCache{
		ristCache: cache,
	}

}
