package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/databases"
	"github.com/laxamore/mineralos/internal/restapi/middleware"
	"github.com/laxamore/mineralos/internal/restapi/rigs"
	"github.com/laxamore/mineralos/internal/restapi/users"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {
	var err error
	config.Config, err = config.LoadConfig()
	if err != nil {
		log.Panic("Load Config Error ", err)
	}

	// Init Database
	databases.Connect(config.Config.DB_USER, config.Config.DB_PASSWORD, config.Config.DB_HOST, config.Config.DB_PORT, config.Config.DB_NAME)
	databases.InitDatabase()

	var client *mongo.Client

	router := gin.Default()

	router.Use(middleware.CORSMiddleware)

	router.POST("/api/v1/newRig", middleware.CheckAuthPrivilege("admin"), rigs.NewRig())
	router.DELETE("/api/v1/deleteRig", middleware.CheckAuth, rigs.DeleteRig(client))
	router.GET("/api/v1/getRigs", middleware.CheckAuth, rigs.GetRigs(client))
	router.GET("/api/v1/getRig/:rig_id", middleware.CheckAuth, rigs.GetRig(client))

	router.POST("/api/v1/newWallet", middleware.CheckAuth, rigs.NewWallet(client))
	router.DELETE("/api/v1/deleteWallet", middleware.CheckAuth, rigs.DeleteWallet(client))
	router.GET("/api/v1/getWallets", middleware.CheckAuth, rigs.GetWallets(client))

	router.PUT("/api/v1/updateOC", middleware.CheckAuth, rigs.UpdateOC(client))

	router.POST("/api/v1/register", middleware.BeforeRegister, users.Register(client))
	router.POST("/api/v1/registerToken", middleware.VerifyAdmin, users.RegisterToken(client))
	router.POST("/api/v1/login", users.Login(client))
	router.POST("/api/v1/refreshToken", users.RefreshToken(client))
	router.POST("/api/v1/logout", middleware.CheckAuth, users.Logout)

	router.Run(fmt.Sprintf(":%d", config.Config.REST_API_PORT))
}
