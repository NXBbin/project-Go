package controller

import (
	"config"
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/gomango/imgtype"
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

	//file.Header // 检查从前端上传的文件是否可靠

	//打开临时文件
	f, _ := file.Open()
	defer f.Close()
	content := make([]byte, file.Size)
	//读取文件全部内容
	f.Read(content)

	//通过文件内容检查文件类型(依赖imgtype第三方包)，得到文件类型
	mime, err := imgtype.DetectBytes(content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	//用md5计算文件内容得到，散列值给文件命名
	// dst := fmt.Sprintf("%x", md5.Sum(content)) + filepath.Ext(file.Filename)
	dst := fmt.Sprintf("%x", md5.Sum(content))

	//利用文件名首两位字符，构建子目录
	subPath := string(dst[0]) + string(os.PathSeparator) + string(dst[1]) + string(os.PathSeparator)
	savePath := config.App["UPLOAD_PATH"] + subPath
	os.MkdirAll(savePath, 0755)
	// fmt.Println("=====", subPath)
	// fmt.Println("-----", savePath)

	//将文件重新编码，保证安全性
	switch mime {
	case "image/jpeg":
		// 打开文件
		srcFile, _ := file.Open()
		defer srcFile.Close()
		img, err := jpeg.Decode(srcFile) //解码
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		//确定文件类型规定后缀名，后编码创建文件，并保存到指定路径下
		dst += ".jpg"
		imgFile, err := os.Create(savePath + dst)
		// log.Println(savePath + dst)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer imgFile.Close()
		//将解码后的img文件，重新编码写入该路径下的文件中。
		err = jpeg.Encode(imgFile, img, &jpeg.Options{36}) //压缩值，值越大文件越大
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}

	case "image/png":
		srcFile, _ := file.Open()
		defer srcFile.Close()
		img, err := png.Decode(srcFile) //解码
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		//确定文件类型规定后缀名，后编码创建文件，并保存到指定路径下
		dst += ".png"
		imgFile, err := os.Create(savePath + dst)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer imgFile.Close()
		err = png.Encode(imgFile, img)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}

	case "image/gif":
		srcFile, _ := file.Open()
		defer srcFile.Close()
		img, err := gif.Decode(srcFile) //解码
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		//确定文件类型规定后缀名，后编码创建文件，并保存到指定路径下
		dst += ".gif"
		imgFile, err := os.Create(savePath + dst)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer imgFile.Close()
		err = gif.Encode(imgFile, img, &gif.Options{NumColors: 256})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	//制作缩略图
	//确定缩放宽高值
	// thumbW,thumbH = config.App["THUMB_SMALL_W"],config.App["THUMB_SMALL_H"]
	small, err := MakeThumb(savePath+dst, 146, 146)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	//log.Println(small)

	//制作中图
	big, err := MakeThumb(savePath+dst, 1460, 1460)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	//log.Println(big)

	//保证目录存在，创建目录
	// if err := os.MkdirAll(savePath, 0755); err != nil {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// 将文件存储到指定的路径下
	// if err := c.SaveUploadedFile(file, savePath+dst); err != nil {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	//存储到数据库

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": map[string]string{
			"Image":      dst,
			"ImageSmall": small,
			"ImageBig":   big,
		},
	})
	return
}

//制作缩略图（原图像，缩放宽度，高度）
func MakeThumb(srcF string, thumbW, thumbH int) (string, error) {
	//打开原图
	srcFile, err := os.Open(srcF)
	if err != nil {
		log.Println("打开原图失败")
		return "", err
	}

	//获取尺寸
	srcConfig, _, err := image.DecodeConfig(srcFile)
	srcFile.Close()
	//原图的宽，高
	srcW, srcH := srcConfig.Width, srcConfig.Height
	//if err != nil {
	//	return "", err
	//}
	// log.Println(srcConfig.Width, srcConfig.Height)

	//计算宽之比，高之比（计算：缩略图有效区域的尺寸）
	var dstW, dstH int
	// 当缩略比例，宽值>高值时，说明是宽图
	if float64(srcW)/float64(thumbW) >= float64(srcH)/float64(thumbH) {
		//缩放的宽度应与缩略值一致
		dstW = thumbW
		//等比例计算出缩放高度,Round四舍五入
		dstH = int(math.Round(float64(dstW) * (float64(srcH) / float64(srcW))))
	} else {
		//宽值<高值时，说明是高图
		dstH = thumbH //缩放的高度应与缩略值一致
		//等比例计算出缩放高度,Round四舍五入
		dstW = int(math.Round(float64(dstH) * (float64(srcW) / float64(srcH))))
	}

	// 计算缩略图居中位置
	var dstX, dstY int
	dstX = (thumbW - dstW) / 2
	dstY = (thumbH - dstH) / 2
	//log.Println(dstW, dstH, dstX, dstY)

	//创建缩略图
	//重新采样，1.从原图上采集新的点，2.重新编码为缩略图图片
	//利用成品包，将缩略图做出来，在放在白色背景上

	//图像的范围
	thumbRect := image.Rect(0, 0, thumbW, thumbH)
	// 具有指定范围的像素色彩信息
	thumb := image.NewRGBA(thumbRect)

	//使用白色进行背景填充
	bgColor := color.RGBA{
		255, 255, 255, 255,
	}

	//填充在白色画布上
	draw.Draw(thumb, thumbRect, &image.Uniform{C: bgColor}, image.Pt(0, 0), draw.Src)

	//缩略图，将原图src图，画到thumb图上
	srcFile1, err := os.Open(srcF)
	if err != nil {
		return "", err
	}
	defer srcFile1.Close()
	src, err := jpeg.Decode(srcFile1) //打开原图
	//使用imaging包中的采样算法
	// thumb := imaging.Fit(src, thumbW, thumbH, imaging.Lanczos)//自适应
	thumSrc := imaging.Resize(src, dstW, dstH, imaging.Lanczos) //直接生成

	//将缩略图居中于白色画布中
	srcRect := image.Rect(dstX, dstY, dstX+dstW, dstY+dstH)
	draw.Draw(thumb, srcRect, thumSrc, image.Pt(0, 0), draw.Src)

	//存储缩略图
	tempFile := os.TempDir() + "/thumb.png" //临时目录
	thumbFile, _ := os.Create(tempFile)     //临时文件
	png.Encode(thumbFile, thumb)
	thumbFile.Close()

	//MD5摘要算法确定文件名
	content, _ := ioutil.ReadFile(tempFile)
	dst := fmt.Sprintf("%x", md5.Sum(content)) + ".png"
	//构建子目录
	subPath := string(dst[0]) + string(os.PathSeparator) + string(dst[1]) + string(os.PathSeparator)
	savePath := config.App["UPLOAD_PATH"] + subPath
	os.MkdirAll(savePath, 0755)

	dstFile, err := os.Create(savePath + dst)
	dstFile.Write(content)
	dstFile.Close()

	return dst, nil
}
