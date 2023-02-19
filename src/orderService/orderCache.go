package orderService

import (
	"customClothing/src/cache"
	"github.com/go-redis/redis"
)

type CacheService interface {
}

type cacheSvc struct {
	cache *redis.Client
}

func MakeCacheService() CacheService {
	return &cacheSvc{
		cache: cache.Cache(),
	}
}
