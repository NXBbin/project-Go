package controller

//wabAPP购物车列表
import (
	// "strings"
	// "config"
	// "time"

	"model"
	"net/http"

	// "strconv"

	"github.com/gin-gonic/gin"
)

//获取商品主信息列表
func CartProduct(c *gin.Context) {
	//获取前端传递的filterIDs[]（由于前端传递的是数组，在Headre里面中括号也是key的一部分。）
	filterIDs := c.QueryArray("filterIDs[]")
	if len(filterIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "产品ID不存在",
		})
		return
	}

	//条件
	condStr := "id in (?)"

	//获取product模型
	products := []model.Product{}
	//获取展示数量和偏移量,输出数据获 // select * from a_products where id in (44,22,51);
	orm.Where(condStr, filterIDs).Find(&products)
	//遍历全部属性，找到关联字段
	for i, _ := range products {
		//查询关联图像
		products[i].Images = []model.Image{}
		orm.Model(&products[i]).Related(&products[i].Images)
		//获取产品型号信息
		productModel(&products[i])
	}

	//响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  products,
	})
}
