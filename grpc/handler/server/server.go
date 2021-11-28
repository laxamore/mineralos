package server

import (
	"context"
	"fmt"
	"time"

	"github.com/laxamore/mineralos/grpc/grpc_type"
	pb "github.com/laxamore/mineralos/grpc/mineralos_proto"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/utils/Linux"
	"github.com/laxamore/mineralos/utils/Log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerRepositoryInterface interface {
	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
	UpdateOne(*mongo.Client, string, string, interface{}, interface{}) (*mongo.UpdateResult, error)
}

type ServerController struct {
	pb.UnimplementedMineralosServer
	Client              *mongo.Client
	RepositoryInterface ServerRepositoryInterface
}

func (s *ServerController) ReportStatus(ctx context.Context, in *pb.Payload) (*pb.ServerReply, error) {
	clientPayload := grpc_type.ClientPayload{
		RigID: in.GetRigId(),
		Drivers: Linux.GPUDriverVersion{
			AMD:    in.GetAmdDriver(),
			NVIDIA: in.GetNvidiaDriver(),
		},
	}

	replyMsg, err := handlePayload(clientPayload, s.Client, s.RepositoryInterface)
	utils.CheckErr(err)
	return &pb.ServerReply{Message: replyMsg.Message}, nil
}

func handlePayload(clientPayload grpc_type.ClientPayload, mongoClient *mongo.Client, repositoryInterface ServerRepositoryInterface) (replyMsg grpc_type.Reply, err error) {
	res := repositoryInterface.FindOne(mongoClient, "mineralos", "rigs", bson.D{
		{
			Key: "rig_id", Value: clientPayload.RigID,
		},
	})

	if len(res) > 0 {
		Log.Printf("Info: Got Client Payload: %v", clientPayload)

		update := bson.D{
			{
				Key: "$set", Value: bson.M{"lastActivity": time.Now().UTC()},
			},
			{
				Key: "$set", Value: bson.M{"status": map[string]interface{}{
					"Drivers": clientPayload.Drivers,
				},
				},
			},
		}

		_, err = repositoryInterface.UpdateOne(mongoClient, "mineralos", "rigs", bson.D{
			{
				Key: "rig_id", Value: clientPayload.RigID,
			},
		}, update)

		if err != nil {
			Log.Printf("error %v", err)
		}

		replyMsg.Message = "ok"
		return
	}

	err = fmt.Errorf("error: rig_id not found")
	return
}
