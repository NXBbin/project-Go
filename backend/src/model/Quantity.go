package model

import (
	"github.com/jinzhu/gorm"
)

//库存表模型

type Quantity struct {
	gorm.Model

	ProductID uint
	Number    int  //库存
	StoreID   uint //哪个仓库
}
