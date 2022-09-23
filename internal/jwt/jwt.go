package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db/models"
)

type LoginClaims struct {
	StandardClaims
	Username       string       `json:"username"`
	Email          string       `json:"email"`
	Role           *models.Role `json:"role"`
	IsRefreshToken bool         `json:"is_refresh_token"`
}

type StandardClaims struct {
	jwt.StandardClaims
}

type Claims interface {
	jwt.Claims
}

type Token struct {
	jwt.Token
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type Cmd interface {
	SignJWT(claims Claims) (string, error)
	VerifyJWT(token string, claims Claims) (tokenParsed *Token, err error)
}

func (j Service) SignJWT(claims Claims) (string, error) {
	signToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := signToken.SignedString([]byte(config.Config.JWT_SECRET))

	return token, err
}

func (j Service) VerifyJWT(token string, claims Claims) (*Token, error) {
	tokenParsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWT_SECRET), nil
	})

	return &Token{*tokenParsed}, err
}
