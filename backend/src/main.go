package main

import (
	"router"

	"github.com/gin-contrib/cors"
)

func main() {
	//调用初始化路由函数，获得路由对象
	r := router.Routerlnit()

	//为路由引擎增加中间件
	// CORS,允许所有来源请求，解决跨域请求问题
	r.Use(cors.Default())

	//启动服务端口
	r.Run(":8088")
}
