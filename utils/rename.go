package utils

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type RenameBody struct {
	Path    string `json:"path"`
	NewName string `json:"newName"`
	OldName string `json:"oldName"`
}

func Rename(c *gin.Context, basePath string) {
	var reqBody RenameBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		if reqBody.Path == "" || reqBody.NewName == "" || reqBody.OldName == "" {
			c.JSON(200, gin.H{
				"ok":  false,
				"msg": "Parameter illegal",
			})
		}
		return
	}
	oldPath := filepath.Join(basePath, reqBody.Path, reqBody.OldName)
	newPath := filepath.Join(basePath, reqBody.Path, reqBody.NewName)

	// 目标名称已存在时直接阻止，避免 os.Rename 在部分平台上覆盖同名文件。
	if _, err := os.Lstat(newPath); err == nil {
		c.JSON(200, gin.H{
			"ok":  false,
			"msg": "A file or directory with the new name already exists",
		})
		return
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		c.JSON(200, gin.H{
			"ok":  false,
			"msg": "Rename failed",
		})
	} else {
		c.JSON(200, gin.H{
			"ok":  true,
			"msg": "",
		})
	}
}
