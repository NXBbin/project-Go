package main

//利用模型生成产品表

import (
	"config"
	"fmt"
	"log"
	"model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	//1.启用配置
	config.InitConfig()

	//初始化GORM
	// 初始化Gorm，处理特殊表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.App["DB_TABLE_PREFIX"] + defaultTableName
	}

	//基于模型 拼接连接数据库配置信息
	// "bin:123456@tcp(localhost:3306)/projecta?charset=utf8mb4&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s",
		config.App["MYSQL_USER"],
		config.App["MYSQL_PASSWORD"],
		config.App["MYSQL_HOST"],
		config.App["MYSQL_PORT"],
		config.App["MYSQL_DBNAME"],
		config.App["MYSQL_CHARSET"],
		config.App["MYSQL_LOC"],
	)
	//2.连接数据库
	orm, dberr := gorm.Open(config.App["DB_DRIVER"], dsn)
	if dberr != nil {
		log.Println(dberr)
		return
	}

	//3.迁移
	//判断表是否已经存在
	// if orm.HasTable(&model.Product{}) {
	// 	// 删除表
	// 	orm.DropTable(&model.Product{})
	// 	log.Println("Product已删除，并重新创建")
	// }

	// //判断表是否已经存在
	// if orm.HasTable(&model.Category{}) {
	// 	//删除表
	// 	orm.DropTable(&model.Category{})
	// 	log.Println("Category已删除，并重新创建")
	// }

	// //迁移（利用模型创建表，migrate）
	orm.AutoMigrate(&model.Product{})
	log.Println("Product已创建")

	//创建categories产品表
	orm.AutoMigrate(&model.Category{})
	log.Println("Category已创建")

	// 创建categories产品表
	orm.AutoMigrate(&model.User{})
	log.Println("User已创建")

	//创建角色表
	orm.AutoMigrate(&model.Role{})
	log.Println("Role已创建")

	// 创建角色权限表
	orm.AutoMigrate(&model.Privilege{})
	log.Println("Privilege已创建")

	//创建角色授权关联表
	orm.AutoMigrate(&model.RolePrivilege{})
	log.Println("RolePrivilege已创建")

	//创建属性类型表
	orm.AutoMigrate(&model.AttrType{})
	log.Println("AttrType属性类型表已创建")

	// 创建属性分组表
	orm.AutoMigrate(&model.AttrGroup{})
	log.Println("AttrGroup属性分组表已创建")

	//创建属性表
	orm.AutoMigrate(&model.Attr{})
	log.Println("Attr属性表已创建")

	//创建产品属性表
	orm.AutoMigrate(&model.ProductAttr{})
	log.Println("ProductAttr产品属性表已创建")

	//创建分组表
	orm.AutoMigrate(&model.Group{})
	log.Println("Group分组表已创建")

	//创建分组表
	orm.AutoMigrate(&model.GroupAttr{})
	log.Println("GroupAttr分组差异属性表已创建")

	//创建产品图像表
	orm.AutoMigrate(&model.Image{})
	log.Println("Image产品图像表已创建")

	//创建购物车表
	orm.AutoMigrate(&model.Cart{})
	log.Println("Cart购物车表已创建")

	//创建收货地址表
	orm.AutoMigrate(&model.Address{})
	log.Println("Address收货地址表已创建")

	//创建订单表
	orm.AutoMigrate(&model.Order{})
	log.Println("Order订单表已创建")

	//创建订单状态表
	orm.AutoMigrate(&model.OrderStatus{})
	log.Println("Order-Status订单状态表已创建")

	//创建支付表
	orm.AutoMigrate(&model.Payment{})
	log.Println("Payment支付表已创建")

	//创建支付状态表
	orm.AutoMigrate(&model.PaymentStatus{})
	log.Println("Payment-Status支付状态表已创建")

	//创建配送表
	orm.AutoMigrate(&model.Shipping{})
	log.Println("Shipping配送表已创建")

	//创建配送状态表
	orm.AutoMigrate(&model.ShippingStatus{})
	log.Println("Shipping-Status配送状态表已创建")

	//创建订单产品关联表
	orm.AutoMigrate(&model.OrderProduct{})
	log.Println("OrderProduct订单产品关联表已创建")
}
