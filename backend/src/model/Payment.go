package model

import (
	"github.com/jinzhu/gorm"
)

//支付表模型

type Payment struct {
	gorm.Model

	Title  string
	Key    string //对应Title
	Intro  string //介绍
	Status int    //0 1 2 分别表示禁用，启用，维护等状态

}
