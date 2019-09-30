package model

import (
	"github.com/jinzhu/gorm"
)

//产品图像表模型

type Image struct {
	gorm.Model

	ProductID  uint
	IsDefault  bool //默认展示图
	SortOrder  int
	Host       string
	Image      string //客户端上传后，保存到服务器的字段
	ImageSmall string //图像缩略图路径
	ImageBig   string //图像缩略图路径
}
