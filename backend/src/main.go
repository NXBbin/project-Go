package main

import (
	"log"
	"router"
	"config"
	// "github.com/jinzhu/gorm"
	"controller"
)

func main() {
	//初始化配置
	config.InitConfig()
	
	//调用初始化路由函数，获得路由对象
	r := router.Routerlnit()

	// 调用初始化方法，（连接数据库，gorm）
	db, err := controller.InitDB()
	if err != nil {
		log.Println("数据库连接失败", err)
		return
	}
	defer db.Close()

	//启动服务端口
	r.Run(config.App["SERVER_ADDR"])
}
