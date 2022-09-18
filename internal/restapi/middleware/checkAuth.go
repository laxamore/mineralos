package middleware

import (
	"github.com/laxamore/mineralos/internal/jwt"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	CheckAuthPrivilege("")
}

func CheckAuthPrivilege(privilege string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if bearerToken != "" {
			token := strings.Split(bearerToken, " ")[1]
			tokenClaims, tokenParsed, err := jwt.VerifyJWT(token)

			if privilege == "" || tokenClaims["privilege"] == privilege {
				if tokenParsed == nil || err != nil {
					logger.Printf("Couldn't handle this token: %v", err)
				} else if tokenParsed.Valid {
					logger.Print("Token Valid")
					c.Set("tokenClaims", tokenClaims)
					c.Writer.WriteHeader(http.StatusOK)
					return
				} else {
					logger.Printf("Couldn't handle this token: %v", err)
				}
			}
		}

		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
	}
}
