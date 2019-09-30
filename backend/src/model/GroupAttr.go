package model

import (
	"github.com/jinzhu/gorm"
)

//产品分组差异属性表模型

type GroupAttr struct {
	gorm.Model

	GroupID uint
	AttrID   uint
}
