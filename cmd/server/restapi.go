package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db/models"
	v1 "github.com/laxamore/mineralos/internal/restapi/v1"
	"github.com/laxamore/mineralos/internal/restapi/v1/middlewares"
	"github.com/laxamore/mineralos/internal/restapi/v1/rigs"
	"github.com/laxamore/mineralos/internal/restapi/v1/users"
)

func restApi() {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware)

	router.POST(v1.BASE_PATH+"/newRig", middlewares.CheckAuthRole(&models.RoleOperator), rigs.NewRig)
	//router.DELETE(v1.BASE_PATH+"/deleteRig", middlewares.CheckAuth, rigs.DeleteRig)
	router.GET(v1.BASE_PATH+"/getRigs", middlewares.CheckAuth, rigs.GetRig)
	//router.GET("/api/v1/getRig/:rig_id", middlewares.CheckAuth, rigs.GetRig(client))
	//
	//router.POST("/api/v1/newWallet", middlewares.CheckAuth, rigs.NewWallet(client))
	//router.DELETE("/api/v1/deleteWallet", middlewares.CheckAuth, rigs.DeleteWallet(client))
	//router.GET("/api/v1/getWallets", middlewares.CheckAuth, rigs.GetWallets(client))
	//
	//router.PUT("/api/v1/updateOC", middlewares.CheckAuth, rigs.UpdateOC(client))
	//
	router.POST(v1.BASE_PATH+"/register", middlewares.BeforeRegister, users.Register)
	//router.POST("/api/v1/registerToken", middlewares.VerifyAdmin, users.RegisterToken(client))
	router.POST(v1.BASE_PATH+"/login", users.Login)
	//router.POST("/api/v1/refreshToken", users.RefreshToken(client))
	//router.POST("/api/v1/logout", middlewares.CheckAuth, users.Logout)

	router.Run(fmt.Sprintf(":%d", config.Config.REST_API_PORT))
}
