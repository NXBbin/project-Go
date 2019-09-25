package model

import (
	"github.com/jinzhu/gorm"
)

//属性表模型

type Attr struct {
	gorm.Model

	AttrGroupID uint
	Name        string
	SortOrder   int
}
