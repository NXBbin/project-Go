package controller

//webAPP用户校验

import (
	"bytes"
	"config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"model"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//新增收货地址
func MemberAddressAdd(c *gin.Context) {
	user := member(c)
	if user == nil {
		//未验证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在该会员",
		})
		return
	}

	address := model.Address{}
	bindErr := c.ShouldBind(&address)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": bindErr.Error(),
		})
		return
	}
	address.UserID = user.ID
	//如果当前地址为默认，应将之前的默认地址去除
	if address.IsDefault {
		orm.Model(&model.Address{}).Where("user_id=? AND is_default=1", user.ID).Update("is_default", 0)
	}
	orm.Create(&address)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  address,
	})
}

//收货地址列表
func MemberAddressList(c *gin.Context) {
	user := member(c)
	if user == nil {
		//未验证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在该会员",
		})
		return
	}
	//查询
	as := []model.Address{}
	orm.Where("user_id=?", user.ID).
		//将默认地址倒叙排列
		Order("is_default desc").
		Find(&as)

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  as,
	})
}

//验证是否会员
func member(c *gin.Context) *model.User {
	//获取前端传递的token
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		return nil
	}

	//存在请求头，取出Token部分
	token := string(bytes.Replace([]byte(authorization), []byte("Bearer "), []byte(""), -1))
	if token == "" {
		return nil
	}

	// 校验token是否被篡改
	tokenObj, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.App["SECRET"]), nil
	})
	//token语法失败
	if parseErr != nil {
		return nil
	}

	//判断校验结果
	//解析验证失敗
	if !tokenObj.Valid {
		return nil
	}

	//验证通过
	userName := ""
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		userName = claims["aud"].(string)
	}

	//验证通过，获取用户名，利用用户名获取用户全部信息
	user := model.User{}
	orm.Where("user=?", userName).Find(&user)
	return &user
}

//将前端购物车信息同步到后端
func MemberCartSync(c *gin.Context) {
	user := member(c)
	if user == nil {
		//未验证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在该会员",
		})
		return
	}

	//登录会员成功，执行同步
	// 将浏览器端的购物车信息和会员已有的购物车数据合并起来。
	//获取浏览器端的购物车信息
	tempCart := c.PostForm("cart")
	//将前段的JSON格式数据转为结构体格式
	tempCartProducts := []struct {
		ProductID   uint
		BuyQuantity int
	}{}
	// 将json数据解码写入一个结构体
	json.Unmarshal([]byte(tempCart), &tempCartProducts)

	//获取服务器端已有的会员购物车信息
	memberCart := model.Cart{}
	orm.Where("user_id=?", user.ID).Find(&memberCart)
	//将取服务器端的JSON格式数据转为结构体格式
	memberCartProducts := []struct {
		ProductID   uint
		BuyQuantity int
	}{}
	// 将json数据解码写入一个结构体
	json.Unmarshal([]byte(memberCart.Content), &tempCartProducts)

	//同步，遍历浏览器端携带购物车信息和服务器端已有的购物车信息是否重复
	for _, tp := range tempCartProducts {
		//假设不重复
		exists := false
		for _, mp := range memberCartProducts {
			// 判断前端和后端购物车产品是否重复
			if tp.ProductID == mp.ProductID {
				exists = true
				break
			}
		}
		//若与服务器数据不重复，则增加商品
		if !exists {
			memberCartProducts = append(memberCartProducts, tp)
		}
	}
	//同步完成，存储到cart表中
	//JSON编码
	content, _ := json.Marshal(memberCartProducts)
	if memberCart.ID == 0 {
		//添加购物车数据
		memberCart.UserID = user.ID
		memberCart.Content = string(content)
		orm.Create(&memberCart)
	} else {
		//更新
		memberCart.Content = string(content)
		orm.Save(&memberCart)
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  string(content),
	})
}

//将后端数据同步前端
func MemberCartSet(c *gin.Context) {
	user := member(c)
	if user == nil {
		//未验证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在该会员",
		})
		return
	}

	//登录会员成功，执行同步
	// 将浏览器端的购物车信息和会员已有的购物车数据合并起来。
	//获取浏览器端的购物车信息
	tempCart := c.PostForm("cart")
	//将前段的JSON格式数据转为结构体格式
	tempCartProducts := []struct {
		ProductID   uint `json:"productID"`
		BuyQuantity int  `json:"buyQuantity"`
	}{}
	// 将json数据解码写入一个结构体
	json.Unmarshal([]byte(tempCart), &tempCartProducts)

	//同步完成，存储到cart表中
	//JSON编码
	content, _ := json.Marshal(tempCartProducts)

	//更新
	memberCart := model.Cart{}
	orm.Where("user_id=?", user.ID).Find(&memberCart)
	memberCart.Content = string(content)
	orm.Save(&memberCart)

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  string(content),
	})
}

//认证会员购物车信息
func MemberCart(c *gin.Context) {
	user := member(c)
	if user == nil {
		//未验证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "不存在该会员",
		})
		return
	}
	//验证通过,查找会员的购物车信息
	cart := model.Cart{}
	//判断是否存在购物车信息
	orm.Where("user_id=?", user.ID).Find(&cart)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  cart.Content,
	})

}

//认证会员
func MemberAuth(c *gin.Context) {
	//获取前端传递的token
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		//没有请求的响应
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无authorization请求头",
		})
		return
	}

	//存在请求头，取出Token部分
	token := string(bytes.Replace([]byte(authorization), []byte("Bearer "), []byte(""), -1))
	if token == "" {
		//没有请求的响应
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无Token请求",
		})
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
		return
	}

	//判断校验结果
	//解析验证失敗
	if !tokenObj.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token解析验证失敗",
		})
		return
	}

	//token 校验通过
	//获取token中包含的用户信息
	userName := ""
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		userName = claims["aud"].(string)
	}

	//验证通过，获取用户名，利用用户名获取用户全部信息
	user := model.User{}
	orm.Where("user=?", userName).Find(&user)
	//屏蔽隐私数据
	user.PasswordSalt = ""
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  user,
	})

	// 继续执行
	// c.Next()
}

//用户校验
func MemberLogin(c *gin.Context) {
	user := model.User{}
	//通过用户名获取用户信息，在比较密码是否正确
	postUser := c.PostForm("User")
	//判断该用户是否不存在，select * table from user where user=?;
	// v := orm.Where("user=?", postUser).Find(&user).RecordNotFound()
	// fmt.Println(postUser)
	if orm.Where("user=?", postUser).Find(&user).RecordNotFound() {
		//没有该用户
		c.JSON(http.StatusOK, gin.H{
			"error": "用户或密码错误(用户名)",
		})
		return
	}

	//用户存在，检查密码
	postPassword := c.PostForm("Password")
	//通过数据库中存储的Salt，处理密码
	pwdFunc := hmac.New(sha256.New, []byte(user.PasswordSalt))
	//将前端传递的密码进行摘要算法加密
	pwdFunc.Write([]byte(postPassword))
	//判断加密后的前端密码和数据库中的密码是否相同
	if user.Password != fmt.Sprintf("%x", pwdFunc.Sum(nil)) {
		//密码错误
		c.JSON(http.StatusOK, gin.H{
			"error": "用户或密码错误(密码)",
		})
		return
	}

	//生成Token签名key随机串
	mySigningKey := []byte(config.App["SECRET"])
	// Token内容
	claims := &jwt.StandardClaims{
		// 有效期，当前时间戳+有效时间(毫秒),有效一个月时间
		ExpiresAt: time.Now().Unix() + (30 * 24 * 3600),
		//发行人
		Issuer: "Backend",
		//对应的用户名
		Audience: user.User,
	}
	// 获取token构建器
	tokenBuilder := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//签名加密生成token
	token, err := tokenBuilder.SignedString(mySigningKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Token生成失败",
		})
		return
	}
	// fmt.Printf("%v %v", token, err)

	//用户名密码正确，认证通过
	user.Password, user.PasswordSalt = "", ""
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"user":  user,
		"token": token,
	})

}
