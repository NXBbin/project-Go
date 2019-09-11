package router

import (
	"controller"

	"github.com/gin-gonic/gin"
)

//路由初始化函数,返回驱动引擎
func Routerlnit() *gin.Engine {
	// 初始化路由引擎对象
	r := gin.Default()

	//定义路由，以及对应的动作处理函数
	r.GET("/ping", controller.Ping)
	r.GET("/category-tree", controller.CategoryTree)

	//返回路由引擎对象
	return r
}
