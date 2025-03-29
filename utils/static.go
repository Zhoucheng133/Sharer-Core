package utils

import (
	"embed"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticHandler(c *gin.Context, f embed.FS) {
	path := c.Request.URL.Path
	if path == "" || path == "/" || strings.HasPrefix(path, "/login") {
		path = "Sharer-Web/dist/index.html"
	} else {
		path = "Sharer-Web/dist" + path
	}
	file, err := f.Open(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("File not found: %s", path)})
		return
	}
	defer file.Close()
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	c.Header("Content-Type", contentType)
	c.DataFromReader(http.StatusOK, -1, path, file, map[string]string{})
}

func DynamicLibStaticHandler(c *gin.Context, basePath string) {
	path := strings.TrimPrefix(c.Request.URL.Path, "/")
	if path == "" || path == "/" || strings.HasPrefix(path, "login") {
		path = "index.html"
	}

	fullPath := filepath.Join(basePath, path)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found: " + fullPath})
		return
	}

	// 直接返回文件
	c.File(fullPath)
}
