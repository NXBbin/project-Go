package model

//产品模型

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model

	//targ注释语法
	Name        string    `gorm:"index;"`
	Price       float64   `gorm:"type:decimal(14,2);index;"`
	Upc         string    `gorm:"index;"`
	Mpn         string    `gorm:"size:127;"`
	IsSale      int       `gorm:"是否在售"`
	SaleTime    time.Time `gorm:"起售时间"`
	IsSubstract int       `gorm:"是否扣减库存"`
	IsShipping  int       `gorm:"是否支持配送"`
	Weight      float64   `gorm:"重量"`
	Description string    `gorm:"描述"`

	//关联定义，多个产品关联一个分类
	//被关联的外键（必须存在一个...ID属性名，否则关联的外键上需要使用targ语法)
	CategoryID uint
	
	//产品属于分类
	Category Category
}
