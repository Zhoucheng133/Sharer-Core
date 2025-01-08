package main

import (
	"sharer-core/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func requestMiddleware(useAuth bool, username string, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") && useAuth {
			if strings.HasPrefix(c.Request.URL.Path, "/api/auth") {
				c.Next()
			} else if utils.TokenCheck(username, password, c.GetHeader("token")) {
				c.Next()
			} else {
				c.JSON(401, gin.H{
					"ok":  false,
					"msg": "Not authorized",
				})
				c.Abort()
			}
		} else {
			c.Next()
		}
	}
}

func main() {

	// -->测试代码<---
	username := "admin"
	password := "123456"
	useAuth := false
	basePath := "/Users/zhoucheng/Desktop"
	//-->测试结束<--

	r := gin.New()
	r.Use(requestMiddleware(useAuth, username, password))
	r.POST("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/list"):
			utils.GetList(c, basePath)
		}
	})
	r.GET("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/raw"):
			utils.GetRaw(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/download"):
			utils.Download(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/login"):
			utils.Login(c)
		case strings.HasPrefix(c.Request.URL.Path, "/api/auth"):
			utils.Auth(useAuth, c)
		}
	})

	r.Run(":8080")
}
