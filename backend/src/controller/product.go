package controller

//产品列表相关功能

import (
	// "log"
	"time"

	"model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//产品列表
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

//产品添加
func ProductCreate(c *gin.Context) {
	//确定模型对象（表）
	m := model.Product{}
	//使用c.ShouldBind(),绑定并解析post数据
	err := c.ShouldBind(&m)
	if err != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	//判断用户是否有设置Upc数据（Upc是唯一键，不设置会报错）
	if "" == m.Upc {
		//未设置的情况下，自动生成随机值填充
		m.Upc = strconv.Itoa(int(time.Now().UnixNano()))
	}

	// 将关联临时关闭，(若不关闭，也会将添加的数据自动添加到关联的表中)，并添加数据
	orm.Set("gorm:save_associations", false).Create(&m)
	if orm.Error != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
		return
	}

	// 查询相关联的表数据
	category := model.Category{}
	//
	category.ID = m.CategoryID
	orm.Find(&category)
	//
	m.Category = category

	//响应正确数据
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  m,
	})
}

//产品更新
func ProductUpdate(c *gin.Context) {
	//获取前端传递的请求更新数据ID
	IDstr := c.DefaultQuery("ID", "")
	// 没有传递时，错误响应
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定更新ID",
		})
		return
	}
	//类型转换（string转int）
	ID, _ := strconv.Atoi(IDstr)

	// 获取需要更新的表数据
	m := model.Product{}
	m.ID = uint(ID)
	orm.Find(&m)

	// 绑定并解析 数据
	err := c.ShouldBind(&m)
	if err != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 临时取消关联,并更新数据
	orm.Set("gorm:save_associations", false).Save(&m)
	if orm.Error != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
		return
	}

	// 查询相关联的表
	category := model.Category{}
	//通过前端传递的关联ID，找到相关联的表数据
	category.ID = m.CategoryID
	//获取表中全部数据	 // SELECT * FROM category;
	orm.Find(&category)
	//将关联的分类的数据，赋值到产品的属性中
	m.Category = category
	// log.Println("-----", m)
	// log.Println("++++++", category)

	//响应正确数据
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  m,
	})
}
