package cache

import (
	"customClothing/src/config"
	"github.com/go-redis/redis"
)

var cache *redis.Client

func Cache() *redis.Client {
	if cache != nil {
		return cache
	}

	InitCache()
	return cache
}

func InitCache() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Cfg().RedisCfg.Addr,
		Password: "", // 密码
		DB:       0,  // 数据库
		PoolSize: 20, // 连接池大小
	})

	if err := rdb.Ping().Err(); err != nil {
		panic("初始化redis失败:" + err.Error())
	}

	cache = rdb
}
