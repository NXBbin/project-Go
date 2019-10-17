package model

import (
	"github.com/jinzhu/gorm"
)

//订单状态表模型

type OrderStatus struct {
	gorm.Model

	Title string
}
