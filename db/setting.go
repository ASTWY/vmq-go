package db

import "strconv"

type Setting struct {
	ID     uint   `gorm:"primaryKey" json:"id"` // ID
	VKey   string `json:"vkey"`                 // 键
	VValue string `json:"vvalue"`               // 值
}

type AppConfig struct {
	AdminPwd      string `json:"adminPwd"`
	AdminUser     string `json:"adminUser"`
	AliPay        string `json:"aliPay"`
	APISecret     string `json:"apiSecret"`
	EmailSMTPfrom string `json:"emailSMTPfrom"`
	EmailSMTPhost string `json:"emailSMTPhost"`
	EmailSMTPport int    `json:"emailSMTPport"`
	EmailSMTPpwd  string `json:"emailSMTPpwd"`
	EmailSMTPssl  bool   `json:"emailSMTPssl"`
	EmailSMTPto   string `json:"emailSMTPto"`
	EmailSMTPuser string `json:"emailSMTPuser"`
	ErrorNotice   bool   `json:"errorNotice"`
	Expire        int    `json:"expire"`
	LastHeart     int    `json:"lastHeart"`
	LastPay       int    `json:"lastPay"`
	MonitorNotice bool   `json:"monitorNotice"`
	NotifyUrl     string `json:"notifyUrl"`
	OrderMaxNum   int    `json:"orderMaxNum"`
	OrderType     int    `json:"orderType"`
	PayNotice     bool   `json:"payNotice"`
	ReturnUrl     string `json:"returnUrl"`
	WechatPay     string `json:"wechatPay"`
}

func GetAppConfig() (AppConfig, error) {
	settings, err := GetSettings()
	if err != nil {
		return AppConfig{}, err
	}

	var appConfig AppConfig

	for _, v := range settings {
		switch v.VKey {
		case "adminPwd":
			appConfig.AdminPwd = v.VValue
		case "adminUser":
			appConfig.AdminUser = v.VValue
		case "aliPay":
			appConfig.AliPay = v.VValue
		case "apiSecret":
			appConfig.APISecret = v.VValue
		case "emailSMTPfrom":
			appConfig.EmailSMTPfrom = v.VValue
		case "emailSMTPhost":
			appConfig.EmailSMTPhost = v.VValue
		case "emailSMTPport":
			appConfig.EmailSMTPport, _ = strconv.Atoi(v.VValue)
		case "emailSMTPpwd":
			appConfig.EmailSMTPpwd = v.VValue
		case "emailSMTPssl":
			appConfig.EmailSMTPssl = v.VValue != "0"
		case "emailSMTPto":
			appConfig.EmailSMTPto = v.VValue
		case "emailSMTPuser":
			appConfig.EmailSMTPuser = v.VValue
		case "errorNotice":
			appConfig.ErrorNotice = v.VValue != "0"
		case "expire":
			appConfig.Expire, _ = strconv.Atoi(v.VValue)
		case "lastHeart":
			appConfig.LastHeart, _ = strconv.Atoi(v.VValue)
		case "lastPay":
			appConfig.LastPay, _ = strconv.Atoi(v.VValue)
		case "monitorNotice":
			appConfig.MonitorNotice = v.VValue != "0"
		case "notifyUrl":
			appConfig.NotifyUrl = v.VValue
		case "orderMaxNum":
			appConfig.OrderMaxNum, _ = strconv.Atoi(v.VValue)
		case "orderType":
			appConfig.OrderType, _ = strconv.Atoi(v.VValue)
		case "payNotice":
			appConfig.PayNotice = v.VValue != "0"
		case "returnUrl":
			appConfig.ReturnUrl = v.VValue
		case "wechatPay":
			appConfig.WechatPay = v.VValue
		}
	}
	return appConfig, nil
}

func Keys() []string {
	return []string{
		"adminPwd",
		"adminUser",
		"aliPay",
		"apiSecret",
		"emailSMTPfrom",
		"emailSMTPhost",
		"emailSMTPport",
		"emailSMTPpwd",
		"emailSMTPssl",
		"emailSMTPto",
		"emailSMTPuser",
		"errorNotice",
		"expire",
		"lastHeart",
		"lastPay",
		"monitorNotice",
		"notifyUrl",
		"orderMaxNum",
		"orderType",
		"payNotice",
		"returnUrl",
		"wechatPay",
	}
}

func (appConfig *AppConfig) VerifyAdmin(username string, password string) bool {
	if appConfig.AdminUser == username && appConfig.AdminPwd == password {
		return true
	}
	return false
}

// 检查表数据是否存在 用于初始化时检查所有的配置项是否存在，不存在则添加
func CheckSetting() {
	keys := Keys()
	for _, key := range keys {
		setting, err := getSetting(key)
		if err != nil {
			AddSetting(key, "")
		}
		if setting.VValue == "" {
			UpdateSetting(key, "")
		}
	}
}

// 获取所有数据
func GetSettings() ([]Setting, error) {
	settings := []Setting{}
	err := DB.Find(&settings).Error
	return settings, err
}

// 获取单个数据
func getSetting(key string) (Setting, error) {
	setting := Setting{}
	err := DB.Where("v_key = ?", key).First(&setting).Error
	return setting, err
}

// 添加数据 返回添加的数据
func AddSetting(key string, value string) (Setting, error) {
	setting := Setting{
		VKey:   key,
		VValue: value,
	}
	if err := DB.Create(&setting).Error; err != nil {
		return Setting{}, err
	}
	// 取最新的一条数据
	var settingNew Setting
	if err := DB.Last(&settingNew).Error; err != nil {
		return Setting{}, err
	}
	return settingNew, nil
}

// 更新数据 返回更新的数据
func UpdateSetting(key string, value string) error {
	// 设置回滚点
	tx := DB.Begin()
	// 更新数据
	if err := tx.Model(&Setting{}).Where("v_key = ?", key).Update("v_value", value).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	tx.Commit()
	return nil
}

// 删除数据
func DeleteSetting(key string) error {
	tx := DB.Begin()
	if err := tx.Where("v_key = ?", key).Delete(&Setting{}).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}
