package model

import (
	"github.com/jinzhu/gorm"
)

//属性分组表模型

type AttrGroup struct {
	gorm.Model

	AttrTypeID uint
	Name string
	SortOrder int
	//关联
	AttrType AttrType
	Attrs []Attr
	
}
