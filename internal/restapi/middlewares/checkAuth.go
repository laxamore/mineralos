package middlewares

import (
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/jwt"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	CheckAuthRole(&models.RoleUser)
}

func CheckAuthRole(role *models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		middlewareCtrl := MiddlewareController{}
		middlewareCtrl.CheckAuthRole(c, role)
	}
}

func (ctrl MiddlewareController) CheckAuthRole(c *gin.Context, role *models.Role) {
	bearerToken := c.GetHeader("Authorization")

	if bearerToken != "" {
		token := strings.Split(bearerToken, " ")[1]
		tokenClaims, tokenParsed, err := jwt.VerifyJWT(token)

		if err != nil {
			logger.Error("Failed to verify JWT ", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if tokenClaims.Role.Level <= role.Level && tokenParsed.Valid {
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
