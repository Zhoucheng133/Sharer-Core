package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func TokenCheck(username string, password string, token string) bool {
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
