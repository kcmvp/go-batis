package cache

import (
	. "github.com/go-redis/redis/v8"
	"time"
)

type redisCache struct {
	*ClusterClient
}

func (r redisCache) Put(key string, value interface{}) (interface{}, error) {
	panic("implement me")
}

func (r redisCache) Get(key string) (interface{}, error) {
	panic("implement me")
}

func (r redisCache) Evict(key string) (interface{}, error) {
	panic("implement me")
}

func (r redisCache) Touch(key string, duration time.Duration) (interface{}, error) {
	panic("implement me")
}

func initRedisCache() Cache {

	rdb := NewClusterClient(&ClusterOptions{
		Addrs: []string{},
	})

	/*
		rdb2 := NewFailoverClusterClient(&FailoverOptions{
			MasterName:    "master-name",
			SentinelAddrs: []string{":9126", ":9127", ":9128"},

		})
	*/

	return &redisCache{
		ClusterClient: rdb,
	}
}
