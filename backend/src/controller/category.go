package controller

import (
	"dao"
	"log"

	"github.com/gin-gonic/gin"
)

//负责 分类相关操作的函数集合文件

//分类树
func CategoryTree(c *gin.Context) {
	//连接数据库，获取全部分类内容
	//数据库参数
	config := map[string]string{
		"username":  "bin",
		"password":  "123456",
		"host":      "127.0.0.1",
		"port":      "3306",
		"dbname":    "projecta",
		"collation": "utf8mb4_general_ci",
	}
	// 调用构造函数，并且传递参数，进行初始化工作
	db, err := dao.NewDao(config)
	if err != nil {
		log.Println("初始化失败：", err)
		return
	}
	// 查询分类的全部数据
	rows, err := db.Table("a_categories").Rows()
	log.Println(rows, err)
	//没有错误的情况下，响应结果
	if err == nil {
		c.JSON(200, gin.H{
			"error": "",
			"data":  rows,
		})
	} else {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	}

}

// 删除分类
func CategoryDelete() {

}
