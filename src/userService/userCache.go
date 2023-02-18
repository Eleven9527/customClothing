package userService

import (
	"context"
	"customClothing/src/cache"
	"customClothing/src/config"
	errors "customClothing/src/error"
	"github.com/go-redis/redis"
	"time"
)

type CacheService interface {
	SetToken(ctx context.Context, req *SetTokenReq) errors.Error
}

type cacheSvc struct {
	cache *redis.Client
}

func MakeCacheService() CacheService {
	return &cacheSvc{
		cache: cache.Cache(),
	}
}

func (c *cacheSvc) SetToken(ctx context.Context, req *SetTokenReq) errors.Error {
	cacheKey := config.Cfg().TokenCfg.CacheKey + req.Phone //token:phone

	err := c.cache.Set(cacheKey, req.Token, time.Duration(config.Cfg().TokenCfg.Timeout)).Err()
	if err != nil {
		return errors.New(errors.INTERNAL_ERROR, "")
	}

	return nil
}
