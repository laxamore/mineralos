package main

import (
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db"
	"log"
)

func main() {
	var err error
	config.Config, err = config.LoadConfig()
	if err != nil {
		log.Panic("Load Config Error ", err)
	}

	// Init Database
	db.ConnectDB()
	db.InitDB()

	db.ConnectRedis()

	restApi()
	grpcServer()
}
