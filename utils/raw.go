package utils

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// resolveFilePath 校验 token 与路径参数，返回服务端上的实际文件路径。
// 校验失败时已写好响应，调用方只需返回。
func resolveFilePath(c *gin.Context, basePath string, username string, password string, secret string) (string, bool) {
	token := c.DefaultQuery("token", "")

	if !TokenCheck(username, password, token, secret) {
		c.JSON(401, gin.H{
			"ok":  false,
			"msg": "Not authorized",
		})
		return "", false
	}

	rawPath := c.DefaultQuery("path", "")
	if rawPath == "" {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Path parameter is required",
		})
		return "", false
	}
	decodedPath, err := url.QueryUnescape(rawPath)
	decodedPath = basePath + "/" + decodedPath
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": fmt.Sprintf("Invalid path: %v", err),
		})
		return "", false
	}
	file, err := os.Open(decodedPath)
	if err != nil {
		c.JSON(404, gin.H{
			"ok":  false,
			"msg": "File not found",
		})
		return "", false
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(404, gin.H{
			"ok":  false,
			"msg": "File not found",
		})
		return "", false
	}
	if fileInfo.IsDir() {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "The path is a directory",
		})
		return "", false
	}
	return decodedPath, true
}

func GetRaw(c *gin.Context, basePath string, username string, password string, secret string) {
	decodedPath, ok := resolveFilePath(c, basePath, username, password, secret)
	if !ok {
		return
	}
	c.File(decodedPath)
}
