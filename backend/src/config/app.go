package config

//数据库配置信息集合

var App map[string]string

func InitConfig() {
	App = map[string]string{
		"DB_DRIVER":       "mysql",
		"DB_TABLE_PREFIX": "a_",
		"MYSQL_HOST":      "localhost",
		"MYSQL_PORT":      "3306",
		"MYSQL_USER":      "bin",
		"MYSQL_PASSWORD":  "123456",
		"MYSQL_DBNAME":    "projecta",
		"MYSQL_CHARSET":   "utf8mb4",
		"MYSQL_LOC":       "Local",
		"MYSQL_PARSETIME": "true", //false默认，不执行time的解析
		//redis配置
		"REDIS_HOST":     "192.168.150.200",
		"REDIS_PORT":     "6379",
		"REDIS_DB":       "0",
		"REDIS_PASSWORD": "",

		"SERVER_ADDR": ":8088",
		"SECRET":      "AllYourBase", //Token签名key随机串
		//存储上传文件的路径
		"UPLOAD_PATH": "D:\\demo\\GO-project\\project-Go\\upload\\",
		//静态资源服务器
		"IMAGE_HOST": "http://localhost:8089/",
		// 小图缩放值
		"THUMB_SMALL_W": "146",
		"THUMB_SMALL_H": "146",
		//中图缩放值
		"THUMB_BIG_W": "1460",
		"THUMB_BIG_H": "1460",
	}
}
