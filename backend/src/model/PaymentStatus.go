package model

import (
	"github.com/jinzhu/gorm"
)

//支付状态表模型

type PaymentStatus struct {
	gorm.Model

	Title string
}
