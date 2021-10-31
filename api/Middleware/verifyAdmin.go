package Middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func VerifyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")

		claims := jwt.MapClaims{}
		tokenParsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if tokenParsed == nil || err != nil {
			log.Print("Couldn't handle this token: ", err)
		} else if tokenParsed.Valid {
			if claims["privilege"] == "admin" {
				c.Set("admin", true)
				return
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				log.Print("Couldn't handle this token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				log.Print("Token Expired.")
			} else {
				log.Print("Couldn't handle this token:", err)
			}
		} else {
			log.Print("Couldn't handle this token:", err)
		}

		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
	}
}
