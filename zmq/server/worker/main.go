//
//  Paranoid Pirate worker.
//

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/utils/Log"
	"github.com/laxamore/mineralos/zmq/server/worker/worker"
	zmq "github.com/pebbe/zmq4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func handlePayload(msg []string, mongoClient *mongo.Client, repositoryInterface worker.WorkerRepositoryInterface) ([]byte, error) {
	type Payload struct {
		RigID  string
		Status interface{}
	}

	type Reply struct {
		Status string
		Config interface{}
	}

	var clientPayload Payload
	json.Unmarshal([]byte(msg[2]), &clientPayload)

	res := repositoryInterface.FindOne(mongoClient, os.Getenv("PROJECT_NAME"), "rigs", bson.D{
		{
			Key: "rig_id", Value: clientPayload.RigID,
		},
	})

	if len(res) > 0 {
		Log.Printf("Info: Got Client Payload: %v", msg[2])

		update := bson.D{
			{
				Key: "$set", Value: bson.M{"lastActivity": time.Now().UTC()},
			},
			{
				Key: "$set", Value: bson.M{"status": clientPayload.Status},
			},
		}

		_, err := repositoryInterface.UpdateOne(mongoClient, os.Getenv("PROJECT_NAME"), "rigs", bson.D{
			{
				Key: "rig_id", Value: clientPayload.RigID,
			},
		}, update)

		if err != nil {
			Log.Printf("error %v", err)
		}

		replyMsg := Reply{
			Status: "ok",
			Config: res["conf"],
		}
		replyMsgMarshal, _ := json.Marshal(replyMsg)

		return replyMsgMarshal, nil
	}

	return nil, fmt.Errorf("error: rig_id not found")
}

func main() {
	// load .env file
	err := godotenv.Load()
	utils.CheckErr(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	defer zmq.AuthStop()

	//  Tell the authenticator to allow any CURVE requests for this domain
	zmq.AuthCurveAdd("*", zmq.CURVE_ALLOW_ANY)

	worker_public, worker_secret, err := zmq.NewCurveKeypair()
	utils.CheckErr(err)

	ROUTER_ENDPOINT := "127.0.0.1"
	if os.Getenv("DOCKER") == "true" {
		ROUTER_ENDPOINT = "zmq_router"
	}

	repo := db.MongoDB{}
	cntrl := worker.WorkerController{
		HEARTBEAT_LIVENESS: 3,
		HEARTBEAT_INTERVAL: 1000 * time.Millisecond,
		INTERVAL_INIT:      1000 * time.Millisecond,
		INTERVAL_MAX:       16000 * time.Millisecond,

		//  Paranoid Pirate Protocol constants
		PPP_READY:     "\001",
		PPP_HEARTBEAT: "\002",

		// Connection Settings
		CONNECTION_ENDPOINT: fmt.Sprintf("tcp://%s:9001", ROUTER_ENDPOINT),
		SERVER_PUB_KEY:      "83<s>=wXS9RXKPR4wp:45?Pmo!y>R!qAy%^:^dDl",

		WORKER_KEY:     worker_secret,
		WORKER_PUB_KEY: worker_public,
	}

	worker, poller := cntrl.S_worker_socket(worker_public, worker_secret)

	//  If liveness hits zero, queue is considered disconnected
	liveness := cntrl.HEARTBEAT_LIVENESS
	interval := cntrl.INTERVAL_INIT

	for {
		cntrl.Worker(&worker, &poller, &liveness, &interval, handlePayload, client, repo)
	}
}
