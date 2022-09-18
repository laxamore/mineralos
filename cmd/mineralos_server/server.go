package main

import (
	"github.com/joho/godotenv"
	"github.com/laxamore/mineralos/cmd/mineralos_server/handler/server"
	pb "github.com/laxamore/mineralos/config/mineralos_proto"
	"github.com/laxamore/mineralos/internal/databases"
	"github.com/laxamore/mineralos/internal/embeding"
	"github.com/laxamore/mineralos/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":9000"
)

func main() {
	// load .env file
	err := godotenv.Load()
	utils.CheckErr(err)

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//client, err := databases.MongoClient(ctx)
	//utils.CheckErr(err)

	var client *mongo.Client

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := databases.MongoDB{}

	creds, _ := embeding.LoadServerTLSCert()
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterMineralosServer(s, &server.ServerController{
		Client:              client,
		RepositoryInterface: repo,
	})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
