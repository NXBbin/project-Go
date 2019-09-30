package controller

import (
	"config"
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//图像上传
func ImageUpload(c *gin.Context) {
	//获取前端上传的文件
	file, fileErr := c.FormFile("file")
	if fileErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": fileErr.Error(),
		})
		return
	}

	//打开临时文件
	f, _ := file.Open()
	defer f.Close()
	content := make([]byte, file.Size)
	//读取文件全部内容
	f.Read(content)

	//用md5计算文件内容得到，散列值给文件命名,并连接上文件后缀
	dst := fmt.Sprintf("%x", md5.Sum(content)) + filepath.Ext(file.Filename)

	//利用文件名首两位字符，构建子目录
	subPath := string(dst[0]) + string(os.PathSeparator) + string(dst[1]) + string(os.PathSeparator)
	savePath := config.App["UPLOAD_PATH"] + subPath
	// fmt.Println("=====", subPath)
	// fmt.Println("-----", savePath)

	//保证目录存在，创建目录
	if err := os.MkdirAll(savePath, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 将文件存储到指定的路径下
	if err := c.SaveUploadedFile(file, savePath+dst); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	//存储到数据库

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": map[string]string{
			"filename": dst,
			// "subpath":string(dst[0]) + "/" string(dst[i]),
		},
	})
	return

}
