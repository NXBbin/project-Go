package model

import (
	"github.com/jinzhu/gorm"
)

//分组表模型

type Group struct {
	gorm.Model

	Counter    int
	Name       string
	SortOrder  int
	AttrTypeID uint

	//前端传递的分组ID
	CheckedProductID []uint `gorm:"-"`
	//前端传递的差异属性ID
	CheckedAttrID []uint `gorm:"-"`
	//关联产品表
	Products []Product
	// 多对多关联
	Attrs []Attr `gorm:"many2many:group_attrs"`
}
