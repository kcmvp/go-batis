package batis

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"time"
)


type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Ttl(key string, value interface{}, duration time.Duration) error
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
