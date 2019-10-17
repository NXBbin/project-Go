package model

import (
	"github.com/jinzhu/gorm"
)

//配送表模型

type Shipping struct {
	gorm.Model

	Title  string
	Key    string
	Intro  string //介绍
	Status int    //0 1 2 分别表示禁用，启用，维护等状态

}
