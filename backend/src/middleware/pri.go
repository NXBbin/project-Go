package middleware

import (
	"config"
	"fmt"
	"log"
	"model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//获取用户权限中间件
func Pri(c *gin.Context) {

	// 初始化Gorm，处理特殊表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.App["DB_TABLE_PREFIX"] + defaultTableName
	}

	//基于模型 连接数据库
	// "bin:123456@tcp(localhost:3306)/projecta?charset=utf8mb4&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s&parseTime=%s",
		config.App["MYSQL_USER"],
		config.App["MYSQL_PASSWORD"],
		config.App["MYSQL_HOST"],
		config.App["MYSQL_PORT"],
		config.App["MYSQL_DBNAME"],
		config.App["MYSQL_CHARSET"],
		config.App["MYSQL_LOC"],
		config.App["MYSQL_PARSETIME"],
	)
	orm, dberr := gorm.Open(config.App["DB_DRIVER"], dsn)
	if dberr != nil {
		log.Println(dberr)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//利用userName得到用户的权限列表
	userName, _ := c.Get("userName")
	user := model.User{}
	if orm.Where("user=?", userName).Find(&user).RecordNotFound() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "该用户不存在",
		})
		return
	}
	pris := []string{}

	//利用user得到权限列表
	role := model.Role{}
	orm.Model(&user).Related(&role, "Role")
	if role.ID != 0 {
		//利用角色获取权限
		cps := []model.Privilege{}
		orm.Model(&role).Related(&cps, "Privileges")
		//拼凑权限Key的切片
		for _, p := range cps {
			pris = append(pris, p.Key)
		}
	}
	// log.Println(pris)
	// 记录当前用户的权限
	c.Set("pris", pris)
}
