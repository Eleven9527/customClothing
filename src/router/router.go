package router

import (
	"customClothing/src/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	//用户模块
	userRoot := r.Group("/user")
	userRoot.POST("", handler.RegisterHandler)    //注册用户
	userRoot.POST("/login", handler.LoginHandler) //登录
}
