package controller

//wabAPP展示商品列表
import (
	"strings"
	// "config"
	// "time"

	"model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//获取同组产品型号差异信息
func productModel(product *model.Product) {
	// product.ModelInfo = "精装版"
	//如果组存在，
	if product.GroupID == 0 {
		return
	}
	// 先获取组信息,
	orm.Model(&product).Related(&product.Group)
	//再获取产品差异属性
	orm.Model(&product.Group).Related(&product.Group.Attrs, "Attrs")

	//检查是否存在差异属性
	aids := []uint{}
	if len(product.Group.Attrs) > 0 {
		//获取当前产品属性
		for _, a := range product.Group.Attrs {
			aids = append(aids, a.ID)
		}
	}
	//获取product_attr关联表中的差异属性
	pas := []model.ProductAttr{}
	orm.Where("product_id=? AND attr_id = ?", product.ID, aids).Find(&pas)
	values := []string{}
	for _, pa := range pas {
		values = append(values, pa.Value)
	}
	//连接差异属性
	product.ModelInfo = strings.Join(values, "|")
}

//获取产品信息
func ProductInfo(c *gin.Context) {
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
	product := model.Product{}
	product.ID = uint(ID)
	orm.Find(&product)

	//获取关联数据
	//查询关联分类Category
	orm.Model(&product).Related(&product.Category)

	//查询关联图像Images
	product.Images = []model.Image{}
	orm.Model(&product).Related(&product.Images)

	//型号信息
	productModel(&product)

	//查询关联组内产品Group
	orm.Model(&product).Related(&product.Group)
	product.Group.Products = []model.Product{}
	//存在分组的情况
	if product.Group.ID != 0 {
		//查找组内产品信息
		orm.Model(&product.Group).Related(&product.Group.Products)
		for i, _ := range product.Group.Products {
			productModel(&product.Group.Products[i])
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  product,
	})

}

//获取推荐商品主图列表
func ProductPromote(c *gin.Context) {
	//搜索(筛选）
	condStr := ""
	condParams := []string{}
	//获取前端传递的筛选名

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
