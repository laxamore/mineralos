package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db/models"
)

type LoginClaims struct {
	jwt.StandardClaims
	Username       string       `json:"username"`
	Email          string       `json:"email"`
	Role           *models.Role `json:"role"`
	IsRefreshToken bool         `json:"is_refresh_token"`
}

func SignJWT(claims jwt.Claims) (string, error) {
	signToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := signToken.SignedString([]byte(config.Config.JWT_SECRET))

	return token, err
}

func VerifyJWT(token string, claims jwt.Claims) (tokenParsed *jwt.Token, err error) {
	tokenParsed, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWT_SECRET), nil
	})
	return tokenParsed, err
}
