package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/joho/godotenv"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/grpc/handler/server"
	pb "github.com/laxamore/mineralos/grpc/mineralos_proto"
	"github.com/laxamore/mineralos/utils"
	"google.golang.org/grpc"
)

const (
	port = ":9000"
)

func main() {
	// load .env file
	err := godotenv.Load()
	utils.CheckErr(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := db.MongoDB{}
	s := grpc.NewServer()
	pb.RegisterMineralosServer(s, &server.ServerController{
		Client:              client,
		RepositoryInterface: repo,
	})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
