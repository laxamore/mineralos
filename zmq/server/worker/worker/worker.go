package worker

import (
	"fmt"
	"time"

	"github.com/laxamore/mineralos/utils/Log"
	zmq "github.com/pebbe/zmq4"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkerRepositoryInterface interface {
	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
	UpdateOne(*mongo.Client, string, string, interface{}, interface{}) (*mongo.UpdateResult, error)
}

type WorkerController struct {
	HEARTBEAT_LIVENESS int           //  3-5 is reasonable
	HEARTBEAT_INTERVAL time.Duration //  msecs
	INTERVAL_INIT      time.Duration //  Initial reconnect
	INTERVAL_MAX       time.Duration //  After exponential backoff

	//  Paranoid Pirate Protocol constants
	PPP_READY     string //  Signals worker is ready
	PPP_HEARTBEAT string //  Signals worker heartbeat

	// Connection Settings
	CONNECTION_ENDPOINT string
	SERVER_PUB_KEY      string
	WORKER_KEY          string
	WORKER_PUB_KEY      string
}

//  Helper function that returns a new configured socket
//  connected to the Paranoid Pirate queue
func (a WorkerController) S_worker_socket(client_public string, client_secret string) (zmq.Socket, zmq.Poller) {
	worker, _ := zmq.NewSocket(zmq.DEALER)
	worker.ClientAuthCurve(a.SERVER_PUB_KEY, client_public, client_secret)
	worker.Connect(a.CONNECTION_ENDPOINT)

	//  Tell queue we're ready for work
	Log.Print("Info: worker ready")
	worker.Send(a.PPP_READY, 0)

	poller := zmq.NewPoller()
	poller.Add(worker, zmq.POLLIN)

	return *worker, *poller
}

//  We have a single task, which implements the worker side of the
//  Paranoid Pirate Protocol (PPP). The interesting parts here are
//  the heartbeating, which lets the worker detect if the queue has
//  died, and vice-versa:
func (a WorkerController) Worker(worker *zmq.Socket, poller *zmq.Poller, liveness *int, interval *time.Duration, handlePayload func([]string, *mongo.Client, WorkerRepositoryInterface) ([]byte, error), client *mongo.Client, repositoryInterface WorkerRepositoryInterface) error {
	//  Send out heartbeats at regular intervals
	heartbeat_at := time.NewTicker(a.HEARTBEAT_INTERVAL)

	sockets, err := poller.Poll(a.HEARTBEAT_INTERVAL)
	if err != nil {
		return fmt.Errorf("worker poller error %v", err)
	}

	if len(sockets) == 1 {
		//  Get message
		//  - 3-part envelope + content -> request
		//  - 1-part HEARTBEAT -> heartbeat
		msg, err := worker.RecvMessage(0)
		if err != nil {
			return fmt.Errorf("worker recv msg error %v", err)
		}

		//  To test the robustness of the queue implementation we //
		//  simulate various typical problems, such as the worker
		//  crashing, or running very slowly. We do this after a few
		//  cycles so that the architecture can get up and running
		//  first:
		if len(msg) == 3 {
			replyMsg, err := handlePayload(msg, client, repositoryInterface)
			if err != nil {
				Log.Printf("%v", err)
			} else {
				msg[2] = string(replyMsg)
				worker.SendMessage(msg)
			}
			*liveness = a.HEARTBEAT_LIVENESS
		} else if len(msg) == 1 {
			//  When we get a heartbeat message from the queue, it means the
			//  queue was (recently) alive, so reset our liveness indicator:
			if msg[0] == a.PPP_HEARTBEAT {
				*liveness = a.HEARTBEAT_LIVENESS
			} else {
				Log.Printf("Error: invalid message: %q\n", msg)
			}
		} else {
			Log.Printf("Error: invalid message: %q\n", msg)
		}
		*interval = a.INTERVAL_INIT
	} else {
		//  If the queue hasn't sent us heartbeats in a while, destroy the
		//  socket and reconnect. This is the simplest most brutal way of
		//  discarding any messages we might have sent in the meantime://
		*liveness--
		if *liveness == 0 {
			Log.Print("Warning: heartbeat failure, can't reach queue")
			Log.Printf("Warning: reconnecting in %d seconds", *interval/time.Second)
			time.Sleep(*interval)

			if *interval < a.INTERVAL_MAX {
				*interval = 2 * *interval
			}

			*worker, *poller = a.S_worker_socket(a.WORKER_PUB_KEY, a.WORKER_KEY)
			*liveness = a.HEARTBEAT_LIVENESS
		}
	}

	//  Send heartbeat to queue if it's time
	select {
	case <-heartbeat_at.C:
		Log.Print("Info: worker heartbeat")
		worker.Send(a.PPP_HEARTBEAT, 0)
	default:
	}

	return nil
}
