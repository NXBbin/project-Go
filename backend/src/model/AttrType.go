package model

import (
	"github.com/jinzhu/gorm"
)

//属性类型表模型

type AttrType struct {
	gorm.Model

	Name      string
	SortOrder int

	//关联属性分组
	AttrGroups []AttrGroup
}
