package batis

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"os"
	"path/filepath"
	"time"
)

var innerCache = defaultCache()

type Cache interface {
	Put(key string, value interface{}) (interface{}, error)
	Get(key string) (interface{}, error)
	Evict(key string) (interface{}, error)
	Touch(key string, duration time.Duration) (interface{}, error)
}

type DbConfig struct {
	DriverName            string `mapstructure:"driverName"`
	Url                   string `mapstructure:"url"`
	MaxOpenConns          int    `mapstructure:"maxOpenConns"`
	MaxIdleConns          int    `mapstructure:"maxIdleConns"`
	MaxTransactionRetries int    `mapstructure:"maxTransactionRetries"`
	MapperDir             string `mapstructure:"mapperDir"`
}

type Config struct {
	DbConfig
	CacheStore Cache
}

var ConfigDefault = Config{
	DbConfig{
		MaxOpenConns:          200,
		MaxIdleConns:          50,
		MaxTransactionRetries: 3,
	},
	innerCache,
}

func (cfg Config) validate() string {
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

func (r ristrettoCache) Put(key string, value interface{}) (interface{}, error) {
	panic("implement me")
}

func (r ristrettoCache) Get(key string) (interface{}, error) {
	panic("implement me")
}

func (r ristrettoCache) Evict(key string) (interface{}, error) {
	panic("implement me")
}

func (r ristrettoCache) Touch(key string, duration time.Duration) (interface{}, error) {
	panic("implement me")
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

