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
