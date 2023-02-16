package main

import (
	"customClothing/src/config"
	"customClothing/src/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	env string = config.DevFlag //dev=开发环境  prd=生产环境
)

func main() {
	r := gin.Default()

	//初始化配置
	config.InitConfig(env)

	//连接数据库
	db.InitDb()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}
