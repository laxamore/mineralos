package main

import (
	"log"
	"mineralos/api/ApiRigs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiServer := gin.Default()

	apiServer.POST("/newrig", ApiRigs.NewRig)
	apiServer.DELETE("/deleterig", ApiRigs.DeleteRig)

	apiServer.Run(":5000")
}
