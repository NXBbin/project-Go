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

	//插入数据前先清空表中数据
	// orm.Exec("truncate a_categories")
	// orm.Exec("truncate a_products")

	//-插入数据测试（seed）
	orm.Create(&model.Payment{
		Title: "微信支付", Key: "wechat-pay", Intro: "基于微信提供的支付系统", Status: 1,
	})
	orm.Create(&model.Payment{
		Title: "支付宝", Key: "alipay", Intro: "基于支付宝提供的支付系统", Status: 1,
	})
	orm.Create(&model.Payment{
		Title: "银联", Key: "yinhan-pay", Intro: "基于银联提供的支付系统", Status: 1,
	})
	log.Println("Payment支付方式表数据插入成功")

	orm.Create(&model.PaymentStatus{
		Title: "支付错误",
	})
	orm.Create(&model.PaymentStatus{
		Title: "未支付",
	})
	orm.Create(&model.PaymentStatus{
		Title: "已支付",
	})
	log.Println("PaymentStatus支付状态表数据插入成功")

	orm.Create(&model.Shipping{
		Title: "顺丰", Key: "sf", Intro: "顺丰快递", Status: 1,
	})
	orm.Create(&model.Shipping{
		Title: "圆通", Key: "yt", Intro: "圆通快递", Status: 1,
	})
	orm.Create(&model.Shipping{
		Title: "韵达", Key: "yd", Intro: "韵达快递", Status: 1,
	})
	log.Println("Shipping配送方式表数据插入成功")

	orm.Create(&model.ShippingStatus{
		Title: "超出配送范围",
	})
	orm.Create(&model.ShippingStatus{
		Title: "未发货",
	})
	orm.Create(&model.ShippingStatus{
		Title: "已发货",
	})
	orm.Create(&model.ShippingStatus{
		Title: "已收货",
	})
	log.Println("ShippingStatus配送状态表数据插入成功")

	orm.Create(&model.OrderStatus{
		Title: "订单异常",
	})
	orm.Create(&model.OrderStatus{
		Title: "完成",
	})
	orm.Create(&model.OrderStatus{
		Title: "取消",
	})
	orm.Create(&model.OrderStatus{
		Title: "删除",
	})
	orm.Create(&model.OrderStatus{
		Title: "确认",
	})
	log.Println("OrderStatus订单状态表数据插入成功")

	// //产品-插入数据测试（seed）
	// orm.Create(&model.Product{
	// 	Name:       "纸质书",
	// 	CategoryID: 2,
	// })
	// orm.Create(&model.Product{
	// 	Name:       "电子档",
	// 	CategoryID: 2,
	// })
	// orm.Create(&model.Product{
	// 	Name:       "电脑",
	// 	CategoryID: 3,
	// })
	// orm.Create(&model.Product{
	// 	Name:       "手机",
	// 	CategoryID: 3,
	// })
	// orm.Create(&model.Product{
	// 	Name:       "相机",
	// 	CategoryID: 3,
	// })
	// orm.Create(&model.Product{
	// 	Name:       "平板",
	// 	CategoryID: 3,
	// })
	// orm.Create(&model.Product{
	// 	Name:       "充电宝",
	// 	CategoryID: 3,
	// })

	// orm.Create(&model.User{
	// 	User: "root",
	// })
	log.Println("测试数据生成成功")
}
