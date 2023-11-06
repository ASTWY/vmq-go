package db

type PayQrcode struct {
	ID     uint64  `gorm:"primaryKey" json:"id"`             // 主键
	PayURL string  `json:"payUrl"`                           // 支付地址
	Price  float64 `gorm:"type:decimal(10,2);" json:"price"` // 金额
	Type   int     `json:"type"`                             // 类型 1:微信 2:支付宝
}

// 获取所有数据
func GetPayQrcodes() ([]PayQrcode, error) {
	payQrcodes := []PayQrcode{}
	err := DB.Find(&payQrcodes).Error
	return payQrcodes, err
}

// 获取单个数据
func GetPayQrcodeByID(id uint64) (PayQrcode, error) {
	payQrcode := PayQrcode{}
	err := DB.Where("id = ?", id).First(&payQrcode).Error
	return payQrcode, err
}

// 获取单个数据 通过类型与金额
func GetPayQrcodeByTypeAndPrice(type_ int, price float64) (PayQrcode, error) {
	payQrcode := PayQrcode{}
	err := DB.Where("type = ? AND price = ?", type_, price).First(&payQrcode).Error
	return payQrcode, err
}

// 获取单个数据 通过PayURL
func GetPayQrcodeByPayURL(payURL string) (PayQrcode, error) {
	payQrcode := PayQrcode{}
	err := DB.Where("pay_url = ?", payURL).First(&payQrcode).Error
	return payQrcode, err
}

// 添加数据 返回添加的数据
func AddPayQrcode(payURL string, price float64, type_ int) error {
	payQrcode := PayQrcode{
		PayURL: payURL,
		Price:  price,
		Type:   type_,
	}
	err := DB.Create(&payQrcode).Error
	return err
}

// 删除数据
func DeletePayQrcode(id uint64) error {
	err := DB.Delete(&PayQrcode{}, id).Error
	return err
}

// 更新数据
func UpdatePayQrcode(id uint64, payURL string, price float64, type_ int) error {
	err := DB.Model(&PayQrcode{}).Where("id = ?", id).Updates(map[string]interface{}{
		"pay_url": payURL,
		"price":   price,
		"type":    type_,
	}).Error
	return err
}

// 获取数据总数
func GetPayQrcodeTotal() (int64, error) {
	var count int64
	err := DB.Model(&PayQrcode{}).Count(&count).Error
	return count, err
}
