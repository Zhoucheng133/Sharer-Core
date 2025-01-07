package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Download(c *gin.Context) {
	rawPath := c.DefaultQuery("path", "")
	if rawPath == "" {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Path parameter is required",
		})
		return
	}
	decodedPath, err := url.QueryUnescape(rawPath)
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(decodedPath)))
	c.File(decodedPath)
}
