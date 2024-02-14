package admin

import (
	"net/http"
	"vmq-go/db"
	"vmq-go/middleware"
	"vmq-go/utils/captcha"
	"vmq-go/utils/jwt"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/login", loginHandler)                              // 登录
		adminGroup.Use(middleware.AuthMiddleware())                          // 验证token 中间件
		adminGroup.GET("/logout", logoutHandler)                             // 登出
		adminGroup.GET("/settings", getSettingsHandler)                      // 获取所有设置
		adminGroup.PUT("/setting", putSettingHandler)                        // 更新设置
		adminGroup.POST("/setting/email", sendEmailHandler)                  // 发送测试邮件
		adminGroup.GET("/orders", getOrderHandler)                           // 获取订单列表
		adminGroup.GET("/order/today", getOrderDataTodayHandler)             // 获取今日订单
		adminGroup.DELETE("/order/:orderId", deleteOrderHandler)             // 删除订单
		adminGroup.POST("/order/:orderId/replenish", reCallbackOrderHandler) // 重新回调订单
		adminGroup.GET("/qrcodes", getQrcodeHandler)                         // 获取二维码列表
		adminGroup.POST("/qrcode", postQrcodeHandler)                        // 创建二维码
		adminGroup.DELETE("/qrcode/:id", deleteQrcodeHandler)                // 删除二维码
		adminGroup.GET("/paylogs", getPaylogHandler)                         // 获取支付日志列表
	}
}

// login model
type login struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaID string `json:"captcha_id"`
}

// login
func loginHandler(c *gin.Context) {
	var login login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if login.Username == "" || login.Password == "" || login.Captcha == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		c.Set("message", "用户名、密码、验证码不能为空")
		return
	}
	if !captcha.VerifyCaptcha(login.CaptchaID, login.Captcha) {
		c.Set("message", "验证码错误")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	appConfig, err := db.GetAppConfig()
	if !appConfig.VerifyAdmin(login.Username, login.Password) {
		c.AbortWithStatus(http.StatusBadRequest)
		c.Set("message", "用户名或密码错误")
		return
	}
	token, err := jwt.GenerateToken(login.Username)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		c.Set("message", "token生成失败")
		return
	}
	c.Set("data", gin.H{"token": token})
}

// logoutHandler
func logoutHandler(c *gin.Context) {
}
