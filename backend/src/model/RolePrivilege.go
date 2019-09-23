package model

import (
	"github.com/jinzhu/gorm"
)

//角色授权表模型

type RolePrivilege struct {
	gorm.Model

	RoleID uint
	PrivilegeID uint

}
