package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type MultiDownloadType struct {
	Path  string   `json:"path"`
	Files []string `json:"files"`
}

func Download(c *gin.Context, basePath string, username string, password string) {

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
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s.zip", url.QueryEscape(filepath.Base(decodedPath))))
		c.Header("Content-Type", "application/zip")
		zipWriter := zip.NewWriter(c.Writer)
		defer zipWriter.Close()
		err := addDirToZip(decodedPath, url.QueryEscape(filepath.Base(decodedPath)), zipWriter)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(500, gin.H{
				"ok":  false,
				"msg": "Failed to add directory",
			})
			return
		}
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(filepath.Base(decodedPath))))
	c.File(decodedPath)
}

// MultiDownloadHandler 处理多文件下载的函数
func MultiDownloadHandler(c *gin.Context, basePath string, data MultiDownloadType) {

	path := filepath.Join(basePath, data.Path)

	for _, fileName := range data.Files {
		filePath := filepath.Join(path, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.Header("Content-Type", "application/json")
			c.JSON(404, gin.H{
				"ok":  false,
				"msg": fmt.Sprintf("File or directory %s does not exist", fileName),
			})
			return
		}
	}

	var fileName string = "多文件下载.zip"
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(fileName)))

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// 递归遍历目录，添加文件到 ZIP 包中
	for _, fileName := range data.Files {
		filePath := filepath.Join(path, fileName)

		// 检查文件或目录是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.Header("Content-Type", "application/json")
			c.JSON(404, gin.H{
				"ok":  false,
				"msg": fmt.Sprintf("File or directory %s does not exist", fileName),
			})
			return
		}

		info, err := os.Stat(filePath)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(500, gin.H{
				"ok":  false,
				"msg": fmt.Sprintf("Unable to stat file %s", fileName),
			})
			return
		}

		if !info.IsDir() {
			err := addFileToZip(filePath, fileName, zipWriter)
			if err != nil {
				c.Header("Content-Type", "application/json")
				c.JSON(500, gin.H{
					"ok":  false,
					"msg": fmt.Sprintf("Failed to add file %s", fileName),
				})
				return
			}
		} else {
			// 如果是目录，递归遍历该目录下的文件
			err := addDirToZip(filePath, fileName, zipWriter)
			if err != nil {
				c.Header("Content-Type", "application/json")
				c.JSON(500, gin.H{
					"ok":  false,
					"msg": fmt.Sprintf("Failed to add directory %s", fileName),
				})
				return
			}
		}
	}
}

// 将文件添加到 ZIP 文件中
func addFileToZip(filePath, fileName string, zipWriter *zip.Writer) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %v", fileName, err)
	}
	defer file.Close()

	// 创建一个新的文件头，表示要将文件加入到 zip 包中
	zipFile, err := zipWriter.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create zip entry for %s", fileName)
	}

	// 将文件内容复制到 zip 文件中
	_, err = io.Copy(zipFile, file)
	if err != nil {
		return fmt.Errorf("failed to write file %s to zip", fileName)
	}

	return nil
}

// 递归将目录下的所有文件添加到 ZIP 文件中
func addDirToZip(dirPath, dirName string, zipWriter *zip.Writer) error {
	// 遍历目录中的文件和子目录
	return filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录本身，添加文件
		if info.IsDir() {
			return nil
		}

		// 计算相对于根目录的相对路径
		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return err
		}

		// 在 ZIP 文件中创建相应的条目
		return addFileToZip(filePath, filepath.Join(dirName, relPath), zipWriter)
	})
}

func MultiDownload(c *gin.Context, basePath string, username string, password string) {

	token := c.DefaultQuery("token", "")

	if !TokenCheck(username, password, token) {
		c.JSON(401, gin.H{
			"ok":  false,
			"msg": "Not authorized",
		})
		return
	}

	data, err := url.QueryUnescape(c.Query("data"))
	if err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Invalid path parameter",
		})
		return
	}
	if len(data) == 0 {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Data parameter is required",
		})
		return
	}
	var filesData MultiDownloadType

	err = json.Unmarshal([]byte(data), &filesData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	MultiDownloadHandler(c, basePath, filesData)
}
