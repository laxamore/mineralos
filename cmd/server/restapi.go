package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/restapi/middlewares"
	"github.com/laxamore/mineralos/internal/restapi/rigs"
	"github.com/laxamore/mineralos/internal/restapi/users"
)

func restApi() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware)

	router.POST("/api/v1/newRig", middlewares.CheckAuthRole(&models.RoleAdmin), rigs.NewRig)
	router.DELETE("/api/v1/deleteRig", middlewares.CheckAuth, rigs.DeleteRig)
	//router.GET("/api/v1/getRigs", middlewares.CheckAuth, rigs.GetRigs(client))
	//router.GET("/api/v1/getRig/:rig_id", middlewares.CheckAuth, rigs.GetRig(client))
	//
	//router.POST("/api/v1/newWallet", middlewares.CheckAuth, rigs.NewWallet(client))
	//router.DELETE("/api/v1/deleteWallet", middlewares.CheckAuth, rigs.DeleteWallet(client))
	//router.GET("/api/v1/getWallets", middlewares.CheckAuth, rigs.GetWallets(client))
	//
	//router.PUT("/api/v1/updateOC", middlewares.CheckAuth, rigs.UpdateOC(client))
	//
	router.POST("/api/v1/register", middlewares.BeforeRegister, users.Register)
	//router.POST("/api/v1/registerToken", middlewares.VerifyAdmin, users.RegisterToken(client))
	//router.POST("/api/v1/login", users.Login(client))
	//router.POST("/api/v1/refreshToken", users.RefreshToken(client))
	//router.POST("/api/v1/logout", middlewares.CheckAuth, users.Logout)

	router.Run(fmt.Sprintf(":%d", config.Config.REST_API_PORT))
}
