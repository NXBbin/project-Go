package model

import (
	"github.com/jinzhu/gorm"
)

//角色表模型

type Role struct {
	gorm.Model

	Name string
	SortOrder int
	Description string
	
	// 多对多关联
	Privileges []Privilege `gorm:"many2many:role_privileges;"`

}
