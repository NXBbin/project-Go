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
	bg := r.Group("/", middleware.JWTToken, middleware.Pri)
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
		bg.GET("products", middleware.Require("product-manager"), controller.ProductList)
		//删除产品
		bg.DELETE("product", middleware.Require("product-manager"), controller.ProductDelete)
		//添加产品
		bg.POST("product", middleware.Require("product-manager"), controller.ProductCreate)
		//更新产品
		bg.PUT("product", middleware.Require("product-manager"), controller.ProductUpdate)
		//复制
		bg.POST("product-copy", middleware.Require("product-manager"), controller.ProductCopy)

		//生成路由代码，脚手架模板
		//品牌 Restful 路由
		bg.GET("brand", middleware.Require("brand-manager"), controller.BrandList)
		bg.DELETE("brand", middleware.Require("brand-manager"), controller.BrandDelete)
		bg.POST("brand", middleware.Require("brand-manager"), controller.BrandCreate)
		bg.PUT("brand", middleware.Require("brand-manager"), controller.BrandUpdate)

		//用户 Restful 路由
		bg.GET("user", middleware.Require("user-manager"), controller.UserList)
		bg.DELETE("user", middleware.Require("user-manager"), controller.UserDelete)
		bg.POST("user", middleware.Require("user-manager"), controller.UserCreate)
		bg.PUT("user", middleware.Require("user-manager"), controller.UserUpdate)

		//角色 Restful 路由
		bg.GET("/role", middleware.Require("role-manager"), controller.RoleList)
		bg.DELETE("/role", middleware.Require("role-manager"), controller.RoleDelete)
		bg.POST("/role", middleware.Require("role-manager"), controller.RoleCreate)
		bg.PUT("/role", middleware.Require("role-manager"), controller.RoleUpdate)
		//展示全部权限信息
		bg.GET("/privilege-role", middleware.Require("role-grant"), controller.PrivilegeRole)
		//展示已授权信息
		bg.PUT("/role-grant", middleware.Require("role-grant"), controller.RoleGrant)

		//权限 Restful 路由
		bg.GET("/privilege", controller.PrivilegeList)
		bg.DELETE("/privilege", controller.PrivilegeDelete)
		bg.POST("/privilege", controller.PrivilegeCreate)
		bg.PUT("/privilege", controller.PrivilegeUpdate)

		//属性类型 Restful 路由
		bg.GET("/attr-type", controller.AttrTypeList)
		bg.DELETE("/attr-type", controller.AttrTypeDelete)
		bg.POST("/attr-type", controller.AttrTypeCreate)
		bg.PUT("/attr-type", controller.AttrTypeUpdate)

		//属性分组 Restful 路由
		bg.GET("/attr-group", controller.AttrGroupList)
		bg.DELETE("/attr-group", controller.AttrGroupDelete)
		bg.POST("/attr-group", controller.AttrGroupCreate)
		bg.PUT("/attr-group", controller.AttrGroupUpdate)

		//属性 Restful 路由
		bg.GET("/attr", controller.AttrList)
		bg.DELETE("/attr", controller.AttrDelete)
		bg.POST("/attr", controller.AttrCreate)
		bg.PUT("/attr", controller.AttrUpdate)

		//产品属性 Restful 路由
		bg.GET("/product-attr", controller.ProductAttrList)
		bg.DELETE("/product-attr", controller.ProductAttrDelete)
		bg.POST("/product-attr", controller.ProductAttrCreate)
		bg.PUT("/product-attr", controller.ProductAttrUpdate)

		//分组 Restful 路由
		bg.GET("/group", controller.GroupList)
		bg.DELETE("/group", controller.GroupDelete)
		bg.POST("/group", controller.GroupCreate)
		bg.PUT("/group", controller.GroupUpdate)

		//图像上传
		bg.POST("/image-upload", controller.ImageUpload)
	}

	//添加用户状态(登录校验)
	r.POST("/user/auth", controller.UserAuth)

	//webAPP项目路由
	//商品页
	r.GET("product-promote", controller.ProductPromote)
	//商品型号选择
	r.GET("product-info", controller.ProductInfo)
	//购物车
	r.GET("cart-product", controller.CartProduct)
	//用户登录
	r.POST("member-login", controller.MemberLogin)
	//认证会员
	r.GET("member-auth", controller.MemberAuth)
	//获取会员购物车信息
	r.GET("member-cart", controller.MemberCart)
	//将前端购物车信息同步到后端
	r.PUT("member-cart-sync", controller.MemberCartSync)
	//将后端数据同步前端
	r.PUT("member-cart-set", controller.MemberCartSet)
	//获取收货地址列表
	r.GET("member-address-list", controller.MemberAddressList)
	//新增地址
	r.POST("member-address-add", controller.MemberAddressAdd)
	//配送列表
	r.GET("shipping-list", controller.ShippingList)
	//创建订单号
	r.POST("order-create", controller.OrderCreate)
	//处理订单
	r.GET("order-result", controller.OrderResult)

	//验证码
	r.GET("check-code", controller.CheckCode)

	//返回路由引擎对象
	return r
}
