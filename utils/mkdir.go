package utils

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FolderBody struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func createHandler(path string, name string, c *gin.Context) {
	dirPath := filepath.Join(path, name)
	err := os.Mkdir(dirPath, 0777)
	if err != nil {
		if os.IsExist(err) {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Directory exists",
			})
			return
		} else {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Create failed",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"ok":  true,
		"msg": "",
	})
}

func CreateFolder(c *gin.Context, basePath string) {
	var reqBody FolderBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		if reqBody.Path == "" || reqBody.Name == "" {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Parameter illegal",
			})
		}
		return
	}
	reqBody.Path = basePath + reqBody.Path
	createHandler(reqBody.Path, reqBody.Name, c)
}
