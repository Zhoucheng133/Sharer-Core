package utils

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context, basePath string) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Unable to parse multipart form",
		})
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "No files uploaded",
		})
		return
	}
	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Unable to create upload directory",
		})
		return
	}
	for _, file := range files {
		// 构建文件存储路径
		dst := filepath.Join(basePath, file.Filename)

		// 保存文件到指定路径
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			c.JSON(400, gin.H{
				"ok":    false,
				"error": "Unable to save file(s)",
			})
			return
		}

		// 返回上传文件信息
		c.JSON(200, gin.H{
			"ok":  true,
			"msg": "",
		})
	}
}
