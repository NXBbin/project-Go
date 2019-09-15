package controller

//负责 分类相关操作的函数集合文件

import (
	"strconv"

	"github.com/gin-gonic/gin"

	// "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	// "dao"
	"log"
	"model"
)

//分类树
func CategoryTree(c *gin.Context) {

	//获取分类结构体对象
	categories := []model.Category{}
	//获取全部分类
	orm.Find(&categories)
	//遍历categors，得到每个分类，利用分类查询关联
	for i, v := range categories {
		//利用orm模型操作查询表中全部字段，找到表中的关联字段
		orm.Model(&categories[i]).Related(&categories[i].Products)
		log.Println(i, v)
	}
	log.Println(&categories)
	//响应数据
	c.JSON(200, gin.H{
		"error": "",
		"data":  categories,
	})

}

//添加分类
func CategoryAdd(c *gin.Context) {

	//得到模型对象（表）
	category := model.Category{}
	//使用c.ShouldBind(),绑定并解析post数据
	err := c.ShouldBind(&category)
	if err != nil {
		//响应错误
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	} else {
		// 解析数据成功
		log.Println("数据解析成功：", category)
		// 校验数据

		// 数据入库
		orm.Create(&category)
		log.Println("数据入库完成：", category)
		//完成响应
		c.JSON(200, gin.H{
			"error": "",
			"data":  category,
		})
	}
}

// 删除分类
func CategoryDelete(c *gin.Context) {

	// 获取删除的ID
	ID := c.Query("ID")
	// 构建模型对象
	category := model.Category{}
	//将前端传递过来的字符串型的ID转成int整型
	id, _ := strconv.Atoi(ID)
	category.ID = uint(id)
	//删除
	orm.Delete(&category)
	//响应
	c.JSON(200, gin.H{
		"error": "",
	})
}

// 更新分类
func CategoryUpdate(c *gin.Context) {

	// 获取更新的ID
	ID := c.Query("ID")
	//得到模型对象（表）
	category := model.Category{}
	//将前端传递过来的字符串型的ID转成int整型
	id, _ := strconv.Atoi(ID)
	category.ID = uint(id)
	//使用c.ShouldBind(),绑定并解析post数据
	err := c.ShouldBind(&category)
	if err != nil {
		//响应错误
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	//更新数据
	orm.Save(&category)
	//响应
	c.JSON(200, gin.H{
		"error": "",
		"data":  category,
	})

}

// //分类树
// func CategoryTree(c *gin.Context) {
//连接数据库，获取全部分类内容
//数据库参数
// config := map[string]string{
// 	"username":  "bin",
// 	"password":  "123456",
// 	"host":      "127.0.0.1",
// 	"port":      "3306",
// 	"dbname":    "projecta",
// 	"collation": "utf8mb4_general_ci",
// }
// // 调用构造函数，并且传递参数，进行初始化工作
// db, err := dao.NewDao(config)
// if err != nil {
// 	log.Println("初始化失败：", err)
// 	return
// }
// // 查询分类的全部数据
// rows, err := db.Table("a_categories").Rows()
// log.Println(rows, err)
// //没有错误的情况下，响应结果
// if err == nil {
// 	c.JSON(200, gin.H{
// 		"error": "",
// 		"data":  rows,
// 	})
// } else {
// 	c.JSON(200, gin.H{
// 		"error": err.Error(),
// 	})
// }

// }
