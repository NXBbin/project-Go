package model

import (
	"github.com/jinzhu/gorm"
)

//品牌表模型

type Brand struct {
	gorm.Model

	Name string
	Logo string
	//官网
	Site string
	//介绍
	Description string
}
