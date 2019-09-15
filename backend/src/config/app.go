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
		"MYSQL_POSSWORD":  "123456",
		"MYSQL_DBNAME":    "projecta",
		"MYSQL_CHARSET":   "utf8mb4",
		"MYSQL_LOC":       "Local",
		"MYSQL_PARSETIME": "true", //false默认，不执行time的解析
		"SERVER_ADDR":     ":8088",
	}
}
