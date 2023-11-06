package admin

import (
	"fmt"
	"strconv"
	"vmq-go/db"

	"github.com/gin-gonic/gin"
)

func getQrcodeHandler(c *gin.Context) {
	qrcodeList, err := db.GetPayQrcodes()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", qrcodeList)
}

func postQrcodeHandler(c *gin.Context) {
	var params map[string]string
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Error(err)
		return
	}
	payUrl := params["payUrl"]
	price, err := strconv.ParseFloat(params["price"], 64)
	if err != nil {
		c.Error(fmt.Errorf("price 类型错误"))
		return
	}
	type_, err := strconv.Atoi(params["type"])
	if err != nil {
		c.Error(fmt.Errorf("type 类型错误"))
		return
	}
	qc, err := db.GetPayQrcodeByPayURL(payUrl)
	if err != nil {
		if err.Error() != "record not found" {
			c.Error(err)
			return
		}
	}
	if qc.ID != 0 {
		c.Set("message", "该二维码已存在")
		return
	}
	if err := db.AddPayQrcode(payUrl, price, type_); err != nil {
		c.Error(err)
		return
	}
	c.Set("message", "添加成功")
}

func deleteQrcodeHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	if err := db.DeletePayQrcode(id); err != nil {
		c.Error(err)
		return
	}
	c.Set("message", "删除成功")
}
