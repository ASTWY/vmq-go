package middleware

import (
	"net/http"
	"strings"
	"time"
	"vmq-go/logger"
	"vmq-go/utils/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.Abort()
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]
		if !jwt.IsTokenValid(tokenString) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		logger.GinLogger.Infof("%3d | %13v | %15s | %s  %s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// Response 结构体用于构建API响应
type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

// JSONMiddleware 用于处理JSON格式的API响应
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(http.StatusBadRequest, Response{
				Code:    http.StatusBadRequest,
				Message: c.Errors.Last().Error(),
				Data:    nil,
			})
			c.Abort()
			return
		}

		status := c.Writer.Status()
		message := c.Keys["message"]
		if message == nil {
			switch {
			case status >= http.StatusInternalServerError: // 500
				message = "服务器内部错误"
			case status >= http.StatusBadRequest: // 400
				message = "请求错误"
			case status >= http.StatusUnauthorized: // 401
				message = "未授权"
			case status >= http.StatusNotFound: // 404
				message = "未找到"
			case status >= http.StatusOK: // 200
				message = "Success"
			default:
				message = "未知错误"
			}
		}
		resp := Response{
			Code:    status,
			Message: message,
			Data:    c.Keys["data"],
		}
		c.JSON(status, resp)
	}
}
