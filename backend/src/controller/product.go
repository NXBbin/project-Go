package controller

//产品列表相关功能

import (
	// "log"
	"model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//产品列表接口
func ProductList(c *gin.Context) {
	//搜索(筛选）
	condStr := ""
	condParams := []string{}
	//获取前端传递的筛选名
	filterName := c.DefaultQuery("filterName", "")
	//若不等于空字符串，说明有传递
	if filterName != "" {
		//SQL筛选条件语句
		condStr = "name like ?"
		condParams = append(condParams, filterName+"%")
		// log.Println("where语句：",condStr,condParams)
	}
	//多条件

	//排序
	//获取排序参数（字段）
	orderStr := ""
	sortProp := c.DefaultQuery("sortProp", "")
	//排序方式：ascending ， descending
	sortOrder := c.DefaultQuery("sortOrder", "")
	//判断用户是否请求了排序
	if sortProp != "" && sortOrder != "" {
		//默认是升序，若传递的是降序请求，则设置为DESC
		sortMethod := "ASC"
		if "descending" == sortOrder {
			sortMethod = "DESC"
		}
		//拼凑：name ASC||DESC
		orderStr = sortProp + " " + sortMethod
		// log.Println("order语句：", orderStr)
	}

	//翻页: /products?currentPage =&pageSize=
	//获取请求的当前页码,默认第一页
	currentPageStr := c.DefaultQuery("currentPage", "1")
	//每页的显示的数量（偏移量）
	pageSizeStr := c.DefaultQuery("pageSize", "5")
	//将从前端获取到的页码数据转换类型（int转string）
	currentPage, pageErr := strconv.Atoi(currentPageStr)
	//若用户传递的参数不是整形数据（不合法数据），则指定页码为1
	if pageErr != nil {
		currentPage = 1
	}

	pageSize, sizeErr := strconv.Atoi(pageSizeStr)
	if sizeErr != nil {
		pageSize = 5
	}

	//获取总记录数
	total := 0
	orm.Model(&model.Product{}).Where(condStr, condParams).Count(&total)
	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//获取product模型
	products := []model.Product{}
	//获取展示数量和偏移量,输出数据获
	orm.Where(condStr, condParams).Order(orderStr).Limit(pageSize).Offset(offset).Find(&products)
	//遍历全部属性，找到关联字段
	for i, _ := range products {
		orm.Model(&products[i]).Related(&products[i].Category)
	}

	//响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  products,
		//页数
		"pager": map[string]int{
			//当前页
			"currentPage": currentPage,
			//偏移量
			"pageSize": pageSize,
			//数据总量
			"total": total,
		},
	})
}

//产品删除
func ProductDelete(c *gin.Context) {
	//获取前端参数（需要删除的ID）
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定ID",
		})
		return
	}
	//确定模型对象（表）
	m := model.Product{}
	//将前端数据进行类型转换（string转uint）
	id, _ := strconv.Atoi(ID)
	m.ID = uint(id)
	//执行删除
	orm.Delete(&m)
	//判断删除是否有错误
	if orm.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
	}
	//无错误响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
	})
}
