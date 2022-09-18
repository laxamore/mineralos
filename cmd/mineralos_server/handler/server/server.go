package server

import (
	"context"
	"fmt"
	"github.com/laxamore/mineralos/internal/logger"
	"time"

	pb "github.com/laxamore/mineralos/config/mineralos_proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServerRepositoryInterface interface {
	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
	UpdateOne(*mongo.Client, string, string, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type ServerController struct {
	pb.UnimplementedMineralosServer
	Client              *mongo.Client
	RepositoryInterface ServerRepositoryInterface
}

func (s *ServerController) ReportStatus(ctx context.Context, in *pb.Payload) (*pb.ServerReply, error) {
	replyMsg, err := handlePayload(in, s.Client, s.RepositoryInterface)
	if err != nil {
		logger.Printf("error: %v", err)
	}
	return replyMsg, nil
}

func handlePayload(clientPayload *pb.Payload, mongoClient *mongo.Client, repositoryInterface ServerRepositoryInterface) (replyMsg *pb.ServerReply, err error) {
	res := repositoryInterface.FindOne(mongoClient, "mineralos", "rigs", bson.D{
		{
			Key: "rig_id", Value: clientPayload.GetRigId(),
		},
	})

	if len(res) > 0 {
		logger.Printf("Info: Got Client Payload: %v", clientPayload)

		update := bson.D{
			{
				Key: "$set", Value: bson.M{"lastActivity": time.Now().UTC()},
			},
			{
				Key: "$set", Value: bson.M{
					"status": clientPayload.GetStatus(),
				},
			},
		}

		_, err = repositoryInterface.UpdateOne(mongoClient, "mineralos", "rigs", bson.D{
			{
				Key: "rig_id", Value: clientPayload.GetRigId(),
			},
		}, update)

		if err != nil {
			logger.Printf("error %v", err)
		}

		replyMsg = &pb.ServerReply{Message: "ok"}
		return
	}

	err = fmt.Errorf("error: rig_id not found")
	return
}
