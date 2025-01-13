package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func TokenCheck(username string, password string, token string) bool {
	if len(username) == 0 && len(password) == 0 {
		return true
	}
	hash := sha256.New()
	hash.Write([]byte(username + password))
	hashedData := hash.Sum(nil)
	hashString := hex.EncodeToString(hashedData)
	return token == hashString
}

func Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"ok":  true,
		"msg": "",
	})
}

func Auth(c *gin.Context, username string, password string) {
	c.JSON(200, gin.H{
		"useAuth": !(len(username) == 0 && len(password) == 0),
	})
}
