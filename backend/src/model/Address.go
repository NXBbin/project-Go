package model

import (
	"github.com/jinzhu/gorm"
)

//收货地址表模型

type Address struct {
	gorm.Model

	UserID       uint
	Tag          string //标签（家，公司）
	ProvinceCode string `gorm:"type:char(6)"` //前端编码
	Province     string //省
	CityCode     string `gorm:"type:char(6)"` //前端编码
	City         string //市
	//解析JSON格式的前端编码，form来自该键
	CountyCode string `gorm:"type:char(6)" form:"areaCode" json:"areaCode"`
	County     string //区

	Addr string `form:"addressDetail" json:"addressDetail"` //详细地址

	IsDefault bool   //默认地址
	Name      string //收货人
	Tel       string //联系电话
	PostCode  string `form:"postalCode" json:"postalCode"` //邮政编码
}
