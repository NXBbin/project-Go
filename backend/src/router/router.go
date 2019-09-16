package router

import (
	"controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//路由初始化函数,返回驱动引擎
func Routerlnit() *gin.Engine {
	// 初始化路由引擎对象
	r := gin.Default()

	//为路由引擎增加中间件
	// cors包,允许所有来源请求，解决跨域请求问题
	r.Use(cors.Default())

	// 分类：
	//定义路由，以及对应的动作处理函数
	r.GET("/ping", controller.Ping)
	r.GET("/category-tree", controller.CategoryTree)
	//添加分类
	r.POST("/category", controller.CategoryAdd)
	//删除分类
	r.DELETE("/category", controller.CategoryDelete)
	// 更新分类
	r.PUT("/category", controller.CategoryUpdate)

	// 产品：
	//获取产品
	r.GET("/products",controller.ProductList)
	//删除产品
	r.DELETE("/product",controller.ProductDelete)
	
	//返回路由引擎对象
	return r
}
