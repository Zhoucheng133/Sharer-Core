package utils

import (
	"github.com/gin-gonic/gin"
)

type Body struct {
	Path string `json:"path"`
}

func GetList(c *gin.Context) {
	var body Body

	if err := c.ShouldBind(&body); err == nil {
		if body.Path == "" {
			c.JSON(400, gin.H{
				"ok":  false,
				"msg": "Path parameter is required",
			})
		} else {
			c.JSON(200, gin.H{
				"ok":  true,
				"msg": body.Path,
			})
		}
		return
	}
	c.JSON(400, gin.H{
		"ok":  false,
		"msg": "Bad request",
	})

}
