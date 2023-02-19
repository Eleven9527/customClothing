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

	//初始化配置
	config.Mode = config.DevFlag
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

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}
