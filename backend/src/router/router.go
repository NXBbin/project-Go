package router

//路由表

import (
	"controller"

	"middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//路由初始化函数,返回驱动引擎
func Routerlnit() *gin.Engine {
	// 初始化路由引擎对象
	r := gin.Default()

	//为路由引擎增加中间件
	// cors包,允许所有来源请求，解决跨域请求问题
	// r.Use(cors.Default())

	// 允许前端的Authorization请求
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	r.Use(cors.New(config))

	// jwt-token
	// r.Use(JWTToken)

	//定义路由，以及对应的动作处理函数,（访问指定路由时携带Token，判断是否有权限访问）
	bg := r.Group("/", middleware.JWTToken)
	{
		// 查询全部分类：
		bg.GET("category-tree", controller.CategoryTree)
		//添加分类
		bg.POST("category", controller.CategoryAdd)
		//删除分类
		bg.DELETE("category", controller.CategoryDelete)
		// 更新分类
		bg.PUT("category", controller.CategoryUpdate)

		// 产品：
		//查询全部产品
		bg.GET("products", controller.ProductList)
		//删除产品
		bg.DELETE("product", controller.ProductDelete)
		//添加产品
		bg.POST("product", controller.ProductCreate)
		//更新产品
		bg.PUT("product", controller.ProductUpdate)

		//生成路由代码，脚手架模板
		//品牌 Restful 路由
		bg.GET("brand", controller.BrandList)
		bg.DELETE("brand", controller.BrandDelete)
		bg.POST("brand", controller.BrandCreate)
		bg.PUT("brand", controller.BrandUpdate)

		//用户 Restful 路由
		bg.GET("user", controller.UserList)
		bg.DELETE("user", controller.UserDelete)
		bg.POST("user", controller.UserCreate)
		bg.PUT("user", controller.UserUpdate)

		//角色 Restful 路由
		bg.GET("/role", controller.RoleList)
		bg.DELETE("/role", controller.RoleDelete)
		bg.POST("/role", controller.RoleCreate)
		bg.PUT("/role", controller.RoleUpdate)

		//权限 Restful 路由
		bg.GET("/privilege", controller.PrivilegeList)
		bg.DELETE("/privilege", controller.PrivilegeDelete)
		bg.POST("/privilege", controller.PrivilegeCreate)
		bg.PUT("/privilege", controller.PrivilegeUpdate)
	}

	//添加用户状态(登录校验)
	r.POST("/user/auth", controller.UserAuth)

	//返回路由引擎对象
	return r
}
