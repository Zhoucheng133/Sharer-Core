package utils

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Move(c *gin.Context, basePath string) {
	var reqBody CopyMoveBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		if reqBody.From == "" || len(reqBody.Files) == 0 || reqBody.To == "" {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Parameter illegal",
			})
		}
		return
	}

	srcBase := filepath.Join(basePath, reqBody.From)
	dstBase := filepath.Join(basePath, reqBody.To)

	if _, err := os.Stat(srcBase); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"ok":  false,
			"msg": "Source path not found",
		})
		return
	}

	for _, file := range reqBody.Files {
		srcPath := filepath.Join(srcBase, file)
		dstPath := filepath.Join(dstBase, file)

		if err := os.Rename(srcPath, dstPath); err != nil {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Failed to move file or directory: " + file,
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"ok":  true,
		"msg": "",
	})
}
