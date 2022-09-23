package users

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/db/models"
	JWT "github.com/laxamore/mineralos/internal/jwt"
	"github.com/laxamore/mineralos/internal/logger"
	"net/http"
	"net/url"
	"os"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	ctrl := UserController{
		DB: db.DB,
	}
	ctrl.Login(c)
}

func (ctrl UserController) Login(c *gin.Context) {
	bodyRaw, err := c.GetRawData()
	if err != nil {
		logger.Errorf("Login Get Raw Data Failed:\n%v", err)
		c.JSON(http.StatusBadRequest, "Bad Request")
	}

	var loginRequest LoginRequest
	json.Unmarshal(bodyRaw, &loginRequest)

	if loginRequest.Username == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s", loginRequest.Password)))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	loginUser := models.User{}
	err = ctrl.DB.First(&loginUser, "username = ? AND password = ?", loginRequest.Username, hashedPassword).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	if loginUser != (models.User{}) {
		exp := config.Config.ACCESS_TOKEN_EXPIRED + time.Now().Unix()
		expRT := config.Config.REFRESH_TOKEN_EXPIRED + time.Now().Unix()

		// Create the Claims for token
		claims := JWT.LoginClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: exp,
				IssuedAt:  time.Now().Unix(),
			},
			Username:       loginUser.Username,
			Email:          loginUser.Email,
			Role:           &loginUser.Role,
			IsRefreshToken: false,
		}
		token, err := JWT.SignJWT(&claims)
		if err != nil {
			logger.Errorf("Login Sign JWT Failed:\n%v", err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// Create the claims for refresh token
		claims = JWT.LoginClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expRT,
				IssuedAt:  time.Now().Unix(),
			},
			Username:       loginUser.Username,
			Email:          loginUser.Email,
			Role:           &loginUser.Role,
			IsRefreshToken: true,
		}
		rtoken, err := JWT.SignJWT(&claims)

		if err != nil {
			logger.Errorf("Login Sign JWT Failed:\n%v", err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "rtoken",
			Value:    url.QueryEscape(rtoken),
			MaxAge:   2592000,
			Path:     "/",
			Domain:   os.Getenv("DOMAIN"),
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
			HttpOnly: true,
		})

		c.JSON(http.StatusOK, gin.H{
			"jwt_token":        token,
			"jwt_token_expiry": exp,
			"r_token":          rtoken,
			"r_token_expiry":   expRT,
		})
	} else {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
	}
}
