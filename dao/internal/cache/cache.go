package cache

import (
	"time"
)

var defaultCache Cache

func init() {
	var cacheType string = ""
	switch cacheType {
	case "Redis":
		defaultCache = initRedisCache()
	default:
		defaultCache = initRistrettoCache()
	}
}

type Cache interface {
	Put(key string, value interface{}) (interface{}, error)
	Get(key string) (interface{}, error)
	Evict(key string) (interface{}, error)
	Touch(key string, duration time.Duration) (interface{}, error)
}

func Put(key string, value interface{}) (interface{}, error) {
	return defaultCache.Put(key, value)
}

func Get(key string) (interface{}, error) {
	return defaultCache.Get(key)
}

func Evict(key string) (interface{}, error) {
	return defaultCache.Evict(key)
}

func Touch(key string, duration time.Duration) (interface{}, error) {
	return defaultCache.Touch(key, duration)
}
