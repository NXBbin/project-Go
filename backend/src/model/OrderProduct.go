package model

import (
	"github.com/jinzhu/gorm"
)

//订单产品关联表模型

type OrderProduct struct {
	gorm.Model

	OrderID     uint
	ProductID   uint
	BuyQuantity int //购买数量
	BuyPice     int //购买价格
}
