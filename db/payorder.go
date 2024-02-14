package db

import (
	"fmt"
	"time"
	"vmq-go/utils"
)

type PayOrder struct {
	ID          uint    `gorm:"primaryKey" json:"id"`                   // ID
	CloseDate   int64   `json:"closeDate"`                              // 关闭时间
	ExpectDate  int64   `json:"expectDate"`                             // 过期时间
	CreateDate  int64   `json:"createDate"`                             // 创建时间
	IsAuto      int     `json:"isAuto"`                                 // 是否不需要输入金额 0否 1是
	NotifyURL   string  `json:"notifyUrl"`                              // 异步回调地址 一般用于接收通知 安全性高
	NotifyRes   string  `json:"notifyRes"`                              // 异步回调结果
	OrderID     string  `json:"orderId"`                                // 订单号
	Param       string  `json:"param"`                                  // 自定义参数
	PayDate     int64   `json:"payDate"`                                // 支付时间
	PayID       string  `json:"payId"`                                  // 商户订单号
	PayURL      string  `json:"payUrl"`                                 // 支付地址
	Price       float64 `gorm:"type:decimal(12,2);" json:"price"`       // 金额 保留两位小数
	ReallyPrice float64 `gorm:"type:decimal(12,2);" json:"reallyPrice"` // 实际支付金额 保留两位小数
	ReturnURL   string  `json:"returnUrl"`                              // 同步回调地址 一般用于跳转 安全性低
	State       int     `gorm:"default:0" json:"state"`                 // 状态 -1过期 0未支付 1已支付 2已支付但通知失败 3 已完成
	Type        int     `json:"type"`                                   // 类型 1:微信 2:支付宝
}

// 获取今日订单
func GetTodayPayOrders() ([]PayOrder, error) {
	payOrders := []PayOrder{}
	today := time.Now().Format("2006-01-02")
	todayDateUnix, err := time.ParseInLocation("2006-01-02", today, time.Local)
	if err != nil {
		return payOrders, err
	}
	err = DB.Where("create_date > ?", todayDateUnix.Unix()).Find(&payOrders).Error
	return payOrders, err
}

// 获取所有数据
func GetPayOrders() ([]PayOrder, error) {
	payOrders := []PayOrder{}
	err := DB.Find(&payOrders).Error
	return payOrders, err
}

// 获取所有数据 带分页 正序
func GetPayOrderList(page int, pageSize int) ([]PayOrder, error) {
	payOrders := []PayOrder{}
	err := DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&payOrders).Error
	return payOrders, err
}

// 获取所有数据 带分页 倒序
func GetPayOrderListDesc(page int, pageSize int) ([]PayOrder, error) {
	payOrders := []PayOrder{}
	err := DB.Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&payOrders).Error
	return payOrders, err
}

// 获取数据总数
func GetPayOrderCount() (int64, error) {
	var count int64
	err := DB.Model(&PayOrder{}).Count(&count).Error
	return count, err
}

// 获取单个数据
func GetPayOrderByPayID(payID string) (PayOrder, error) {
	payOrder := PayOrder{}
	err := DB.Where("pay_id = ?", payID).First(&payOrder).Error
	return payOrder, err
}

// 获取单个数据
func GetPayOrderByOrderID(orderID string) (PayOrder, error) {
	payOrder := PayOrder{}
	err := DB.Where("order_id = ?", orderID).First(&payOrder).Error
	return payOrder, err
}

// 添加数据 返回添加的数据
func AddPayOrder(payId string, type_ int, price float64, param string, notifyUrl string, returnUrl string) error {
	createDate := utils.GetUnix10()
	appConfig, err := GetAppConfig()
	if err != nil {
		return err
	}
	expectDate := createDate + int64(appConfig.Expire*60)
	payqrcode, err := GetPayQrcodeByTypeAndPrice(type_, price)
	if err != nil {
		if err.Error() == "record not found" {
			payqrcode.ID = 0
		} else {
			return err
		}
	}
	var payUrl string
	var isAuto int
	if payqrcode.ID == 0 {
		switch type_ {
		case 1:
			payUrl = appConfig.WechatPay
		case 2:
			payUrl = appConfig.AliPay
		}
		isAuto = 0
	} else {
		payUrl = payqrcode.PayURL
		isAuto = 1
	}
	// 取未支付的订单 过滤器为 price type 并以真实支付金额升序排序 取出所有的真实支付金额
	var reallyPriceList []float64
	err = DB.Model(&PayOrder{}).Where("state = 0 AND price = ? AND type = ?", price, type_).Order("really_price asc").Pluck("really_price", &reallyPriceList).Error
	if err != nil {
		return err
	}
	consecutive, missPrices := utils.CeckConsecutive(reallyPriceList)
	var reallyPrice float64
	switch appConfig.OrderType {
	case 1:
		// 金额递减
		if len(reallyPriceList) == 0 || !utils.BinarySearch(reallyPriceList, price) {
			reallyPrice = price
		} else {
			if consecutive {
				reallyPrice = reallyPriceList[0] - 0.01
			} else {
				// 如果missPrices中没有值为price的值 则取missPrices中的第一个值
				reallyPrice = missPrices[0]
			}
		}
	case 2:
		// 金额递增
		if len(reallyPriceList) == 0 || !utils.BinarySearch(reallyPriceList, price) {
			reallyPrice = price
		} else {
			if consecutive {
				reallyPrice = reallyPriceList[0] + 0.01
			} else {
				// 如果missPrices中没有值为price的值 则取missPrices中的第一个值
				reallyPrice = missPrices[0]
			}
		}
	}
	if len(reallyPriceList) >= appConfig.OrderMaxNum || reallyPrice < 0.01 {
		return fmt.Errorf("订单已满, 请稍后再试")
	}
	payOrder := PayOrder{
		PayID:       payId,
		PayURL:      payUrl,
		IsAuto:      isAuto,
		ExpectDate:  expectDate,
		CreateDate:  createDate,
		Type:        type_,
		Price:       price,
		ReallyPrice: reallyPrice,
		Param:       param,
		NotifyURL:   notifyUrl,
		ReturnURL:   returnUrl,
		State:       0,
		OrderID:     utils.GenerateOrderNo(),
	}
	if err := DB.Create(&payOrder).Error; err != nil {
		return err
	}
	return nil
}

// 更新数据
func UpdatePayOrder(payOrder PayOrder) error {
	tx := DB.Begin()
	if err := tx.Model(&PayOrder{}).Where("id = ?", payOrder.ID).Updates(payOrder).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 删除数据
func DeletePayOrder(payOrder PayOrder) error {
	tx := DB.Begin()
	if err := tx.Where("id = ?", payOrder.ID).Delete(&payOrder).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
