package Middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/utils/JWT"
	"github.com/laxamore/mineralos/utils/Log"
)

func CheckAuth(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")

	if bearerToken != "" {
		token := strings.Split(bearerToken, " ")[1]

		tokenClaims, tokenParsed, err := JWT.VerifyJWT(token)

		if tokenParsed == nil || err != nil {
			Log.Printf("Couldn't handle this token: %v", err)
		} else if tokenParsed.Valid {
			Log.Print("Token Valid")
			c.Set("tokenClaims", tokenClaims)
			c.Writer.WriteHeader(http.StatusOK)
			return
		} else {
			Log.Printf("Couldn't handle this token: %v", err)
		}
	}

	c.Abort()
	c.Writer.WriteHeader(http.StatusUnauthorized)
	c.Writer.Write([]byte("Unauthorized"))
}
