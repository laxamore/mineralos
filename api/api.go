package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/laxamore/mineralos/api/ApiRigs"
	"github.com/laxamore/mineralos/api/ApiUsers"
	"github.com/laxamore/mineralos/api/Middleware"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
)

func main() {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Panicf("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	router := gin.Default()

	router.Use(Middleware.CORSMiddleware())

	router.POST("/api/v1/newRig", Middleware.CheckAuth, ApiRigs.NewRig(client))
	router.DELETE("/api/v1/deleteRig", Middleware.CheckAuth, ApiRigs.DeleteRig(client))
	router.GET("/api/v1/getRigs", Middleware.CheckAuth, ApiRigs.GetRigs(client))
	router.GET("/api/v1/getRig/:rig_id", Middleware.CheckAuth, ApiRigs.GetRig(client))

	router.POST("/api/v1/newWallet", Middleware.CheckAuth, ApiRigs.NewWallet(client))
	router.DELETE("/api/v1/deleteWallet", Middleware.CheckAuth, ApiRigs.DeleteWallet(client))
	router.GET("/api/v1/getWallets", Middleware.CheckAuth, ApiRigs.GetWallets(client))

	router.PUT("/api/v1/updateOC", Middleware.CheckAuth, ApiRigs.UpdateOC(client))

	router.POST("/api/v1/register", Middleware.BeforeRegister, ApiUsers.Register(client))
	router.POST("/api/v1/registerToken", Middleware.VerifyAdmin, ApiUsers.RegisterToken(client))
	router.POST("/api/v1/login", ApiUsers.Login(client))
	router.POST("/api/v1/refreshToken", ApiUsers.RefreshToken(client))
	router.POST("/api/v1/logout", Middleware.CheckAuth, ApiUsers.Logout)

	router.Run(":5000")
}
