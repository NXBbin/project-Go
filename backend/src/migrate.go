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
		config.App["MYSQL_POSSWORD"],
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

	//迁移（利用模型创建表，migrate）
	orm.AutoMigrate(&model.Product{})
	log.Println("Product已创建")

	//创建categories产品表
	orm.AutoMigrate(&model.Category{})
	log.Println("Category已创建")

	// log.Println(orm.Error)
}
