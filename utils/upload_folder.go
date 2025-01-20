package utils

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFolder(c *gin.Context, basePath string) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Unable to parse multipart form",
		})
		return
	}
	files := form.File["files"]
	relativePaths := form.Value["paths"]

	if len(files) != len(relativePaths) {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Mismatch between files and paths",
		})
		return
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		relativePath := relativePaths[i]
		// fullPath := basePath + relativePath
		fullPath := filepath.Join(basePath, relativePath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			c.JSON(500, gin.H{
				"ok":  false,
				"msg": "Unable to create directory: " + dir,
			})
			return
		}
		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			c.JSON(500, gin.H{
				"ok":  false,
				"msg": "Failed to save file: " + file.Filename,
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"ok":  true,
		"msg": "",
	})
}
