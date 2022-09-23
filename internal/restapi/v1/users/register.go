package users

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
	"net/mail"
	"regexp"
)

type RegisterRequest struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	RegisterToken string `json:"register_token"`
}

func Register(c *gin.Context) {
	ctrl := UserController{
		DB: db.DB,
	}
	ctrl.Register(c)
}

func (ctrl UserController) Register(c *gin.Context) {
	role, _ := c.Get("role")

	bodyRaw, err := c.GetRawData()

	if err != nil {
		logger.Errorf("Register Get Raw Data Failed:\n%v", err)
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	var registerRequest RegisterRequest
	json.Unmarshal(bodyRaw, &registerRequest)

	if registerRequest.Username != "" && registerRequest.Email != "" && registerRequest.Password != "" {
		_, err = mail.ParseAddress(registerRequest.Email)
		isEmailValid := err == nil

		regexUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]{4,20}$`)
		correctUsername := regexUsername.Match([]byte(fmt.Sprintf("%s", registerRequest.Username)))
		corectPasswordLength := (len(fmt.Sprintf("%s", registerRequest.Password)) >= 8)

		if isEmailValid && corectPasswordLength && correctUsername {
			hasher := sha256.New()
			hasher.Write([]byte(registerRequest.Password))
			hashedPassword := hex.EncodeToString(hasher.Sum(nil))

			if err != nil {
				logger.Errorf("Register Get Role Failed:\n%v", err)
				c.JSON(http.StatusInternalServerError, "Internal Server Error")
				return
			}

			// Insert User Into Database
			user := models.User{
				Username: registerRequest.Username,
				Email:    registerRequest.Email,
				Password: hashedPassword,
				Role:     role.(models.Role),
			}
			err = ctrl.DB.Create(&user).Error

			var mysqlError *mysql.MySQLError
			if errors.As(err, &mysqlError) {
				if mysqlError.Number == 1062 {
					c.JSON(http.StatusConflict, "Conflict")
					return
				}
			} else if err != nil {
				logger.Errorf("Register Create User Failed:\n%v", err)
				c.JSON(http.StatusInternalServerError, "Internal Server Error")
				return
			} else {
				c.JSON(http.StatusOK, "OK")
				return
			}
		}
	}

	c.JSON(http.StatusBadRequest, "Bad Request")
}
