package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//返回中间件的函数
func Require(pri string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//全部权限
		prisT, _ := c.Get("pris")
		pris := prisT.([]string)
		//判断其中是否有pri，需要的权限即可
		has := false
		for _, p := range pris {
			if p == pri {
				has = true
				break
			}
		}
		if !has {
			//无权限
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"error": "无相关权限",
			})
			return
		}
	}
}
