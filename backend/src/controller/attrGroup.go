package controller

//AttrGroup 表控制器（增删改查）代码，脚手架模板

import (
	"model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//AttrGroup列表
func AttrGroupList(c *gin.Context) {
	//搜索(筛选）
	condStr := ""
	cond := []string{}
	condParams := []string{}
	//确定搜索条件
	filterAttrTypeID := c.DefaultQuery("filterAttrTypeID", "")
	if filterAttrTypeID != "" {
		cond = append(cond, "attr_type_id = ?")
		condParams = append(condParams, filterAttrTypeID)
	}
	condStr = strings.Join(cond, "AND")

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
	orm.Model(&model.AttrGroup{}).Where(condStr, condParams).Count(&total)
	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//获取product模型
	ms := []model.AttrGroup{}
	//获取展示数量和偏移量,输出数据获
	orm.Where(condStr, condParams).Order(orderStr).Limit(pageSize).Offset(offset).Find(&ms)
	//遍历全部属性，找到关联字段
	withAttr := c.DefaultQuery("withAttr","")
	for i, m := range ms {
		//关联类型
		orm.Model(&m).Related(&ms[i].AttrType)
		//查询关联属性
		if withAttr == "yes"{
			orm.Model(&m).Related(&ms[i].Attrs)
		}
	}

	//响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  ms,
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

//AttrGroup删除
func AttrGroupDelete(c *gin.Context) {
	//获取前端参数（需要删除的ID）
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定ID",
		})
		return
	}
	//确定模型对象（表）
	m := model.AttrGroup{}
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

//AttrGroup添加
func AttrGroupCreate(c *gin.Context) {
	//确定模型对象（表）
	m := model.AttrGroup{}
	//使用c.ShouldBind(),绑定并解析post数据
	err := c.ShouldBind(&m)
	if err != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	//特定数据设置

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

	//响应正确数据
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  m,
	})
}

//AttrGroup更新
func AttrGroupUpdate(c *gin.Context) {
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
	m := model.AttrGroup{}
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

	//响应正确数据
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  m,
	})
}