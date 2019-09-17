package main

import (
	"bufio"
	"config"
	"fmt"
	"io/ioutil"
	"log"
	"model"
	"os"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//脚手架(代码生成器)

//全局变量
var mName, mTitle, rName, cName string

func main() {
	//获取命令行传递的参数
	mName = os.Args[1]
	mTitle = os.Args[2]
	// log.Println(mName, mTitle)

	// 利用参数替换路由
	//定义正则表达式规则
	p := `([A-Z])`
	re, _ := regexp.Compile(p)
	//正则替换 | 转小写 | 去首-: (传 BrandRoom 生成 brand-room)
	rName = strings.ToLower(re.ReplaceAllString(mName, "-${1}"))[1:]
	//控制器名字前缀
	cName = mName
	// log.Println(rName, cName)

	//模型构建表
	genTable()

	// 生成路由代码
	genRouter()

	//生成控制器代码
	genController()

}

//生成路由代码
func genRouter() {
	// 读取模板内容
	tplFile := "./scaffoldTemplate/router"
	content, _ := ioutil.ReadFile(tplFile)
	//将读取到的字节型切片数据转成字符串
	code := string(content)
	//占位符替换
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%r-name%", rName)
	code = strings.ReplaceAll(code, "%c-name%", cName)
	//将数据写入到指定文件中
	codeFile := "./scaffoldTemplate/routerCode"
	//打开文件追加数据，不存在则创建
	handle, _ := os.OpenFile(codeFile, os.O_APPEND|os.O_CREATE, 0)
	defer handle.Close()
	//获取写入缓冲器
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("路由代码已生成请复制")
}

//生成控制器代码
func genController() {
	// 读取模板内容
	tplFile := "./scaffoldTemplate/controller"
	content, _ := ioutil.ReadFile(tplFile)
	code := string(content)
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%m-name%", mName)
	code = strings.ReplaceAll(code, "%c-name%", cName)
	//统一首字母小写
	cfName := strings.ToLower(mName[0:1]) + mName[1:]
	codeFile := "./controller/" + cfName + ".go"
	//每次打开清空文件，重新写入
	handle, _ := os.OpenFile(codeFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("控制器代码已生成")
}

//模型构建表
func genTable() {
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

	//迁移（利用模型创建表，migrate）
	orm.AutoMigrate(&model.Brand{})
	log.Println("brand表已创建")

}
