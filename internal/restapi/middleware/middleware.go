package middleware

import "github.com/laxamore/mineralos/internal/databases"

type MiddlewareController struct {
	DB databases.DBInterface
}
