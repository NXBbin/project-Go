package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//前端临时的订单数据
type TempOrder struct {
	AddressID    uint
	BuyProductID []uint
	ShippingID   uint
	UserID       uint
}

//处理订单
func OrderResult(c *gin.Context) {
	// 获取前端请求传递的sn，向redis查询是否存在数据
	sn := c.DefaultQuery("sn", "")
	result, err := Rds.Do("HGET", "orderResult", sn)
	if err != nil || result == nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "订单不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  string(result.([]byte)),
	})

}

//订单生成
func OrderCreate(c *gin.Context) {
	user := member(c)
	if user == nil {
		//未验证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在该会员",
		})
		return
	}

	//1.生成订单SN（日期+当天序号）
	now := time.Now()
	//得到redis中的计数器key，以天为标识
	key := fmt.Sprintf("%d%d%d",
		now.Year(),
		now.Month(),
		now.Day(),
	)
	//在redis使用的序号递增key
	counterKey := "counter" + key
	n, err := Rds.Do("incr", counterKey)
	// log.Println("--------", n, err)
	if err != nil {
		//没有序号，生成随机序号
		n = rand.Int63n(10000000)
	}
	//转为字符串类型
	ns := strconv.Itoa(int(n.(int64)))
	// 固定长度补齐
	if l := len(ns); l < 8 {
		// 当长度小于8位时补0
		ns = strings.Repeat("0", 8-l) + ns
	}
	// 订单号
	sn := key + ns
	log.Println("订单号：", sn)

	// 2.临时存储订单信息
	tempOrder := TempOrder{}
	c.ShouldBind(&tempOrder)
	//记录订单对应的用户信息
	tempOrder.UserID = user.ID
	//记录订单中的产品信息

	//JSON化
	toj, err := json.Marshal(tempOrder)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "提交订单数据错误",
		})
		return
	}

	// 利用hash结构，存储在redis中。sn为key键，toj为value值
	_, hseterr := Rds.Do("HSET", "tempOrder", sn, toj)
	if hseterr != nil {
		//redis缓存失败，应使用其他方案
	}

	// 3.将订单放入队列（等待处理）
	// 先获取队列长度，提示：
	waitLen, _ := Rds.Do("XLEN", "orderQueue")
	//放入队列
	_, addErr := Rds.Do("XADD", "orderQueue", "*", "content", sn)
	if addErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "订单队列错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   "",
		"data":    sn,
		"waitLen": waitLen,
	})
}

//订单列表
func ShippingList(c *gin.Context) {
	user := member(c)
	if user == nil {
		//未验证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在该会员",
		})
		return
	}
	data := []model.Shipping{}
	// 查询状态为1的数据
	orm.Where("status=1").Find(&data)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  data,
	})
}
