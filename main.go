package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/laxamore/mineralos/api/ApiRigs"
	"github.com/laxamore/mineralos/api/ApiUsers"
	"github.com/laxamore/mineralos/api/Middleware"
)

func main() {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Panicf("Error loading .env file")
	}

	router := gin.Default()

	router.Use(Middleware.CORSMiddleware())

	router.POST("/api/v1/newRig", Middleware.CheckAuth, ApiRigs.NewRig)
	router.DELETE("/api/v1/deleteRig", Middleware.CheckAuth, ApiRigs.DeleteRig)
	router.GET("/api/v1/getRigs", Middleware.CheckAuth, ApiRigs.GetRigs)
	router.GET("/api/v1/getRig/:rig_id", Middleware.CheckAuth, ApiRigs.GetRig)

	router.POST("/api/v1/register", Middleware.BeforeRegister, ApiUsers.Register)
	router.POST("/api/v1/registerToken", Middleware.VerifyAdmin, ApiUsers.RegisterToken)
	router.POST("/api/v1/login", ApiUsers.Login)
	router.POST("/api/v1/refreshToken", ApiUsers.RefreshToken)
	router.POST("/api/v1/logout", Middleware.CheckAuth, ApiUsers.Logout)

	router.Run(":5000")
}
