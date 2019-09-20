package model

import (
	"github.com/jinzhu/gorm"
)

//品牌表模型

type Brand struct {
	gorm.Model

	Name string
	Logo string
	Site string
	Description string

}
