package main

import (
	"embed"
	"sharer-core/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func requestMiddleware(useAuth bool, username string, password string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// -->允许跨域内容，开始<--
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// -->允许跨域内容，结束<--

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

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

//go:embed Sharer-Web/dist/*
var staticFiles embed.FS

func main() {

	// -->测试代码<---
	username := "admin"
	password := "123456"
	useAuth := false
	basePath := "/Users/zhoucheng/Desktop"
	//-->测试结束<--

	// 所有路径请求path需要添加头/
	r := gin.New()
	r.Use(requestMiddleware(useAuth, username, password))
	r.POST("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/list"):
			utils.GetList(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/multidownload"):
			utils.MultiDownload(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/upload"):
			utils.Upload(c, basePath)
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
		default:
			utils.StaticHandler(c, staticFiles)
		}
	})

	r.Run(":8080")
}
