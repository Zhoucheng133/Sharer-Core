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
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Parameter illegal",
			})
		}
		return
	}
	oldPath := filepath.Join(basePath, reqBody.Path, reqBody.OldName)
	newPath := filepath.Join(basePath, reqBody.Path, reqBody.NewName)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		c.JSON(400, gin.H{
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
