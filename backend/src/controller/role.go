package controller

//Role 表控制器（增删改查）代码，脚手架模板

import (
	// "log"
	"model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//授权管理
func RoleGrant(c *gin.Context) {
	//确定资源
	IDstr := c.DefaultQuery("ID", "")
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "请输入资源ID",
		})
		return
	}
	ID, _ := strconv.Atoi(IDstr)
	role := model.Role{}
	role.ID = uint(ID)
	orm.Find(&role)

	//得到目标关联对象切片
	checked := c.PostFormArray("checked")
	ps := []model.Privilege{}
	for _, pid := range checked {
		p := model.Privilege{}
		pID, _ := strconv.Atoi(pid)
		p.ID = uint(pID)
		ps = append(ps, p)
	}

	//更新关联
	orm.Model(&role).Association("Privileges").Replace(ps)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
	})
}

//获取权限和已选
func PrivilegeRole(c *gin.Context) {
	//全部权限
	ps := []model.Privilege{}
	orm.
		Order("sort_order").
		Find(&ps)
	//角色已选权限
	IDstr := c.DefaultQuery("ID", "")
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "请输入资源ID",
		})
		return
	}
	ID, _ := strconv.Atoi(IDstr)
	//确定角色模型
	role := model.Role{}
	role.ID = uint(ID)
	// 查询与之关联的全部权限，已选权限
	cps := []model.Privilege{}
	orm.Model(&role).Related(&cps, "Privileges")
	rs := []uint{}
	for _, p := range cps {
		rs = append(rs, p.ID)
	}

	//组合数据
	data := struct {
		Privileges []model.Privilege
		Checked    []uint
	}{
		ps,
		rs,
	}
	// log.Printf("------%v\n", ps)
	// log.Printf("+++++%v\n", rs)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  data,
	})
}

//Role列表
func RoleList(c *gin.Context) {
	//搜索(筛选）
	condStr := ""
	condParams := []string{}
	//确定搜索条件

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

	//将从前端获取到的页码数据转换类型（int转string）
	currentPage, pageErr := strconv.Atoi(currentPageStr)
	//若用户传递的参数不是整形数据（不合法数据），则指定页码为1
	if pageErr != nil {
		currentPage = 1
	}
	
	//每页的显示的数量（偏移量）
	pageSizeStr := c.DefaultQuery("pageSize", "5")
	pageSize, sizeErr := strconv.Atoi(pageSizeStr)
	if sizeErr != nil {
		pageSize = 5
	}

	//获取总记录数
	total := 0
	orm.Model(&model.Role{}).Where(condStr, condParams).Count(&total)
	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//获取product模型
	ms := []model.Role{}
	//获取展示数量和偏移量,输出数据获
	orm.Where(condStr, condParams).Order(orderStr).Limit(pageSize).Offset(offset).Find(&ms)
	//遍历全部属性，找到关联字段

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

//Role删除
func RoleDelete(c *gin.Context) {
	//获取前端参数（需要删除的ID）
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定ID",
		})
		return
	}
	//确定模型对象（表）
	m := model.Role{}
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

//Role添加
func RoleCreate(c *gin.Context) {
	//确定模型对象（表）
	m := model.Role{}
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

//Role更新
func RoleUpdate(c *gin.Context) {
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
	m := model.Role{}
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
