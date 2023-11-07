package old

import (
	"vmq-go/middleware"
	"vmq-go/router/api"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter(route *gin.Engine) {
	routeGroup := route.Group("/")
	routeGroup.Use(middleware.JSONMiddleware())
	// 心跳
	routeGroup.GET("/appHeart", api.HeartHandler)
	// 收到推送
	routeGroup.GET("/appPush", api.AppPushHandler)
}
