package controller

//产品列表相关功能

import (
	"config"
	"time"

	"model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//复制
func ProductCopy(c *gin.Context) {
	//获取前端传递的产品信息ID
	IDstr := c.DefaultQuery("ID", "")
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定资源ID",
		})
		return
	}
	ID, _ := strconv.Atoi(IDstr)
	//获取全部源产品信息
	src := model.Product{}
	src.ID = uint(ID)
	orm.Find(&src)
	//包含关联数据，（ProductAttr)
	// orm.Model(&src).Related(&src.ProductAttrs)

	//拷贝新产品
	dst := src //目标新产品
	//清理相关的标识属性
	dst.ID = 0
	dst.Upc = strconv.Itoa(int(time.Now().UnixNano()))
	//将新产品数据存入数据库
	orm.Create(&dst)

	//拷贝关联数据
	orm.Model(&src).Related(&src.ProductAttrs)
	//遍历
	for _, pa := range src.ProductAttrs {
		dstPa := model.ProductAttr{}
		//复制关联属性表的数据，并存入数据库
		dstPa.AttrID = pa.AttrID
		dstPa.Value = pa.Value
		dstPa.ProductID = dst.ID
		orm.Create(&dstPa)
	}
	orm.Model(&dst).Related(&dst.ProductAttrs).Related(&dst.Category)

	//保证属性与值的映射关系
	dst.AttrValue = map[uint]string{}
	for _, pa := range dst.ProductAttrs {
		dst.AttrValue[pa.AttrID] = pa.Value
	}
	dst.ProductAttrs = nil
	if !orm.Model(&dst).Related(&dst.AttrType).RecordNotFound() {
		orm.Model(&dst.AttrType).Related(&dst.AttrType.AttrGroups)
		for ii, ag := range dst.AttrType.AttrGroups {
			orm.Model(&ag).Related(&dst.AttrType.AttrGroups[ii].Attrs)
			for _, a := range dst.AttrType.AttrGroups[ii].Attrs {
				if _, exists := dst.AttrValue[a.ID]; !exists {
					dst.AttrValue[a.ID] = ""
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  dst,
	})

}

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
		// 查询全部产品属性
		products[i].AttrValue = map[uint]string{}
		if !orm.Model(&products[i]).Related(&products[i].ProductAttrs).RecordNotFound() {
			for _, pa := range products[i].ProductAttrs {
				products[i].AttrValue[pa.AttrID] = pa.Value
			}
			products[i].ProductAttrs = nil
		}
		//从产品的属性类型考虑，考虑全部属性
		if !orm.Model(&products[i]).Related(&products[i].AttrType).RecordNotFound() {
			//存在类型
			//根据类型确定全部属性分组
			orm.Model(&products[i].AttrType).Related(&products[i].AttrType.AttrGroups)
			//通过group找到对应的全部属性
			for ii, ag := range products[i].AttrType.AttrGroups {
				//查询到关联的属性
				orm.Model(&ag).Related(&products[i].AttrType.AttrGroups[ii].Attrs)
				for _, a := range products[i].AttrType.AttrGroups[ii].Attrs {
					//得到该产品的全部属性,将不存在的值设置为空
					if _, exists := products[i].AttrValue[a.ID]; !exists {
						products[i].AttrValue[a.ID] = ""
					}
				}
			}
		}

		//查询关联图像
		products[i].Images = []model.Image{}
		orm.Model(&products[i]).Related(&products[i].Images)
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

	//product处理成功，处理关联数据
	//更新product-attr表数据，先获取该表全部属性，再根据前端传递的值，更新对应的字段
	if !orm.Model(&m).Related(&m.AttrType).RecordNotFound() {
		//存在类型
		//根据类型确定全部属性分组
		orm.Model(&m.AttrType).Related(&m.AttrType.AttrGroups)
		//通过group找到对应的全部属性
		for i, ag := range m.AttrType.AttrGroups {
			//查询到关联的属性
			orm.Model(&ag).Related(&m.AttrType.AttrGroups[i].Attrs)
			for _, a := range m.AttrType.AttrGroups[i].Attrs {
				//根据属性选择更新或插入
				pa := model.ProductAttr{}
				if orm.Model(&model.ProductAttr{}).
					Where("product_id=? AND attr_id=?", m.ID, a.ID).
					Find(&pa).RecordNotFound() {
					//不存在
					pa.Value = m.AttrValue[a.ID]
					pa.ProductID = m.ID
					pa.AttrID = a.ID
					orm.Create(&pa)
				} else {
					//已经存在
					pa.Value = m.AttrValue[a.ID]
					orm.Save(&pa)
				}
			}
		}

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

	//更新product-attr表数据，先获取该产品应具有的属性，形成对应的记录，若存在更新置，则执行更新
	if !orm.Model(&m).Related(&m.AttrType).RecordNotFound() {
		//存在类型
		//根据类型确定全部属性分组
		orm.Model(&m.AttrType).Related(&m.AttrType.AttrGroups)
		//通过group找到对应的全部属性
		for i, ag := range m.AttrType.AttrGroups {
			//查询到关联的属性
			orm.Model(&ag).Related(&m.AttrType.AttrGroups[i].Attrs)
			for _, a := range m.AttrType.AttrGroups[i].Attrs {
				//根据属性选择更新或插入
				pa := model.ProductAttr{}
				if orm.Model(&model.ProductAttr{}).
					Where("product_id=? AND attr_id=?", m.ID, a.ID).
					Find(&pa).RecordNotFound() {
					//不存在
					pa.Value = m.AttrValue[a.ID]
					pa.ProductID = m.ID
					pa.AttrID = a.ID
					orm.Create(&pa)
				} else {
					//已经存在
					pa.Value = m.AttrValue[a.ID]
					orm.Save(&pa)
				}
			}
		}
	}

	//更新image表
	for i, img := range m.UploadedImage {
		//存储格式：a/b/xxxx.jpg
		image := model.Image{}
		image.ProductID = m.ID
		image.Host = config.App["IMAGE_HOST"]
		image.Image = string(img[0]) + "/" + string(img[1]) + "/" + img
		image.ImageSmall = string(m.UploadedImageSmall[i][0]) + "/" + string(m.UploadedImageSmall[i][1]) + "/" + m.UploadedImageSmall[i]
		image.ImageBig = string(m.UploadedImageBig[i][0]) + "/" + string(m.UploadedImageBig[i][1]) + "/" + m.UploadedImageBig[i]
		orm.Create(&image)
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
