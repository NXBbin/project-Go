package model

import (
	"github.com/jinzhu/gorm"
)

//分组表模型

type Group struct {
	gorm.Model

	Counter int
	Name string
	SortOrder int

}
