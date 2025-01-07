package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	fileInfo, _ := os.Stat(decodedPath)
	if fileInfo.IsDir() {
		// c.JSON(400, gin.H{
		// 	"ok":  false,
		// 	"msg": "The path is a directory",
		// })
		// return
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(decodedPath)))
		c.Header("Content-Type", "application/zip")
		zipWriter := zip.NewWriter(c.Writer)
		defer zipWriter.Close()
		err := filepath.Walk(decodedPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 跳过根目录自身
			if path == decodedPath {
				return nil
			}

			// 为每个文件创建一个新的 ZIP 条目
			zipEntry, err := zipWriter.Create(strings.TrimPrefix(path, decodedPath+"/"))
			if err != nil {
				return err
			}

			// 打开文件并将其内容写入 ZIP 条目
			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				// 将文件数据写入到 zipEntry
				_, err = io.Copy(zipEntry, file)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			c.JSON(500, gin.H{
				"ok":  false,
				"msg": fmt.Sprintf("Failed to add files to zip: %v", err),
			})
			return
		}
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(decodedPath)))
	c.File(decodedPath)
}
