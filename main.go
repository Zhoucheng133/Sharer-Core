package main

import (
	"strings"
	"virtual-directory-core/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.POST("/api/list", utils.GetList)
	r.GET("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/raw"):
			utils.GetRaw(c)
		}
	})

	r.Run(":8080")
}
