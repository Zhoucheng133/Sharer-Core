package utils

import (
	"embed"
	"fmt"
	"mime"
	"net/http"
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
	path := c.Request.URL.Path
	if path == "" || path == "/" || strings.HasPrefix(path, "/login") {
		path = filepath.Join(basePath, "index.html")
	} else {
		path = filepath.Join(basePath, path)
	}

	// 打开文件
	file, err := http.Dir(basePath).Open(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("File not found: %s", path)})
		return
	}
	defer file.Close()

	// 解析文件类型
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream" // 默认类型
	}

	// 设置响应头
	c.Header("Content-Type", contentType)

	// 读取文件数据并返回
	c.File(path)
}
