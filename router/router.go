package router

import (
	"vmq-go/middleware"
	"vmq-go/router/api"
	"vmq-go/router/frontend"
	"vmq-go/router/old"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
	}))
	frontend.SetupFrontendRoutes(router)
	old.InitRouter(router)
	api.InitRouter(router)
	return router
}
