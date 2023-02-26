package router

import (
	_ "customClothing/docs"
	"customClothing/src/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	//文档
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//用户模块
	userRoot := r.Group("/user")
	handler.RegisterUserHandlers(userRoot)

	//订单模块
	orderRoot := r.Group("/order")
	handler.RegisterOrderHandlers(orderRoot)
}
