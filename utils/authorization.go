package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func Token(c *gin.Context, username string, password string, tokenString string, secret string) {
	c.JSON(200, gin.H{
		"ok":  TokenCheck(username, password, tokenString, secret),
		"msg": "",
	})
}

func TokenCheck(username string, password string, tokenString string, secret string) bool {
	if len(username) == 0 && len(password) == 0 {
		return true
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return err == nil && token.Valid && claims.Username == username
}

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context, username string, password string, secret string) {

	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"ok":  false,
			"msg": "Invalid JSON",
		})
		return
	}

	if user.Username == username && user.Password == password {

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Username: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(500, gin.H{
				"ok":  false,
				"msg": "Internal server error",
			})
			return
		}

		c.JSON(200, gin.H{
			"ok":  true,
			"msg": tokenString,
		})
	} else {
		c.JSON(200, gin.H{
			"ok":  false,
			"msg": "Invalid username or password",
		})
	}
}
func Auth(c *gin.Context, username string, password string) {
	c.JSON(200, gin.H{
		"useAuth": !(len(username) == 0 && len(password) == 0),
	})
}
