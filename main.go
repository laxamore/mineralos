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

	router.POST("/newrig", ApiRigs.NewRig)
	router.DELETE("/deleterig", ApiRigs.DeleteRig)
	router.POST("/register", Middleware.BeforeRegister(), ApiUsers.Register)
	router.POST("/registerToken", Middleware.VerifyAdmin(), ApiUsers.RegisterToken)
	router.POST("/login", ApiUsers.Login)

	router.Run(":5000")
}
