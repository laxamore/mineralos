package middlewares

import (
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/jwt"
)

type MiddlewareController struct {
	DB         db.IDB
	RDB        db.IRedis
	JWTService jwt.Cmd
}
