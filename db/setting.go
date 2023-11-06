package db

type Setting struct {
	ID     uint   `gorm:"primaryKey" json:"id"` // ID
	VKey   string `json:"vkey"`                 // 键
	VValue string `json:"vvalue"`               // 值
}

// 获取所有数据
func GetSettings() ([]Setting, error) {
	settings := []Setting{}
	err := DB.Find(&settings).Error
	return settings, err
}

// 获取单个数据
func GetSetting(key string) (Setting, error) {
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

// VerifyAdmin
func VerifyAdmin(username string, password string) bool {
	// adminUser adminPwd
	setting, err := GetSetting("adminUser")
	if err != nil {
		return false
	}
	if setting.VValue != username {
		return false
	}
	setting, err = GetSetting("adminPwd")
	if err != nil {
		return false
	}
	if setting.VValue != password {
		return false
	}
	return true
}
