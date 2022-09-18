package middleware

import (
	JWT "github.com/laxamore/mineralos/internal/jwt"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type VerifyAdminController struct{}

func (a VerifyAdminController) TryVerifyAdmin(c *gin.Context) {
	c.Set("admin", false)
	token := c.GetHeader("Token")

	claims, tokenParsed, err := JWT.VerifyJWT(token)

	if tokenParsed == nil || err != nil {
		logger.Printf("Couldn't handle this token: %v", err)
	} else if tokenParsed.Valid {
		if claims["privilege"] == "admin" {
			c.Set("admin", true)
			c.Writer.WriteHeader(http.StatusOK)
			return
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			logger.Print("Couldn't handle this token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			logger.Print("Token Expired.")
		} else {
			logger.Printf("Couldn't handle this token: %v", err)
		}
	} else {
		logger.Printf("Couldn't handle this token: %v", err)
	}

	c.Abort()
	c.Writer.WriteHeader(http.StatusUnauthorized)
	c.Writer.Write([]byte("Unauthorized"))
}

func VerifyAdmin(c *gin.Context) {
	cntrl := VerifyAdminController{}
	cntrl.TryVerifyAdmin(c)
}
