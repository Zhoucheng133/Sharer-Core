package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"sharer-core/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func requestMiddleware(username string, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

	gin.SetMode(gin.ReleaseMode)

	// -->参数<--
	// -port 8080				端口
	// -d /Users/admin/desktop	路径
	// -u admin					用户 (可以忽略，如果-p值存在则不能忽略)
	// -p 123456				密码 (可以忽略，如果-u值存在则不能忽略)
	// -->参数结束<--

	port := flag.String("port", "", "端口号")
	basePath := flag.String("d", "", "分享路径")
	username := flag.String("u", "", "用户名")
	password := flag.String("p", "", "密码")

	flag.Parse()

	if *port == "" {
		fmt.Println("缺少参数: -port，用于指定服务运行端口")
		os.Exit(1)
	} else if *basePath == "" {
		fmt.Println("缺少参数: -d，用于指定分享路径")
		os.Exit(1)
	} else if *username != "" && *password == "" {
		fmt.Println("缺少参数: -p，用于指定分享的密码")
		os.Exit(1)
	} else if *username == "" && *password != "" {
		fmt.Println("缺少参数: -u，用于指定分享的用户名")
		os.Exit(1)
	}

	// -->测试代码<---
	// username := "admin"
	// password := "123456"
	// useAuth := true
	// 如果username=="" && password==""表明不需要验证
	// basePath := "/Users/zhoucheng/Downloads"
	//-->测试结束<--

	// 所有路径请求path需要添加头/
	r := gin.New()
	r.Use(requestMiddleware(*username, *password))
	r.POST("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/list"):
			utils.GetList(c, *basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/uploadFolder"):
			utils.UploadFolder(c, *basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/upload"):
			utils.Upload(c, *basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/del"):
			utils.DelRequest(c, *basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/mkdir"):
			utils.CreateFolder(c, *basePath)
		case strings.HasPrefix(c.Request.URL.Path, "/api/rename"):
			utils.Rename(c, *basePath)
		}
	})
	r.GET("/*path", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.URL.Path, "/api/raw"):
			utils.GetRaw(c, *basePath, *username, *password)
		case strings.HasPrefix(c.Request.URL.Path, "/api/download"):
			utils.Download(c, *basePath, *username, *password)
		case strings.HasPrefix(c.Request.URL.Path, "/api/multidownload"):
			utils.MultiDownload(c, *basePath, *username, *password)
		case strings.HasPrefix(c.Request.URL.Path, "/api/login"):
			utils.Login(c)
		case strings.HasPrefix(c.Request.URL.Path, "/api/auth"):
			utils.Auth(c, *username, *password)
		default:
			utils.StaticHandler(c, staticFiles)
		}
	})

	fmt.Println(fmt.Sprint("服务运行在: \n➜ http://", utils.GetIp(), ":", *port, "\n➜ http://127.0.0.1:", *port))
	r.Run(fmt.Sprint(":", *port))
}
