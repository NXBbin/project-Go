package controller

import (
	"crypto/md5"
	"fmt"
	"image/png"
	"strings"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
)

//登录验证码

func CheckCode(c *gin.Context) {
	// #设置 captcha 对象,设置字体样式
	cap := captcha.New()
	font := "../src/config/Lobster-Regular.ttf"
	cap.SetFont(font)

	//创建验证码为4个字符长度，返回图片和图片中对应的值
	img, code := cap.Create(4, captcha.NUM)
	// 将图像作为响应主体
	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, img)

	// 将验证码码值存储在redis（为了验证），格式：IP + 浏览器信息
	// [127.0.0.1]:64756 + [Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:69.0) Gecko/20100101
	//截取IP，去除端口部分，因每次请求端口会不一样
	i := strings.Index(c.Request.RemoteAddr, "]")
	//以md5摘要的形式存储
	key := fmt.Sprintf("%x", md5.Sum([]byte(c.Request.RemoteAddr[1:i]+c.Request.Header["User-Agent"][0])))
	//写入redis，并且设置有效期2分钟
	Rds.Do("SET", "code_"+key, code)
	Rds.Do("expire", "code_"+key, 2*60)

}
