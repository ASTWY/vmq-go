package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"vmq-go/db"
	"vmq-go/task"

	"github.com/gin-gonic/gin"
)

func getOrderHandler(c *gin.Context) {
	var orderList []db.PayOrder
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
	orderList, err = db.GetPayOrderListDesc(pageInt, pageSizeInt)
	if err != nil {
		c.Error(err)
		return
	}
	count, err := db.GetPayOrderCount()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", gin.H{
		"list":  orderList,
		"total": count,
	})
}

func reCallbackOrderHandler(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		c.Error(fmt.Errorf("orderId is empty"))
		return
	}
	order, err := db.GetPayOrderByOrderID(orderId)
	if err != nil {
		c.Error(err)
		return
	}
	task.Notify(order)
	c.Set("code", http.StatusOK)
}

func deleteOrderHandler(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		c.Error(fmt.Errorf("orderId is empty"))
		return
	}
	order, err := db.GetPayOrderByOrderID(orderId)
	if err != nil {
		c.Error(err)
		return
	}
	err = db.DeletePayOrder(order)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("code", http.StatusOK)
}

func getOrderDataTodayHandler(c *gin.Context) {
	orderList, err := db.GetTodayPayOrders()
	if err != nil {
		c.Error(err)
		return
	}
	var todaydataMap = make(map[string]string)
	var income float64
	var orders int = len(orderList)
	var successOrders int
	var failedOrders int
	for _, order := range orderList {
		if order.State == 3 {
			successOrders++
			income += order.ReallyPrice
		} else if order.State == -1 {
			failedOrders++
		}
	}
	todaydataMap["income"] = fmt.Sprintf("%.2f", income)
	todaydataMap["orders"] = strconv.Itoa(orders)
	todaydataMap["successOrders"] = strconv.Itoa(successOrders)
	todaydataMap["failedOrders"] = strconv.Itoa(failedOrders)
	c.Set("data", todaydataMap)
}
