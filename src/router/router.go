package router

import (
	"customClothing/src/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	//用户模块
	userRoot := r.Group("/user")
	handler.RegisterUserHandlers(userRoot)

	//订单模块
	orderRoot := r.Group("/order")
	handler.RegisterOrderHandlers(orderRoot)
}
