package task

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"vmq-go/db"
	"vmq-go/logger"
	"vmq-go/utils"
	"vmq-go/utils/hash"
)

// sendEmail
// param subject 邮件主题
// param body 邮件内容
func SendEmailTask(subject string, body string) {
	// STMP 服务器地址
	setting, err := db.GetSetting("emailSMTPhost")
	if err != nil {
		logger.Logger.Errorf("获取SMTP服务器地址失败, %s", err)
		return
	}
	smtpHost := setting.VValue
	// SMTP 服务器端口
	setting, err = db.GetSetting("emailSMTPport")
	if err != nil {
		logger.Logger.Errorf("获取SMTP服务器端口失败, %s", err)
		return
	}
	smtpPort, err := strconv.Atoi(setting.VValue)
	if err != nil {
		logger.Logger.Errorf("类型转换失败, %s", err)
		return
	}
	// SMTP 服务器账号
	setting, err = db.GetSetting("emailSMTPuser")
	if err != nil {
		logger.Logger.Errorf("获取SMTP服务器账号失败, %s", err)
		return
	}
	smtpUser := setting.VValue
	// SMTP 服务器密码
	setting, err = db.GetSetting("emailSMTPpwd")
	if err != nil {
		logger.Logger.Errorf("获取SMTP服务器密码失败, %s", err)
		return
	}
	smtpPassword := setting.VValue
	// SMTP 服务器发件人邮箱
	setting, err = db.GetSetting("emailSMTPfrom")
	if err != nil {
		logger.Logger.Errorf("获取SMTP服务器发件人邮箱失败, %s", err)
		return
	}
	smtpSender := setting.VValue
	// 收件人邮箱
	setting, err = db.GetSetting("emailSMTPto")
	if err != nil {
		logger.Logger.Errorf("获取收件人邮箱失败, %s", err)
		return
	}
	email := setting.VValue
	// emailSMTPssl
	setting, err = db.GetSetting("emailSMTPssl")
	if err != nil {
		logger.Logger.Errorf("获取emailSMTPssl失败, %s", err)
		return
	}
	emailSMTPssl, err := strconv.Atoi(setting.VValue)
	if err != nil {
		logger.Logger.Errorf("类型转换失败, %s", err)
		return
	}
	var emailSMTPsslBool bool
	if emailSMTPssl == 0 {
		emailSMTPsslBool = false
	} else {
		emailSMTPsslBool = true
	}
	err = utils.SendEmailUsingSMTP(subject, body, smtpHost, smtpPort, smtpUser, smtpPassword, smtpSender, email, emailSMTPsslBool)
	if err != nil {
		logger.Logger.Errorf("发送邮件失败, %s", err)
		return
	}
}

// 检查订单是否过期
func CheckOrderExpire() {
	var payOrders []db.PayOrder
	err := db.DB.Where("state = ?", 0).Find(&payOrders).Error
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	count := 0
	for _, payOrder := range payOrders {
		if payOrder.ExpectDate < utils.GetUnix10() {
			payOrder.State = -1
			payOrder.CloseDate = utils.GetUnix10()
			if err := db.UpdatePayOrder(payOrder); err != nil {
				logger.Logger.Error(err)
			}
			count++
		}
	}
}

// 检查心跳
func CheckHeart() bool {
	logger.Logger.Info("检查心跳")
	lastHeart, err := db.GetSetting("lastHeart")
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	lastHeartInt, err := strconv.ParseInt(lastHeart.VValue, 10, 64)
	if err != nil {
		logger.Logger.Error(err)
		return false
	}
	if utils.GetUnix13()-lastHeartInt > 30000 {
		logger.Logger.Info("检查心跳完成, 心跳超时")
		setting, err := db.GetSetting("monitorNotice")
		if err != nil {
			logger.Logger.Error(err)
		} else {
			if setting.VValue == "1" {
				timeStr := time.Unix(lastHeartInt/1000, 0).Format("2006-01-02 15:04:05")
				go SendEmailTask("监控端掉线", "监控端已掉线，上次心跳时间："+timeStr)
			}
		}
		return false
	}
	logger.Logger.Info("检查心跳完成, 心跳正常")
	return true
}

func Notify(order db.PayOrder) {
	// payId=1547130349673&param=vone666&type=2&price=0.1&reallyPrice=0.1
	// 获取异步通知地址
	var notifyURL string
	if order.NotifyURL != "" {
		notifyURL = order.NotifyURL
	} else {
		setting, err := db.GetSetting("notifyURL")
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		notifyURL = setting.VValue
	}
	params := map[string]string{
		"payId":       order.PayID,
		"price":       utils.Float64ToSting(order.Price),
		"type":        strconv.Itoa(order.Type),
		"param":       order.Param,
		"reallyPrice": utils.Float64ToSting(order.ReallyPrice),
	}
	apiSecret, err := db.GetSetting("apiSecret")
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	params["sign"] = hash.GetMD5Hash(fmt.Sprintf("%s%s%s%s%s", order.PayID, order.Param, fmt.Sprintf("%d", order.Type), utils.Float64ToSting(order.Price), utils.Float64ToSting(order.ReallyPrice)) + apiSecret.VValue)
	// 发送异步通知 GET  使用net/http
	httpClient := &http.Client{}
	var paramsStr string
	for k, v := range params {
		paramsStr += fmt.Sprintf("%s=%s&", k, v)
	}
	paramsStr = paramsStr[:len(paramsStr)-1]
	req, err := http.NewRequest("GET", notifyURL+"?"+paramsStr, nil)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Logger.Error(err)
		order.State = 2
		order.NotifyRes = string(err.Error())
		if err := db.UpdatePayOrder(order); err != nil {
			logger.Logger.Error(err)
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Logger.Error("异步通知失败, 状态码", resp.StatusCode)
		order.State = 2
		if err := db.UpdatePayOrder(order); err != nil {
			logger.Logger.Error(err)
		}
		setting, err := db.GetSetting("errorNotice")
		if err != nil {
			logger.Logger.Error(err)
		} else {
			if setting.VValue != "0" {
				go SendEmailTask("异步通知失败", fmt.Sprintf("异步通知失败，状态码：%d，订单ID：%d", resp.StatusCode, order.ID))
			}
		}
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		logger.Logger.Info("异步通知结果", string(body))
		if string(body) != "success" {
			order.State = 2
		} else {
			order.State = 3
		}
		order.NotifyRes = string(body)
		if err := db.UpdatePayOrder(order); err != nil {
			logger.Logger.Error(err)
		}
	}
}

func AppPush(type_ int, price float64, metadata string) {
	logger.Logger.Info("app推送", type_, price)
	// 添加到收款记录
	payLog, err := db.AddPaylog(price, strconv.Itoa(type_), metadata)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	// 修改最后收款时间
	setting, err := db.GetSetting("lastPay")
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	setting.VValue = fmt.Sprintf("%d", utils.GetUnix10())
	if err := db.UpdateSetting(setting.VKey, setting.VValue); err != nil {
		logger.Logger.Error(err)
		return
	}
	// 收款通知
	setting, err = db.GetSetting("payNotice")
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	if setting.VValue != "0" {
		switch type_ {
		case 1:
			go SendEmailTask("收款通知", "收到微信收款"+strconv.FormatFloat(price, 'f', 2, 64)+"元")
		case 2:
			go SendEmailTask("收款通知", "收到支付宝收款"+strconv.FormatFloat(price, 'f', 2, 64)+"元")
		}
	}
	// 根据类型与金额获取未过期订单
	var order db.PayOrder
	err = db.DB.Where("state = ? AND type = ? AND really_price = ?", 0, type_, price).First(&order).Error
	if err != nil {
		if err.Error() == "record not found" {
			order.ID = 0
		} else {
			logger.Logger.Error(err)
			return
		}
	}
	if order.ID == 0 {
		logger.Logger.Info("app推送完成, 未找到订单")
		setting, err := db.GetSetting("errorNotice")
		if err != nil {
			logger.Logger.Error(err)
		} else {
			if setting.VValue != "0" {
				switch type_ {
				case 1:
					go SendEmailTask("收款异常", "微信收款异常，未找到订单，金额："+strconv.FormatFloat(price, 'f', 2, 64))
				case 2:
					go SendEmailTask("收款异常", "支付宝收款异常，未找到订单，金额："+strconv.FormatFloat(price, 'f', 2, 64))
				}
			}
		}
		return
	} else {
		logger.Logger.Info("app推送完成, 找到订单")
		payLog.OrderID = order.OrderID
		payLog.Status = 1
		if err := db.UpdatePaylog(payLog); err != nil {
			logger.Logger.Error(err)
			return
		}
	}
	// 更新订单状态
	order.State = 1
	order.CloseDate = utils.GetUnix10()
	if err := db.UpdatePayOrder(order); err != nil {
		logger.Logger.Error(err)
		return
	}
	// 发送异步通知
	if order.ID == 0 {
		logger.Logger.Info("app推送完成, 未找到订单")
		return
	}
	go Notify(order)
	logger.Logger.Info("app推送完成")
}

func DownloadFrontend() {
	logger.Logger.Info("下载前端文件")
	// 下载前端文件
	saveFilePath := "./tmp/frontend.zip"
	uri := "https://github.com/astwy/vmq-go/releases/download/latest/frontend.zip"
	err := utils.DownloadFile(uri, saveFilePath)
	if err != nil {
		logger.Logger.Errorf("下载前端文件失败, %s", err)
		return
	}
}
