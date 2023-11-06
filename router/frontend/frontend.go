package frontend

import (
	"os"
	"vmq-go/logger"
	"vmq-go/task"

	"github.com/gin-gonic/gin"
)

func SetupFrontendRoutes(router *gin.Engine) {
	// 检查web目录是否存在
	if _, err := os.Stat("./web"); os.IsNotExist(err) {
		logger.Logger.Info("web目录不存在，下载前端文件")
		// 下载前端文件
		task.DownloadFrontend()
	}
	// 挂载静态资源
	router.Static("/assets", "./web/assets")
	// 挂载 favicon.ico
	router.StaticFile("/favicon.ico", "./web/favicon.ico")
	// 返回 index.html
	router.StaticFile("/", "./web/index.html")
}
