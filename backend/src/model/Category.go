package model

import "github.com/jinzhu/gorm"

//模型结构体定义，对应数据库中的表结构
type Category struct {
	//嵌套结构体(创建，更新，删除时间)
	gorm.Model

	ParentId        uint
	Name            string
	Logo            string
	Description     string
	SortOrder       int
	MetaTitle       string
	MetaKeywords    string
	MetaDescription string
	//外键 (拥有多个：has many),类型为模型切片
	Products []Product
	//外键 (拥有1个：has one),类型为模型切片
	//belongs to(属于某个
	//many to many
}
