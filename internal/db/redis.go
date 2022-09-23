package db

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/laxamore/mineralos/config"
)

var (
	RDB *redis.Client
)

type IRedis interface {
	redis.Cmdable
}

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.REDIS_HOST, config.Config.REDIS_PORT),
		Password: config.Config.REDIS_PASSWORD, // no password set
		Username: config.Config.REDIS_USER,
		DB:       0, // use default DB
	})
}
