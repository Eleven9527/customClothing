package userService

import (
	"context"
	"customClothing/src/cache"
	"customClothing/src/config"
	errors "customClothing/src/error"
	"customClothing/src/utils/token"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type CacheService interface {
	SetToken(ctx context.Context, req *SetTokenReq) errors.Error
	GetToken(tk []byte) (string, errors.Error)
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

//token:phone
func (c *cacheSvc) GetToken(tk []byte) (string, errors.Error) {
	phone, err := token.DecodeToken(tk)
	if err != nil {
		return "", errors.New(errors.TOKEN_VERIFY_ERROR, "token验证失败")
	}

	cacheKey := config.Cfg().TokenCfg.CacheKey + phone //token:phone
	data, err := c.cache.Get(cacheKey).Result()        //
	if err != nil {
		if err == redis.Nil {
			return "", errors.New(errors.USER_NOT_LOGIN, "用户未登录")
		}
		return "", errors.New(errors.INTERNAL_ERROR, "")
	}

	fmt.Println("data = ", data)

	return data, nil
}
