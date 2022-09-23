package users

import (
	"github.com/laxamore/mineralos/internal/db"
	"github.com/laxamore/mineralos/internal/jwt"
)

type UserController struct {
	DB         db.IDB
	JWTService jwt.Cmd
}
