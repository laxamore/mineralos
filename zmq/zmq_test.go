package main

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/laxamore/mineralos/utils"
	"github.com/laxamore/mineralos/utils/Log"
	"github.com/laxamore/mineralos/zmq/client/client"
	"github.com/laxamore/mineralos/zmq/server/router/router"
	"github.com/laxamore/mineralos/zmq/server/worker/worker"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type ZMQRepositoryMock struct {
	mock.Mock
}

var quitRouterHandler chan bool
var quitWorkerHandler chan bool
var quitClientHandler chan bool

var clientResult chan []string

func (a ZMQRepositoryMock) FindOne(client *mongo.Client, db_name string, collection_name string, filter interface{}) map[string]interface{} {
	return map[string]interface{}{
		"test": "test",
	}
}

func (a ZMQRepositoryMock) UpdateOne(client *mongo.Client, db_name string, collection_name string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return nil, nil
}

func routerHandler() {
	server_secret := "2v-hp%0n^HwLK7I+S-)=N0*Z)LPzG&AX#fUaZ]d*"

	clientRtr, _ := zmq.NewSocket(zmq.ROUTER)
	clientRtr.ServerAuthCurve("*", server_secret)
	defer clientRtr.Close()

	workerRtr, _ := zmq.NewSocket(zmq.ROUTER)
	workerRtr.ServerAuthCurve("*", server_secret)
	defer workerRtr.Close()

	clientRtr.Bind("tcp://*:9003") //  For clients
	workerRtr.Bind("tcp://*:9004") //  For workers

	//  List of available workers
	workers := make([]router.Worker_t, 0)

	poller1 := zmq.NewPoller()
	poller1.Add(workerRtr, zmq.POLLIN)
	poller2 := zmq.NewPoller()
	poller2.Add(workerRtr, zmq.POLLIN)
	poller2.Add(clientRtr, zmq.POLLIN)

	repo := ZMQRepositoryMock{}
	cntrl := router.RouterController{
		DATABASE_NAME: "mineralos",
	}

	Log.Print("Info: Starting Router...")

	for {
		select {
		case <-quitRouterHandler:
			return
		default:
			workers, _, _, _ = cntrl.Router(clientRtr, workerRtr, workers, poller1, poller2, nil, repo)
		}
	}
}

func workerHandler() {
	worker_public, worker_secret, err := zmq.NewCurveKeypair()
	utils.CheckErr(err)

	repo := ZMQRepositoryMock{}
	cntrl := worker.WorkerController{
		HEARTBEAT_LIVENESS: 3,
		HEARTBEAT_INTERVAL: 1000 * time.Millisecond,
		INTERVAL_INIT:      1000 * time.Millisecond,
		INTERVAL_MAX:       16000 * time.Millisecond,

		//  Paranoid Pirate Protocol constants
		PPP_READY:     "\001",
		PPP_HEARTBEAT: "\002",

		// Connection Settings
		CONNECTION_ENDPOINT: "tcp://127.0.0.1:9004",
		SERVER_PUB_KEY:      "83<s>=wXS9RXKPR4wp:45?Pmo!y>R!qAy%^:^dDl",

		WORKER_KEY:     worker_secret,
		WORKER_PUB_KEY: worker_public,
	}

	workerSocket, poller := cntrl.S_worker_socket(worker_public, worker_secret)

	//  If liveness hits zero, queue is considered disconnected
	liveness := cntrl.HEARTBEAT_LIVENESS
	interval := cntrl.INTERVAL_INIT

	Log.Print("Info: Starting Worker...")

	for {
		select {
		case <-quitWorkerHandler:
			return
		default:
			cntrl.Worker(&workerSocket, &poller, &liveness, &interval, func(s1 []string, c *mongo.Client, wri worker.WorkerRepositoryInterface) ([]byte, error) {
				return []byte("ok"), nil
			}, nil, repo)
		}
	}
}

func clientHandler() (res string) {
	cntrl := client.ClientController{
		REQUEST_TIMEOUT: 2500 * time.Millisecond, //  msecs, (> 1000!)
		SERVER_ENDPOINT: "tcp://127.0.0.1:9003",

		HEARTBEAT_INTERVAL: 1000 * time.Millisecond, //  msecs
		RIG_ID:             "b073badf-c10f-4a94-9d04-e0bb35b26d18",
		ClientKey:          "xOJ)9GNipQOkJWd7k^m:&fF9BSKlu0v73#JCyana",
		ClientPubKey:       "i!OLA?M*RxcyAwZ{#Kvn+ri^F3x-H-lrki=5n6xP",
	}

	client, poller, err := cntrl.NewClientConnection(cntrl.ClientPubKey, cntrl.ClientKey)
	if err != nil {
		panic(err)
	}

	Log.Print("Info: Starting Client...")

	for {
		select {
		case <-quitClientHandler:
			clientResult <- []string{"timeout"}
			return
		default:
			lastPayload, workerResponse, err := cntrl.Client(client, poller)

			if err != nil {
				Log.Printf("%v", err)
				Log.Printf("info: retrying...")

				//  Old socket is confused; close it and open a new one
				client.Close()
				client, poller, _ = cntrl.NewClientConnection(cntrl.ClientPubKey, cntrl.ClientKey)

				//  Send request again, on new socket
				client.SendMessage(lastPayload)
			}

			if len(workerResponse) > 0 {
				quitRouterHandler <- true
				quitWorkerHandler <- true
				clientResult <- workerResponse
				return
			}
		}
	}
}

func TestZMQ(t *testing.T) {
	timeout := time.After(10 * time.Second)

	//  Start authentication engine
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	defer zmq.AuthStop()
	zmq.AuthAllow("*", "0.0.0.0/0")

	//  Tell the authenticator to allow any CURVE requests for this domain
	zmq.AuthCurveAdd("*", zmq.CURVE_ALLOW_ANY)

	quitRouterHandler = make(chan bool)
	quitWorkerHandler = make(chan bool)
	quitClientHandler = make(chan bool)

	clientResult = make(chan []string)

	go routerHandler()
	go workerHandler()
	go clientHandler()

	select {
	case <-timeout:
		quitRouterHandler <- true
		quitWorkerHandler <- true
		quitClientHandler <- true
		t.FailNow()
	case clientResultData := <-clientResult:
		Log.Print(clientResultData[0])
		assert.Equal(t, clientResultData[0], "ok")
	}
}
