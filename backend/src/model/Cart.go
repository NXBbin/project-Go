package model

import (
	"github.com/jinzhu/gorm"
)

//购物车表模型

type Cart struct {
	gorm.Model

	UserID  uint   `gorm:"unique_index"`
	Content string `gorm:"type:text"`
}
