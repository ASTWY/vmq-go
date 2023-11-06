package db

import (
	"fmt"
	"vmq-go/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 初始化数据库连接
func InitDB(dsn string) error {
	var err error = nil
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("连接数据库失败：%s", err.Error())
	}
	return err
}

// 迁移数据库
func Migrate() error {
	err := DB.AutoMigrate(&PayOrder{}, &PayQrcode{}, &Setting{}, &Paylog{})
	return err
}

// 初始化数据
func initializeData() error {
	// 初始化 setting 数据
	// 如果setting表中没有数据，则初始化数据
	var settingCount int64
	if err := DB.Model(&Setting{}).Count(&settingCount).Error; err != nil {
		return err
	}
	if settingCount > 0 {
		return nil
	}
	settingData := map[string]string{
		"adminUser":     "admin",                            // 管理员账号
		"adminPwd":      "21232f297a57a5a743894a0e4a801fc3", // 管理员密码
		"notifyUrl":     "",                                 // 异步回调地址
		"returnUrl":     "",                                 // 同步回调地址
		"apiSecret":     "",                                 // 通讯密钥
		"lastHeart":     "0",                                // 最后心跳时间
		"lastPay":       "0",                                // 最后支付时间
		"expire":        "5",                                // 过期时间 单位分钟
		"orderType":     "1",                                // 订单区分 1：金额递减 2：金额递增
		"orderMaxNum":   "10",                               // 同金额订单最大数量
		"wechatPay":     "",                                 // 微信收款链接
		"aliPay":        "",                                 // 支付宝收款链接
		"emailSMTPhost": "",                                 // 邮箱SMTP地址
		"emailSMTPport": "",                                 // 邮箱SMTP端口
		"emailSMTPuser": "",                                 // 邮箱SMTP账号
		"emailSMTPpwd":  "",                                 // 邮箱SMTP密码
		"emailSMTPfrom": "",                                 // 邮箱SMTP发件人
		"emailSMTPto":   "",                                 // 邮箱SMTP收件人
		"emailSMTPssl":  "1",                                // 邮箱SMTP是否开启SSL 0否 1是
		"payNotice":     "0",                                // 收款通知 0否 1是
		"errorNotice":   "1",                                // 异常通知 0否 1是
		"monitorNotice": "1",                                // 监控通知 0否 1是
	}
	keys, data := utils.DictionaryOrderSort(settingData)
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		value := data[key]
		setting := Setting{
			VKey:   key,
			VValue: value,
		}
		if err := DB.Create(&setting).Error; err != nil {
			return err
		}
	}
	return nil
}

// 初始化数据库
func SetupDatabase(dsn string) error {
	if err := InitDB(dsn); err != nil {
		return err
	}
	// 迁移数据库
	if err := Migrate(); err != nil {
		return err
	}
	// 初始化数据
	if err := initializeData(); err != nil {
		return err
	}
	return nil
}
