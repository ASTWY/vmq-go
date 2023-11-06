package admin

import (
	"fmt"
	"vmq-go/db"
	"vmq-go/task"

	"github.com/gin-gonic/gin"
)

// getSettingsHandler
func getSettingsHandler(c *gin.Context) {
	settings, err := db.GetSettings()
	if err != nil {
		c.Error(err)
		return
	}
	var dataMap = make(map[string]string)
	for _, setting := range settings {
		if setting.VKey == "adminPwd" {
			dataMap[setting.VKey] = "******"
		} else {
			dataMap[setting.VKey] = setting.VValue
		}
	}
	c.Set("data", dataMap)
}

// putSettingHandler
func putSettingHandler(c *gin.Context) {
	var paramMap map[string]string
	if err := c.ShouldBindJSON(&paramMap); err != nil {
		c.Error(err)
		return
	}
	if paramMap["key"] == "" || paramMap["value"] == "" {
		c.Error(fmt.Errorf("key or value is empty"))
		return
	}
	err := db.UpdateSetting(paramMap["key"], paramMap["value"])
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("message", "更新成功")
}

// sendEmailHandler
func sendEmailHandler(c *gin.Context) {
	subject := "测试邮件"
	body := "这是一封来自VMQ-Go的测试邮件，如果您收到了这封邮件，说明您的邮件服务器配置正确。"
	task.SendEmailTask(subject, body)
	c.Set("message", "发送成功")
}
