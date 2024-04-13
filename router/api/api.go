package api

import (
	"encoding/base64"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"vmq-go/db"
	"vmq-go/middleware"
	"vmq-go/router/api/admin"
	"vmq-go/task"
	"vmq-go/utils"
	"vmq-go/utils/captcha"
	"vmq-go/utils/hash"
	"vmq-go/utils/qrcode"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter(route *gin.Engine) {
	routeGroup := route.Group("/api")
	// 创建订单
	routeGroup.POST("/order", creatOrderHandler)
	// qrcode
	routeGroup.GET("/qrcode", qrcodeGetHandler)
	routeGroup.Use(middleware.JSONMiddleware())
	admin.SetupAdminRoutes(routeGroup)
	// 解析二维码
	routeGroup.POST("/qrcode", qrcodePostHandler)
	// 查询订单状态
	routeGroup.GET("/order/:orderId/state", getOrderStateGetHandler)
	// 验证码
	routeGroup.GET("/captcha", captchaHandler)
	// 心跳
	routeGroup.GET("/appHeart", HeartHandler)
	// 收到推送
	routeGroup.GET("/appPush", AppPushHandler)
	// 查询订单
	routeGroup.GET("/order/:orderId", getOrderGetHandler)
	routeGroup.Use(middleware.AuthMiddleware())
	// 重新回调订单
	routeGroup.PUT("/order/:orderId", reCallbackOrderHandler)
}

func qrcodeGetHandler(c *gin.Context) {
	content := c.Query("content")
	format := c.DefaultQuery("format", "json")
	if content == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "content is empty",
		})
		return
	}
	base64Str, err := qrcode.QrcodeFromStr(content)
	if err != nil {
		c.JSON(
			200,
			gin.H{
				"code": -1,
				"msg":  err.Error(),
			},
		)
		return
	}
	switch format {
	case "image":
		c.Writer.Header().Set("Content-Type", "image/png")
		c.Request.Header.Set("Content-Type", "image/png")
		// 将base64Str 转为 []byte
		buf, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			c.JSON(
				200,
				gin.H{
					"code": -1,
					"msg":  err.Error(),
				},
			)
			return
		}
		c.Writer.Write(buf)
	default:
		c.JSON(200, gin.H{"qrcode": fmt.Sprintf("data:image/png;base64,%s", base64Str)})
	}
}

func qrcodePostHandler(c *gin.Context) {
	// 从请求中获取图片
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
		return
	}
	// 读取file 2 []byte
	src, err := file.Open()
	if err != nil {
		c.Error(err)
		return
	}
	defer src.Close()
	buf := make([]byte, file.Size)
	_, err = src.Read(buf)
	if err != nil {
		c.Set("code", http.StatusInternalServerError)
		c.Error(err)
		return
	}
	// []byte 2 base64
	base64Str := base64.StdEncoding.EncodeToString(buf)
	// 解读二维码
	content, err := qrcode.DecodeQrcodeFromStr(base64Str)
	if err != nil {
		c.Set("code", http.StatusInternalServerError)
		c.Error(err)
		return
	}
	c.Set("data", gin.H{"content": content})
}

type CreateOrderParams struct {
	PayId     string  `json:"payId" binding:"required"`
	Type      int     `json:"type" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Sign      string  `json:"sign" binding:"required"`
	Param     string  `json:"param"`
	NotifyUrl string  `json:"notifyUrl"`
	ReturnUrl string  `json:"returnUrl"`
}

func creatOrderHandler(c *gin.Context) {
	// 检查订单是否过期
	task.CheckOrderExpire()
	// 检查心跳
	heart := task.CheckHeart()
	if !heart {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "heart error",
		})
		return
	}
	// 获取参数
	var params CreateOrderParams
	if c.ContentType() == "application/x-www-form-urlencoded" {
		payId := c.PostForm("payId")
		typeStr := c.PostForm("type")
		priceStr := c.PostForm("price")
		signstr := c.PostForm("sign")
		param := c.PostForm("param")
		notifyUrl := c.PostForm("notifyUrl")
		returnUrl := c.PostForm("returnUrl")
		if payId == "" || typeStr == "" || priceStr == "" || signstr == "" {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "param error",
			})
			return
		}
		typeInt, err := strconv.Atoi(typeStr)
		if err != nil {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "type error",
			})
			return
		}
		priceFloat, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "price error",
			})
			return
		}
		params = CreateOrderParams{
			PayId:     payId,
			Type:      typeInt,
			Price:     priceFloat,
			Sign:      signstr,
			Param:     param,
			NotifyUrl: notifyUrl,
			ReturnUrl: returnUrl,
		}
		appConfig, err := db.GetAppConfig()
		if err != nil {
			c.Error(err)
			return
		}
		// 1. 验证签名
		sign := hash.GetMD5Hash(payId + param + typeStr + priceStr + appConfig.APISecret)
		if sign != params.Sign {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "sign error",
			})
			return
		}
	} else {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "content-type error",
		})
		return
	}
	// 创建订单
	// 2. 验证订单是否存在
	_, err := db.GetPayOrderByPayID(params.PayId)
	if err == nil || err.Error() != "record not found" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "payId is exist",
		})
		return
	}
	err = nil
	// 3. 创建订单
	err = db.AddPayOrder(params.PayId, params.Type, params.Price, params.Param, params.NotifyUrl, params.ReturnUrl)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	// 返回结果
	order, err := db.GetPayOrderByPayID(params.PayId)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	timeout := (order.ExpectDate - order.CreateDate) / 60
	c.IndentedJSON(200, gin.H{
		"code": 1,
		"msg":  "success",
		"data": gin.H{
			"payId":       order.PayID,
			"orderId":     order.OrderID,
			"payType":     order.Type,
			"price":       order.Price,
			"reallyPrice": order.ReallyPrice,
			"payUrl":      order.PayURL,
			"isAuto":      order.IsAuto,
			"state":       order.State,
			"createDate":  order.CreateDate,
			"expectDate":  order.ExpectDate,
			"timeOut":     timeout,
			"redirectUrl": fmt.Sprintf("/payment/%s", order.OrderID),
		},
	})
}

func getOrderGetHandler(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		c.Error(fmt.Errorf("orderId is empty"))
		return
	}
	order, err := db.GetPayOrderByOrderID(orderId)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", gin.H{
		"payId":       order.PayID,
		"orderId":     order.OrderID,
		"payType":     order.Type,
		"price":       order.Price,
		"reallyPrice": order.ReallyPrice,
		"payUrl":      order.PayURL,
		"isAuto":      order.IsAuto,
		"state":       order.State,
		"createDate":  order.CreateDate,
		"expectDate":  order.ExpectDate,
	})
}

func getOrderStateGetHandler(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		c.Error(fmt.Errorf("orderId is empty"))
		return
	}
	order, err := db.GetPayOrderByOrderID(orderId)
	if err != nil {
		if err.Error() == "record not found" {
			c.Error(fmt.Errorf("order not found"))
		} else {
			c.Error(err)
		}
		return
	}
	paramMap := map[string]string{
		"payId":       order.PayID,
		"param":       order.Param,
		"type":        fmt.Sprintf("%d", order.Type),
		"price":       utils.Float64ToSting(order.Price),
		"reallyPrice": utils.Float64ToSting(order.ReallyPrice),
	}
	appConfig, err := db.GetAppConfig()
	if err != nil {
		c.Error(err)
		return
	}
	// sign := hash.GetMD5Hash(payId + param + typeStr + priceStr + key.VValue)
	sign := hash.GetMD5Hash(fmt.Sprintf("%s%s%s%s%s", order.PayID, order.Param, fmt.Sprintf("%d", order.Type), utils.Float64ToSting(order.Price), utils.Float64ToSting(order.ReallyPrice)) + appConfig.APISecret)
	// 将map转为get参数 用于跳转
	paramStr := ""
	for k, v := range paramMap {
		paramStr += fmt.Sprintf("%s=%s&", k, v)
	}
	paramStr += fmt.Sprintf("sign=%s", sign)
	returnUrl := order.ReturnURL
	if returnUrl == "" {
		returnUrl = appConfig.ReturnUrl
	}
	var state int
	if order.State >= 1 {
		state = 1
		returnUrl = fmt.Sprintf("%s?%s", returnUrl, paramStr)
	} else {
		state = order.State
		returnUrl = ""
	}
	c.Set("data", gin.H{
		"state":     state,
		"returnUrl": returnUrl,
	})
}

func reCallbackOrderHandler(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		c.Error(fmt.Errorf("orderId is empty"))
		return
	}
	order, err := db.GetPayOrderByOrderID(orderId)
	if err != nil {
		c.Error(err)
		return
	}
	if order.State != 1 {
		c.Error(fmt.Errorf("order state error"))
		return
	}
	task.Notify(order)
	c.Set("code", http.StatusOK)
}

func captchaHandler(c *gin.Context) {
	id, b64s, err := captcha.GenerateCaptcha()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", gin.H{
		"id":      id,
		"captcha": b64s,
	})
}

func HeartHandler(c *gin.Context) {
	time := c.Query("t")
	if time == "" {
		c.Error(fmt.Errorf("t is empty"))
		return
	}
	timeInt, err := strconv.ParseInt(time, 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("time error"))
		return
	}
	timeNow := utils.GetUnix13()
	// 如果时间差大于10秒
	if math.Abs(float64(timeNow-timeInt)) > 10000 {
		c.Error(fmt.Errorf("time error"))
		return
	}
	sign := c.Query("sign")
	if sign == "" {
		c.Error(fmt.Errorf("sign is empty"))
		return
	}
	appConfig, err := db.GetAppConfig()
	if err != nil {
		c.Error(err)
		return
	}
	if hash.GetMD5Hash(time+appConfig.APISecret) != sign {
		c.Error(fmt.Errorf("sign error"))
		return
	}
	err = db.UpdateSetting("lastHeart", time)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("code", http.StatusOK)
	c.Set("data", "success")
}

func AppPushHandler(c *gin.Context) {
	t := c.Query("t")
	if t == "" {
		c.Error(fmt.Errorf("t is empty"))
		return
	}
	_type := c.Query("type") // 1:微信 2:支付宝
	if _type == "" {
		c.Error(fmt.Errorf("type is empty"))
		return
	}
	if _type != "1" && _type != "2" {
		c.Error(fmt.Errorf("type error"))
		return
	}
	typeInt, err := strconv.Atoi(_type)
	if err != nil {
		c.Error(err)
		return
	}
	price := c.Query("price")
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.Error(err)
		return
	}
	sign := c.Query("sign")
	if sign == "" {
		c.Error(fmt.Errorf("sign is empty"))
		return
	}
	metdata := c.DefaultQuery("metadata", "")
	appConfig, err := db.GetAppConfig()
	if err != nil {
		c.Error(err)
		return
	}
	if hash.GetMD5Hash(_type+price+t+appConfig.APISecret) != sign {
		c.Error(fmt.Errorf("sign error"))
		return
	}
	go task.AppPush(typeInt, priceFloat, metdata)
	c.Set("code", http.StatusOK)
}
