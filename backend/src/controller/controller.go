package controller

//连接数据库配置信息

import (
	"time"
	"config"
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

//基础控制器代码
var orm *gorm.DB

//初始化MySQL
func InitDB() (*gorm.DB, error) {
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
	db, dberr := gorm.Open(config.App["DB_DRIVER"], dsn)
	if dberr != nil {
		log.Println(dberr)
		return nil, dberr
	}
	orm = db
	return db, nil
}

//初始化redis
var rds redis.Conn

func InitRedis() (redis.Conn, error) {
	//设置选项
	options := []redis.DialOption{}
	//设置1秒超时时间
	options = append(options,redis.DialConnectTimeout(1*time.Second))
	
	if db, ok := config.App["REDIS_DB"]; ok {
		dbInt, _ := strconv.Atoi(db)
		options = append(options, redis.DialDatabase(dbInt))
	}
	if pwd, ok := config.App["REDIS_PASSWORD"]; ok && pwd != "" {
		options = append(options, redis.DialPassword(pwd))
	}

	//建立连接
	c, err := redis.Dial("tcp", config.App["REDIS_HOST"]+":"+config.App["REDIS_PORT"], options...)
	if err != nil {
		return nil, err
	}

	rds = c

	return c, nil
}
