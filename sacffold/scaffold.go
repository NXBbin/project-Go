package main

import (
	"bufio"
	// "config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	// "model"
	"os"
	"regexp"
	"strings"

	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
)

//脚手架(代码生成器)

//全局变量
var jsonName, mName, mTitle, rName, cName string
var tConfig map[string]interface{}
var tFields []map[string]string

func main() {
	//获取命令行传递的参数
	jsonName = os.Args[1]

	//解析JSON
	parseJSON()

	// 模型名
	mName = tConfig["modelName"].(string)
	// 标题
	mTitle = tConfig["title"].(string)
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
	// genTable()

	//模型代码生成
	genModel()

	// 生成路由代码
	genRouter()

	//生成控制器代码
	genController()

	// # 5 生成前端路由
	genManagerRouter()
	// # 6 生成前端组件
	genManagerCP()

}

//生成前端组件
func genManagerCP() {
	// 从模板读取内容
	tplFile := "./scaffoldTemplate/manager/listComponent"
	content, _ := ioutil.ReadFile(tplFile) // []byte
	code := string(content)
	// 占位符替换
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%r-name%", rName)
	cpName := mName + "List"
	code = strings.ReplaceAll(code, "%cp-name%", cpName)
	// 展示列列表，筛选字段
	tcList := ""
	filterList := ""
	setFieldList := ""
	for _, field := range tFields { //map[name:Name type:string, isList:yes, label:]
		if field["isList"] == "yes" {
			// 需要列表展示的情况下
			sortAble := ""
			if field["isSort"] == "yes" {
				sortAble = `sortable="custom"`
			}
			tcList += fmt.Sprintf("\t\t\t\t\t\t\t\t<el-table-column prop=\"%s\" label=\"%s\" %s></el-table-column>\n",
				field["name"], field["label"], sortAble,
			)
		}

		if field["isFilter"] == "yes" {
			//需要筛选
			filterList += fmt.Sprintf(`
						<el-form-item label="%s">
							<el-input v-model="filterForm.filter%s"></el-input>
						</el-form-item>`,
				field["label"], field["name"],
			)
		}

		if field["isSet"] == "yes" {
			// 需要出现在表单中
			setFieldList += fmt.Sprintf(`
						<el-form-item label="%s" prop="%s">
							<el-input v-model="itemSetForm.%s"></el-input>
						</el-form-item>`,
				field["label"], field["name"], field["name"],
			)
		}

	}
	code = strings.ReplaceAll(code, "%tc-list%", tcList)
	code = strings.ReplaceAll(code, "%filter-list%", filterList)
	code = strings.ReplaceAll(code, "%set-field-list%", setFieldList)

	// 将 code 写入到指定文件
	cPath := strings.ToLower(mName[0:1]) + mName[1:]
	codeFile := "../manager/src/components/" + cPath + "/" + cpName + ".vue"
	// 子目录
	os.Mkdir("../manager/src/components/"+cPath, 0755)
	handle, _ := os.OpenFile(codeFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("生成的代码位于 components/ 中")
}

// 生成前端路由代码
func genManagerRouter() {
	code := `
	{ path: '%r-name%', component: ()=>import('../components/%c-path%/%c-name%List.vue'), },`
	// 占位符替换
	code = strings.ReplaceAll(code, "%r-name%", rName)
	cPath := strings.ToLower(mName[0:1]) + mName[1:]
	code = strings.ReplaceAll(code, "%c-path%", cPath)
	code = strings.ReplaceAll(code, "%c-name%", cName)

	codeFile := "./ManagerRouterCode"
	handle, _ := os.OpenFile(codeFile, os.O_APPEND|os.O_CREATE, 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("前端路由代码已生成，位于ManagerRouterCode中，请将代码拷贝 ")

}

//生成路由代码
func genRouter() {
	// 读取模板内容
	tplFile := "./scaffoldTemplate/backend/router"
	content, _ := ioutil.ReadFile(tplFile)
	//将读取到的字节型切片数据转成字符串
	code := string(content)
	//占位符替换
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%r-name%", rName)
	code = strings.ReplaceAll(code, "%c-name%", cName)
	//将数据写入到指定文件中
	codeFile := "./routerCode"
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
	tplFile := "./scaffoldTemplate/backend/controller"
	content, _ := ioutil.ReadFile(tplFile)
	code := string(content)
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%m-name%", mName)
	code = strings.ReplaceAll(code, "%c-name%", cName)
	//统一首字母小写
	cfName := strings.ToLower(mName[0:1]) + mName[1:]
	codeFile := "../backend/src/controller/" + cfName + ".go"
	//每次打开清空文件，重新写入
	handle, _ := os.OpenFile(codeFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()

	log.Println("控制器代码已生成")
}

//解析JSON文件
func parseJSON() {
	// 读取json文件
	file := "./config/" + jsonName + ".json"
	content, _ := ioutil.ReadFile(file)
	//解码内容，到map里
	json.Unmarshal(content, &tConfig)
	// log.Println(tConfig)
	//遍历模板的字段信息断言后存入变量
	fields := tConfig["fields"].([]interface{})
	for _, f := range fields {
		field := map[string]string{}
		for k, v := range f.(map[string]interface{}) {
			field[k] = v.(string)
		}
		tFields = append(tFields, field)
	}
	// log.Println(fields)
}

//模型代码生成
func genModel() {
	// 读取模板内容
	tplFile := "./scaffoldTemplate/backend/model"
	content, _ := ioutil.ReadFile(tplFile) //得到[]byte
	code := string(content)
	// 替换占位符
	code = strings.ReplaceAll(code, "%m-title%", mTitle)
	code = strings.ReplaceAll(code, "%m-name%", mName)
	//替换field-list
	fieldList := ""
	for _, field := range tFields {
		fieldList += fmt.Sprintf("\t%s %s\n", field["name"], field["type"])
	}
	code = strings.ReplaceAll(code, "%field-list%", fieldList)
	// log.Println(code)

	// 将替换完成的代码，写入文件
	codeFile := "D:/demo/GO-project/project-Go/backend/src/model/" + mName + ".go"
	handle, _ := os.OpenFile(codeFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0)
	defer handle.Close()
	writer := bufio.NewWriter(handle)
	writer.WriteString(code)
	writer.Flush()
	log.Println("模型代码已生成")
}

//模型构建表
// func genTable() {
// 	//1.启用配置
// 	config.InitConfig()

// 	//初始化GORM
// 	// 初始化Gorm，处理特殊表名前缀
// 	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
// 		return config.App["DB_TABLE_PREFIX"] + defaultTableName
// 	}

// 	//基于模型 拼接连接数据库配置信息
// 	// "bin:123456@tcp(localhost:3306)/projecta?charset=utf8mb4&loc=Local"
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s",
// 		config.App["MYSQL_USER"],
// 		config.App["MYSQL_POSSWORD"],
// 		config.App["MYSQL_HOST"],
// 		config.App["MYSQL_PORT"],
// 		config.App["MYSQL_DBNAME"],
// 		config.App["MYSQL_CHARSET"],
// 		config.App["MYSQL_LOC"],
// 	)
// 	//2.连接数据库
// 	orm, dberr := gorm.Open(config.App["DB_DRIVER"], dsn)
// 	if dberr != nil {
// 		log.Println(dberr)
// 		return
// 	}

// //迁移（利用模型创建表，migrate）
// orm.AutoMigrate(&model.Brand{})
// log.Println("brand表已创建")

// }
