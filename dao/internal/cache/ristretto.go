package cache

import (
	"github.com/dgraph-io/ristretto"
	"time"
)




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

func initRistrettoCache() Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Metrics: true,
	})
	if err != nil {
		panic(err)
	}

	return &ristrettoCache{
		ristCache: cache,
	}

}
