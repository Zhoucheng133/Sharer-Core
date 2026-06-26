package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type CopyMoveBody struct {
	From  string   `json:"from"`
	Files []string `json:"files"`
	To    string   `json:"to"`
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0777); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func Copy(c *gin.Context, basePath string) {
	var reqBody CopyMoveBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		if reqBody.From == "" || len(reqBody.Files) == 0 || reqBody.To == "" {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Parameter illegal",
			})
		}
		return
	}

	srcBase := filepath.Join(basePath, reqBody.From)
	dstBase := filepath.Join(basePath, reqBody.To)

	if _, err := os.Stat(srcBase); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"ok":  false,
			"msg": "Source path not found",
		})
		return
	}

	for _, file := range reqBody.Files {
		srcPath := filepath.Join(srcBase, file)
		dstPath := filepath.Join(dstBase, file)

		info, err := os.Stat(srcPath)
		if err != nil {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "File or directory not found: " + file,
			})
			return
		}

		if info.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				c.JSON(400, gin.H{
					"ok":  false,
					"msg": "Failed to copy directory: " + file,
				})
				return
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				c.JSON(400, gin.H{
					"ok":  false,
					"msg": "Failed to copy file: " + file,
				})
				return
			}
		}
	}

	c.JSON(200, gin.H{
		"ok":  true,
		"msg": "",
	})
}
