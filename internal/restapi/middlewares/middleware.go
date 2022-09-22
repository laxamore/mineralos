package middlewares

import (
	"github.com/laxamore/mineralos/internal/db"
)

type MiddlewareController struct {
	DB  db.IDB
	RDB db.IRedis
}
