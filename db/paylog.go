package db

import (
	"vmq-go/utils"
)

// 数据表结构体
type Paylog struct {
	ID         uint    `gorm:"primary_key" json:"id"`            // ID
	Price      float64 `gorm:"type:decimal(10,2);" json:"price"` // 金额 保留两位小数
	Type       string  `json:"type"`                             // 类型 1:微信 2:支付宝
	Status     int     `gorm:"default:0" json:"state"`           // 状态 0:未匹配到订单 1:已匹配到订单 默认为0
	OrderID    string  `json:"orderId"`                          // 订单号 默认为空 未匹配到订单时为空
	Createtime int64   `json:"createtime"`                       // 创建时间 10位时间戳
	Metadata   string  `json:"metadata"`                         // 元数据
}

// 添加数据 返回添加的数据
func AddPaylog(price float64, payType string, metadata string) (Paylog, error) {
	paylog := Paylog{
		Price:      price,
		Type:       payType,
		Status:     0,
		OrderID:    "",
		Createtime: utils.GetUnix10(),
		Metadata:   metadata,
	}
	if err := DB.Create(&paylog).Error; err != nil {
		return Paylog{}, err
	}
	// 取最新的一条数据
	var paylogNew Paylog
	if err := DB.Last(&paylogNew).Error; err != nil {
		return Paylog{}, err
	}
	return paylogNew, nil
}

// 更新数据
func UpdatePaylog(paylog Paylog) error {
	tx := DB.Begin()
	if err := tx.Model(&Paylog{}).Where("id = ?", paylog.ID).Updates(paylog).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 获取所有数据 带分页 以id倒序
func GetPaylogListDesc(page int, pageSize int) ([]Paylog, error) {
	paylogs := []Paylog{}
	err := DB.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&paylogs).Error
	return paylogs, err
}

// 获取数据总数
func GetPaylogCount() (int64, error) {
	var count int64
	err := DB.Model(&Paylog{}).Count(&count).Error
	return count, err
}
