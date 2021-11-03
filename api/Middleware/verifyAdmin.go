package Middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	Log "github.com/laxamore/mineralos/log"
)

type VerifyAdminController struct{}

func (a VerifyAdminController) TryVerifyAdmin(c *gin.Context) {
	c.Set("admin", false)
	token := c.GetHeader("Token")

	claims := jwt.MapClaims{}
	tokenParsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if tokenParsed == nil || err != nil {
		Log.Printf("Couldn't handle this token: %v", err)
	} else if tokenParsed.Valid {
		if claims["privilege"] == "admin" {
			c.Set("admin", true)
			c.Writer.WriteHeader(http.StatusOK)
			return
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			Log.Print("Couldn't handle this token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			Log.Print("Token Expired.")
		} else {
			Log.Printf("Couldn't handle this token: %v", err)
		}
	} else {
		Log.Printf("Couldn't handle this token: %v", err)
	}

	c.Abort()
	c.Writer.WriteHeader(http.StatusUnauthorized)
	c.Writer.Write([]byte("Unauthorized"))
}

func VerifyAdmin(c *gin.Context) {
	cntrl := VerifyAdminController{}
	cntrl.TryVerifyAdmin(c)
}
