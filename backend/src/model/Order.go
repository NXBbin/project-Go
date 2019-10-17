package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

//订单表模型

type Order struct {
	gorm.Model

	Sn               string //订单号
	AddressID        uint
	UserID           uint
	OrderStatusID    uint
	PaymentID        uint
	PaymentStatusID  uint
	PaymentSn        string
	ShippingID       uint
	ShippingStatusID uint
	ShippingSn       string
	ShippingAmount   int //配送费
	Amount           int
	TaxID            uint //税类型
	TaxAmount        int  //税费
	ProductAmount    int
	OrderTime        time.Time
	Note             string //备注
}
