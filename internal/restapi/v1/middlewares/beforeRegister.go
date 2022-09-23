package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/logger"
	"github.com/laxamore/mineralos/internal/restapi/v1/users"
	"gorm.io/gorm"
	"io"
	"net/http"
)

func BeforeRegister(c *gin.Context) {
	ctrl := MiddlewareController{
		DB:  db.DB,
		RDB: db.RDB,
	}
	ctrl.BeforeRegister(c)
}

func (ctrl MiddlewareController) BeforeRegister(c *gin.Context) {
	registerRequest := users.RegisterRequest{}
	bodyRaw, err := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))
	err = json.Unmarshal(bodyRaw, &registerRequest)

	if err != nil {
		logger.Errorf("Before Register Unmarshal Failed:\n%v", err)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	var u models.User
	err = ctrl.DB.First(&u).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Errorf("BeforeRegister Get User Failed:\n%v", err)
		c.Abort()
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
	}

	if u != (models.User{}) {
		// Check Registration Token On Redis
		val, err := ctrl.RDB.Get(c, registerRequest.RegisterToken).Result()
		if err != nil {
			if err == redis.Nil {
				logger.Errorf("BeforeRegister Get RegToken Failed:\n%v", err)
			} else {
				logger.Errorf("BeforeRegister Get Registration Token On Redis Failed:\n%v", err)
			}
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}

		role := models.Role{}
		err = json.Unmarshal([]byte(val), &role)
		if err != nil {
			logger.Errorf("BeforeRegister Unmarshal Role Failed:\n%v", err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
		c.Set("role", role)
	} else {
		c.Set("role", models.RoleAdmin)
	}

	c.Writer.WriteHeader(http.StatusOK)
	return
}
