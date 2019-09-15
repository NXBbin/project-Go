package main

//利用模型插入产品数据

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

	//插入数据前先清空表中数据
	orm.Exec("truncate a_categories")
	orm.Exec("truncate a_products")

	//分类-插入数据测试（seed）
	orm.Create(&model.Category{
		Name:     "未分类",
		ParentId: 0,
	})
	orm.Create(&model.Category{
		Name:     "图书",
		ParentId: 0,
	})
	orm.Create(&model.Category{
		Name:     "数码产品",
		ParentId: 0,
	})

	//产品-插入数据测试（seed）
	orm.Create(&model.Product{
		Name:       "纸质书",
		CategoryID: 2,
	})
	orm.Create(&model.Product{
		Name:       "电子档",
		CategoryID: 2,
	})
	orm.Create(&model.Product{
		Name:       "电脑",
		CategoryID: 3,
	})
	orm.Create(&model.Product{
		Name:       "手机",
		CategoryID: 3,
	})
	// log.Println(orm.Error)
}
