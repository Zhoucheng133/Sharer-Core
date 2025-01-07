package main

import (
	"virtual-directory-core/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.POST("/api/list", utils.GetList)

	r.Run(":8080")
}
