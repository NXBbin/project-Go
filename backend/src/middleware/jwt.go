package middleware

import (
	"bytes"
	"config"

	// "fmt"
	// "model"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//JWT 校验Token中间件
func JWTToken(c *gin.Context) {
	//获取前端传递的token
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		//没有请求的响应
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无authorization请求头",
		})
		c.Abort()
		return
	}

	//存在请求头，取出Token部分
	token := string(bytes.Replace([]byte(authorization), []byte("Bearer "), []byte(""), -1))
	if token == "" {
		//没有请求的响应
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无Token请求",
		})
		c.Abort()
		return
	}

	// 校验token是否被篡改
	tokenObj, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.App["SECRET"]), nil
	})
	//token语法失败
	if parseErr != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": parseErr.Error(),
		})
		c.Abort()
		return
	}

	//判断校验结果
	if !tokenObj.Valid {
		c.JSON(http.StatusOK, gin.H{
			"error": parseErr.Error(),
		})
		c.Abort()
		return
	}

	//token 校验通过
	//获取token中包含的用户信息
	userName := ""
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		userName = claims["aud"].(string)
	}
	//先在中间件存储起来
	c.Set("userName", userName)

	// 继续执行
	// c.Next()
}
