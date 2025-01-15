package utils

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type DelBody struct {
	Path  string   `json:"path"`
	Files []string `json:"files"`
}

func delHandler(c *gin.Context, data DelBody) {
	if _, err := os.Stat(data.Path); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"ok":  false,
			"msg": "Path not found",
		})
		return
	}
	for _, file := range data.Files {
		fullPath := filepath.Join(data.Path, file)

		// 删除文件或文件夹
		if err := os.RemoveAll(fullPath); err != nil {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "failed to delete file or directory",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"ok":  true,
		"msg": "",
	})
}

func DelRequest(c *gin.Context, basePath string) {
	var reqBody DelBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		if reqBody.Path == "" || len(reqBody.Files) == 0 {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Parameter illegal",
			})
		}
		return
	}
	reqBody.Path = basePath + reqBody.Path
	delHandler(c, reqBody)
}
