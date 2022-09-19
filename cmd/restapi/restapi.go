package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/databases"
	"github.com/laxamore/mineralos/internal/databases/models/user"
	"github.com/laxamore/mineralos/internal/restapi/middleware"
	"github.com/laxamore/mineralos/internal/restapi/rig"
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

	router := gin.Default()
	router.Use(middleware.CORSMiddleware)

	router.POST("/api/v1/newRig", middleware.CheckAuthRole(&user.RoleAdmin), rig.NewRig)
	router.DELETE("/api/v1/deleteRig", middleware.CheckAuth, rig.DeleteRig)
	//router.GET("/api/v1/getRigs", middleware.CheckAuth, rig.GetRigs(client))
	//router.GET("/api/v1/getRig/:rig_id", middleware.CheckAuth, rig.GetRig(client))
	//
	//router.POST("/api/v1/newWallet", middleware.CheckAuth, rig.NewWallet(client))
	//router.DELETE("/api/v1/deleteWallet", middleware.CheckAuth, rig.DeleteWallet(client))
	//router.GET("/api/v1/getWallets", middleware.CheckAuth, rig.GetWallets(client))
	//
	//router.PUT("/api/v1/updateOC", middleware.CheckAuth, rig.UpdateOC(client))
	//
	//router.POST("/api/v1/register", middleware.BeforeRegister, users.Register(client))
	//router.POST("/api/v1/registerToken", middleware.VerifyAdmin, users.RegisterToken(client))
	//router.POST("/api/v1/login", users.Login(client))
	//router.POST("/api/v1/refreshToken", users.RefreshToken(client))
	//router.POST("/api/v1/logout", middleware.CheckAuth, users.Logout)

	router.Run(fmt.Sprintf(":%d", config.Config.REST_API_PORT))
}
