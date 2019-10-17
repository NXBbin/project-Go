package model

import (
	"github.com/jinzhu/gorm"
)

//配送状态表模型

type ShippingStatus struct {
	gorm.Model

	Title string
}
