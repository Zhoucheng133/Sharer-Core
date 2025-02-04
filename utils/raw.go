package utils

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func GetRaw(c *gin.Context, basePath string, username string, password string) {

	token := c.DefaultQuery("token", "")

	if !TokenCheck(username, password, token) {
		c.JSON(401, gin.H{
			"ok":  false,
			"msg": "Not authorized",
		})
		return
	}

	rawPath := c.DefaultQuery("path", "")
	if rawPath == "" {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Path parameter is required",
		})
		return
	}
	decodedPath, err := url.QueryUnescape(rawPath)
	decodedPath = basePath + "/" + decodedPath
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": fmt.Sprintf("Invalid path: %v", err),
		})
		return
	}
	file, err := os.Open(decodedPath)
	if err != nil {
		c.JSON(404, gin.H{
			"ok":  false,
			"msg": "File not found",
		})
		return
	}
	defer file.Close()
	fileInfo, _ := os.Stat(decodedPath)
	if fileInfo.IsDir() {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "The path is a directory",
		})
		return
	}
	c.File(decodedPath)
}
