package controller

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context){
	//响应成功返回状态码和Json格式数据
		c.JSON(200,gin.H{
			"message":"pong",
		})
	}