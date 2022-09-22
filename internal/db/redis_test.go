package db

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/logger"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()

	config.Config = &config.ConfigStruct{
		REDIS_HOST:     "localhost",
		REDIS_PORT:     6377,
		REDIS_PASSWORD: "admin",
	}

	ConnectRedis()

	// struct to json string
	roleAdminString, err := json.Marshal(&models.RoleAdmin)

	RDB.Set(ctx, "nani", roleAdminString, 0)
	val, err := RDB.Get(ctx, "nani").Result()
	if err != nil {
		panic(err)
	}
	logger.Print(val)
	//require.Contains(t, val, "39")

	newRoleAdmin := models.Role{}
	err = json.Unmarshal([]byte(val), &newRoleAdmin)

	logger.Print(newRoleAdmin)

	// Sleep for 3 seconds
	time.Sleep(time.Second * 3)

	val, err = RDB.Get(ctx, "nani").Result()
	if err == redis.Nil {
		require.NotContains(t, val, "39")
	}
}
