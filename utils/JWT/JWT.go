package JWT

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func SignJWT(claims jwt.MapClaims, exp int64) (string, error) {
	claims["exp"] = exp
	claims["iat"] = time.Now().Unix()

	signToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := signToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return token, err
}

func VerifyJWT(token string) (jwt.MapClaims, *jwt.Token, error) {
	claims := jwt.MapClaims{}
	tokenParsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return claims, tokenParsed, err
}
