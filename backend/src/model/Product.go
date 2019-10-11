package model

//产品表模型

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

	//解析属性值字段,忽略该字段不在表中创建该字段
	AttrValue map[uint]string `gorm:"-"`
	//上传图像
	UploadedImage      []string `gorm:"-"`
	UploadedImageSmall []string `gorm:"-"`
	UploadedImageBig   []string `gorm:"-"`

	//关联定义，多个产品关联一个分类
	//被关联的外键（必须存在一个...ID属性名，否则关联的外键上需要使用targ语法)
	CategoryID uint
	AttrTypeID uint
	GroupID    uint

	//产品属于分类
	Category Category

	// 属性类型
	AttrType AttrType
	//产品属性
	ProductAttrs []ProductAttr
	//关联图像表
	Images []Image
	//关联分组产品表
	Group Group
	//webAPP同组产品差异信息额外字段
	ModelInfo string `gorm:"-"`
}
