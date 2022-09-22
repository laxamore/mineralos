package jwt

import (
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTCustomClaims struct {
	jwt.StandardClaims
	Username string       `json:"username"`
	Email    string       `json:"email"`
	Role     *models.Role `json:"role"`
}

func SignJWT(claims JWTCustomClaims, exp int64) (string, error) {
	claims.ExpiresAt = exp
	claims.IssuedAt = time.Now().Unix()

	signToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := signToken.SignedString([]byte(config.Config.JWT_SECRET))

	return token, err
}

func VerifyJWT(token string) (claims JWTCustomClaims, tokenParsed *jwt.Token, err error) {
	tokenParsed, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWT_SECRET), nil
	})
	return claims, tokenParsed, err
}
