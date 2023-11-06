package admin

import (
	"strconv"

	"vmq-go/db"

	"github.com/gin-gonic/gin"
)

func getPaylogHandler(c *gin.Context) {
	var paylogList []db.Paylog
	var err error
	page := c.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.Error(err)
		return
	}
	pageSize := c.DefaultQuery("pageSize", "10")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	paylogList, err = db.GetPaylogListDesc(pageInt, pageSizeInt)
	if err != nil {
		c.Error(err)
		return
	}
	count, err := db.GetPaylogCount()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", gin.H{
		"list":  paylogList,
		"total": count,
	})
}
