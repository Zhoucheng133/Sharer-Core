package utils

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type ListBody struct {
	Path string `json:"path"`
}

type Item struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

type RequestResponse struct {
	Ok    bool   `json:"ok"`
	Items []Item `json:"items"`
	Msg   string `json:"msg"`
}

func listHandler(path string) RequestResponse {
	files, err := os.ReadDir(path)
	if err != nil {
		return RequestResponse{
			Ok:    false,
			Items: nil,
			Msg:   "Path does not exist or is not accessible",
		}
	}
	var items []Item
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}
		var size int64
		if file.IsDir() {
			size = 0
		} else {
			size = fileInfo.Size()
		}
		items = append(items, Item{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Size:  size,
		})
	}

	return RequestResponse{
		Ok:    true,
		Items: items,
		Msg:   "",
	}
}

func GetList(c *gin.Context, basePath string) {
	var body ListBody

	if err := c.ShouldBind(&body); err == nil {
		if body.Path == "" {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Path parameter is required",
			})
		} else {
			c.JSON(200, listHandler(fmt.Sprint(basePath, body.Path)))
		}
		return
	}
	c.JSON(400, gin.H{
		"ok":  false,
		"msg": "Bad request",
	})

}
