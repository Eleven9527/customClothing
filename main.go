package main

import (
	"customClothing/src/cache"
	"customClothing/src/config"
	"customClothing/src/db"
	"customClothing/src/orderService"
	"customClothing/src/router"
	"customClothing/src/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	//跨域
	r.Use(Cors())

	//初始化配置
	config.Mode = config.PrdFlag
	config.InitConfig()

	//连接数据库
	db.InitDb()

	//初始化缓存
	cache.InitCache()

	//初始化数据表
	userService.InitUserServiceDb()
	orderService.InitOrderServiceDb()

	//注册路由
	router.RegisterRoutes(r)

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
