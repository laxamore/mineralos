package main

const (
	port = ":9000"
)

func grpcServer() {
	//// load .env file
	//err := godotenv.Load()
	//utils.CheckErr(err)
	//
	////ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	////defer cancel()
	////client, err := db.MongoClient(ctx)
	////utils.CheckErr(err)
	//
	//var client *mongo.Client
	//
	//lis, err := net.Listen("tcp", port)
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//
	//repo := db.MongoDB{}
	//
	//creds, _ := embeding.LoadServerTLSCert()
	//s := grpc.NewServer(grpc.Creds(creds))
	//pb.RegisterMineralosServer(s, &server.ServerController{
	//	Client:              client,
	//	RepositoryInterface: repo,
	//})
	//
	//log.Printf("server listening at %v", lis.Addr())
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
}
