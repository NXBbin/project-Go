package controller

//User 表控制器（增删改查）代码，脚手架模板

import (
	"config"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"model"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//用户校验
func UserAuth(c *gin.Context) {
	//获取用户在表单中输入的验证码信息
	postCode := c.PostForm("Code")
	if postCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "请输入验证码",
		})
		return
	}
	//于redis中的数据对比
	i := strings.Index(c.Request.RemoteAddr, "]")
	key := fmt.Sprintf("%x", md5.Sum([]byte(c.Request.RemoteAddr[1:i]+c.Request.Header["User-Agent"][0])))
	result, _ := Rds.Do("get", "code_"+key)
	code := string(result.([]byte))
	//对比前后端的验证码信息是否一致
	if postCode != code {
		c.JSON(http.StatusOK, gin.H{
			"error": "验证码错误",
		})
		return
	}

	//验证码正确，立即删除redis中对应的数据，避免重用
	Rds.Do("DEL", "code_"+key)

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
		// 有效期，当前时间戳+有效时间(毫秒)
		ExpiresAt: time.Now().Unix() + 7200,
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

//列表
func UserList(c *gin.Context) {
	//搜索(筛选）
	condStr := ""
	condParams := []string{}
	//确定搜索条件

	//排序
	//获取排序参数（字段）
	orderStr := ""
	sortProp := c.DefaultQuery("sortProp", "")
	//排序方式：ascending ， descending
	sortOrder := c.DefaultQuery("sortOrder", "")
	//判断用户是否请求了排序
	if sortProp != "" && sortOrder != "" {
		//默认是升序，若传递的是降序请求，则设置为DESC
		sortMethod := "ASC"
		if "descending" == sortOrder {
			sortMethod = "DESC"
		}
		//拼凑：name ASC||DESC
		orderStr = sortProp + " " + sortMethod
		// log.Println("order语句：", orderStr)
	}

	//翻页: /products?currentPage =&pageSize=
	//获取请求的当前页码,默认第一页
	currentPageStr := c.DefaultQuery("currentPage", "1")
	//每页的显示的数量（偏移量）
	pageSizeStr := c.DefaultQuery("pageSize", "5")
	//将从前端获取到的页码数据转换类型（int转string）
	currentPage, pageErr := strconv.Atoi(currentPageStr)
	//若用户传递的参数不是整形数据（不合法数据），则指定页码为1
	if pageErr != nil {
		currentPage = 1
	}

	pageSize, sizeErr := strconv.Atoi(pageSizeStr)
	if sizeErr != nil {
		pageSize = 5
	}

	//获取总记录数
	total := 0
	orm.Model(&model.User{}).Where(condStr, condParams).Count(&total)
	//计算偏移量
	offset := (currentPage - 1) * pageSize

	//获取product模型
	ms := []model.User{}
	//获取展示数量和偏移量,输出数据获
	orm.Where(condStr, condParams).Order(orderStr).Limit(pageSize).Offset(offset).Find(&ms)
	//隐私数据，制空属性，前端不显示密码和密码摘要
	for i, _ := range ms {
		ms[i].PasswordSalt = ""
		ms[i].Password = ""
	}

	//响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  ms,
		//页数
		"pager": map[string]int{
			//当前页
			"currentPage": currentPage,
			//偏移量
			"pageSize": pageSize,
			//数据总量
			"total": total,
		},
	})
}

//删除
func UserDelete(c *gin.Context) {
	//获取前端参数（需要删除的ID）
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定ID",
		})
		return
	}
	//确定模型对象（表）
	m := model.User{}
	//将前端数据进行类型转换（string转uint）
	id, _ := strconv.Atoi(ID)
	m.ID = uint(id)
	//执行删除
	orm.Delete(&m)
	//判断删除是否有错误
	if orm.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
	}
	//无错误响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
	})
}

//添加
func UserCreate(c *gin.Context) {
	//确定模型对象（表）
	m := model.User{}
	//使用c.ShouldBind(),绑定并解析post数据
	err := c.ShouldBind(&m)
	if err != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	//特定数据设置
	//生成 密码salt
	saltChars := "zxcasdqweghjn@*&^%#1234567POLK." //salt随机串
	saltLen := 6                                   //salt长度
	salt := ""
	for i := 0; i < saltLen; i++ {
		// 随机数
		index := rand.Int31n(int32(len(saltChars)))
		// log.Println(index)
		salt += string(saltChars[index])
	}
	//将salt摘要记录
	m.PasswordSalt = salt
	log.Printf("Salt摘要：%s\n", salt)

	//为密码做信息摘要（获取加密器）
	pwdFunc := hmac.New(sha256.New, []byte(salt))
	//加密数据
	pwdFunc.Write([]byte(m.Password))                //向原密码写入摘要密串
	m.Password = fmt.Sprintf("%x", pwdFunc.Sum(nil)) //计算生成摘要值
	log.Printf("加密后结果：%s\n", m.Password)

	// 将关联临时关闭，(若不关闭，也会将添加的数据自动添加到关联的表中)，并添加数据
	orm.Set("gorm:save_associations", false).Create(&m)
	if orm.Error != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
		return
	}

	// 查询相关联的表数据

	//响应正确数据
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  m,
	})
}

//更新
func UserUpdate(c *gin.Context) {
	//获取前端传递的请求更新数据ID
	IDstr := c.DefaultQuery("ID", "")
	// 没有传递时，错误响应
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定更新ID",
		})
		return
	}
	//类型转换（string转int）
	ID, _ := strconv.Atoi(IDstr)

	// 获取需要更新的表数据
	m := model.User{}
	m.ID = uint(ID)
	orm.Find(&m)

	// 绑定并解析 数据
	err := c.ShouldBind(&m)
	if err != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 特定数据设置

	// 临时取消关联,并更新数据
	orm.Model(&m).
		Set("gorm:save_associations", false).
		//仅更新非空字段（前端编辑时无法更改密码和Salt）
		Update(m)

	if orm.Error != nil {
		//响应错误
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
		return
	}

	//响应正确数据
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data":  m,
	})
}
