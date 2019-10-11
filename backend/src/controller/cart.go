package controller

//wabAPP展示商品列表
import (
	// "strings"
	// "config"
	// "time"

	"model"
	"net/http"

	// "strconv"

	"github.com/gin-gonic/gin"
)

//获取推荐商品主图列表
func CartProduct(c *gin.Context) {
	//获取前端传递的filterIDs[]
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
	//获取展示数量和偏移量,输出数据获
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
