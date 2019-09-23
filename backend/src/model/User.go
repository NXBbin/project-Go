package model

import (
	"github.com/jinzhu/gorm"
)

//用户表模型

type User struct {
	gorm.Model

	User         string
	Email        string
	Tel          string
	Password     string
	PasswordSalt string
	Token        string
	//关联字段
	RoleID uint
	Role Role
}
