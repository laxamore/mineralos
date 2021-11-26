package main

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/laxamore/mineralos/db"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/zmq/server/router/router"
	zmq "github.com/pebbe/zmq4"
)

//  The main task is a load-balancer with heartbeating on workers so we
//  can detect crashed or blocked worker tasks:
func main() {
	// load .env file
	err := godotenv.Load()
	utils.CheckErr(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := db.MongoClient(ctx)
	utils.CheckErr(err)

	//  Start authentication engine
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	defer zmq.AuthStop()
	zmq.AuthAllow("*", "0.0.0.0/0")

	//  Tell the authenticator to allow any CURVE requests for this domain
	zmq.AuthCurveAdd("*", zmq.CURVE_ALLOW_ANY)
	server_secret := os.Getenv("SERVER_SECRET")

	clientRtr, _ := zmq.NewSocket(zmq.ROUTER)
	clientRtr.SetRouterHandover(true)
	clientRtr.ServerAuthCurve("*", server_secret)
	defer clientRtr.Close()

	workerRtr, _ := zmq.NewSocket(zmq.ROUTER)
	workerRtr.ServerAuthCurve("*", server_secret)
	defer workerRtr.Close()

	clientRtr.Bind("tcp://*:9000") //  For clients
	workerRtr.Bind("tcp://*:9001") //  For workers

	//  List of available workers
	workers := make([]router.Worker_t, 0)

	poller1 := zmq.NewPoller()
	poller1.Add(workerRtr, zmq.POLLIN)
	poller2 := zmq.NewPoller()
	poller2.Add(workerRtr, zmq.POLLIN)
	poller2.Add(clientRtr, zmq.POLLIN)

	repo := db.MongoDB{}
	cntrl := router.RouterController{
		DATABASE_NAME: os.Getenv("PROJECT_NAME"),
	}

	for {
		workers, _, _, _ = cntrl.Router(clientRtr, workerRtr, workers, poller1, poller2, client, repo)
	}
}
