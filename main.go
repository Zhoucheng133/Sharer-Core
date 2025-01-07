package main

import (
	"strings"
	"virtual-directory-core/utils"

	"github.com/gin-gonic/gin"
)

func requestMiddleware(useAuth bool, username string, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/login") || !useAuth {
			c.Next()
		} else {
			if utils.TokenCheck(username, password, c.GetHeader("token")) {
				c.Next()
			} else {
				c.JSON(401, gin.H{
					"ok":  false,
					"msg": "Not authorized",
				})
				c.Abort()
			}
		}
	}
}

func main() {

	// -->测试代码<---
	username := "admin"
	password := "123456"
	useAuth := true
	//-->测试结束<--

	r := gin.New()
	r.Use(requestMiddleware(useAuth, username, password))
	r.POST("/api/list", utils.GetList)
	r.GET("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/raw"):
			utils.GetRaw(c)
		case strings.HasPrefix(c.Request.URL.Path, "/api/download"):
			utils.Download(c)
		case strings.HasPrefix(c.Request.URL.Path, "/api/login"):
			utils.Login(c)
		}
	})

	r.Run(":8080")
}
