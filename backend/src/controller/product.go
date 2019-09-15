package controller

import (
	"net/http"
	"model"
	"github.com/gin-gonic/gin"
)

//产品列表接口
func ProductList(c *gin.Context){
	//搜索
	//排序
	//翻页
	
	//获取product模型
	products := []model.Product{}
	//获取模型全部属性
	orm.Find(&products)
	//遍历全部属性，找到关联字段
	for i,_ := range products{
		orm.Model(&products[i]).Related(&products[i].Category)
	}
	
	//响应
	c.JSON(http.StatusOK, gin.H{
		"error":"",
		"data":products,
	})
}