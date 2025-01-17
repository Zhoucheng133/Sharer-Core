package main

import (
	"embed"
	"sharer-core/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func requestMiddleware(username string, password string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// -->允许跨域内容，开始<--
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Disposition, File-Name")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// -->允许跨域内容，结束<--

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		switch {
		case !strings.HasPrefix(c.Request.URL.Path, "/api"):
			c.Next()
		case len(username) == 0 && len(password) == 0:
			c.Next()
		case strings.HasPrefix(c.Request.URL.Path, "/api/auth"),
			strings.HasPrefix(c.Request.URL.Path, "/api/download"),
			strings.HasPrefix(c.Request.URL.Path, "/api/multidownload"),
			strings.HasPrefix(c.Request.URL.Path, "/api/raw"):
			// 这些操作将token写在URL param中
			c.Next()
		default:
			// 其余操作token写在header中
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

//go:embed Sharer-Web/dist/*
var staticFiles embed.FS

func main() {

	// -->测试代码<---
	username := "admin"
	password := "123456"
	// useAuth := true
	// 如果username=="" && password==""表明不需要验证
	basePath := "/Users/zhoucheng/Downloads"
	//-->测试结束<--

	// 所有路径请求path需要添加头/
	r := gin.New()
	r.Use(requestMiddleware(username, password))
	r.POST("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/list"):
			utils.GetList(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/upload"):
			utils.Upload(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/del"):
			utils.DelRequest(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/mkdir"):
			utils.CreateFolder(c, basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/rename"):
			utils.Rename(c, basePath)
		}
	})
	r.GET("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/raw"):
			utils.GetRaw(c, basePath, username, password)
		case strings.HasPrefix(c.Request.URL.Path, "/api/download"):
			utils.Download(c, basePath, username, password)
		case strings.HasPrefix(c.Request.URL.Path, "/api/multidownload"):
			utils.MultiDownload(c, basePath, username, password)
		case strings.HasPrefix(c.Request.URL.Path, "/api/login"):
			utils.Login(c)
		case strings.HasPrefix(c.Request.URL.Path, "/api/auth"):
			utils.Auth(c, username, password)
		default:
			utils.StaticHandler(c, staticFiles)
		}
	})

	r.Run(":8080")
}
