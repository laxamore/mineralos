//
//  Paranoid Pirate queue.
//

package router

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"fmt"
	"time"
)

const (
	HEARTBEAT_LIVENESS = 3                       //  3-5 is reasonable
	HEARTBEAT_INTERVAL = 1000 * time.Millisecond //  msecs

	PPP_READY     = "\001" //  Signals worker is ready
	PPP_HEARTBEAT = "\002" //  Signals worker heartbeat
)

type RouterRepositoryInterface interface {
	FindOne(*mongo.Client, string, string, interface{}) map[string]interface{}
}

type RouterController struct {
	DATABASE_NAME string
}

//  Here we define the worker class; a structure and a set of functions that
//  as constructor, destructor, and methods on worker objects:
type Worker_t struct {
	identity  string    //  Identity of worker
	id_string string    //  Printable identity
	expire    time.Time //  Expires at this time
}

func (a RouterController) Router(clientRtr *zmq.Socket, workerRtr *zmq.Socket, workers []Worker_t, poller1 *zmq.Poller, poller2 *zmq.Poller, client *mongo.Client, repositoryInterface RouterRepositoryInterface) ([]Worker_t, []string, []string, error) {
	var workerMsg []string
	var clientMsg []string

	heartbeat_at := time.NewTicker(HEARTBEAT_INTERVAL)

	var sockets []zmq.Polled
	var err error
	if len(workers) > 0 {
		sockets, err = poller2.Poll(HEARTBEAT_INTERVAL)
	} else {
		sockets, err = poller1.Poll(HEARTBEAT_INTERVAL)
	}

	if err != nil {
		return workers, workerMsg, clientMsg, err //  Interrupted
	}

	for _, socket := range sockets {
		switch socket.Socket {
		case workerRtr:
			//  Handle worker activity on workerRtr
			//  Use worker identity for load-balancing
			msg, err := workerRtr.RecvMessage(0)
			if err != nil {
				break //  Interrupted
			}

			//  Any sign of life from worker means it's ready
			identity, msg := unwrap(msg)
			workerMsg = msg
			workers = s_worker_ready(s_worker_new(identity), workers)

			//  Validate control message, or return reply to client
			if len(msg) == 1 {
				if msg[0] != PPP_READY && msg[0] != PPP_HEARTBEAT {
					fmt.Println("Error: invalid message from worker", msg)
				}
			} else {
				clientRtr.SendMessage(msg)
			}
		case clientRtr:
			//  Now get next client request, route to next worker
			msg, err := clientRtr.RecvMessage(0)
			if err != nil {
				break //  Interrupted
			}
			clientMsg = msg

			type ClientPayload struct {
				RigID  string
				Key    string
				PubKey string
				Status interface{}
			}

			var clientPayload ClientPayload
			json.Unmarshal([]byte(msg[2]), &clientPayload)

			res := repositoryInterface.FindOne(client, a.DATABASE_NAME, "rigs", bson.D{
				{
					Key: "rig_id", Value: clientPayload.RigID,
				},
				{
					Key: "key", Value: clientPayload.Key,
				},
				{
					Key: "pubkey", Value: clientPayload.PubKey,
				},
			})

			if len(res) > 0 {
				workerRtr.SendMessage(workers[0].identity, msg)
				workers = workers[1:]
			}
		}
	}

	//  We handle heartbeating after any socket activity. First we send
	//  heartbeats to any idle workers if it's time. Then we purge any
	//  dead workers:

	select {
	case <-heartbeat_at.C:
		for _, worker := range workers {
			workerRtr.SendMessage(worker.identity, PPP_HEARTBEAT)
		}
	default:
	}

	workers = s_workers_purge(workers)
	return workers, workerMsg, clientMsg, err
}

//  Construct new worker
func s_worker_new(identity string) Worker_t {
	return Worker_t{
		identity:  identity,
		id_string: identity,
		expire:    time.Now().Add(HEARTBEAT_INTERVAL * HEARTBEAT_LIVENESS),
	}
}

//  The ready method puts a worker to the end of the ready list:
func s_worker_ready(self Worker_t, workers []Worker_t) []Worker_t {
	for i, worker := range workers {
		if self.id_string == worker.id_string {
			if i == 0 {
				workers = workers[1:]
			} else if i == len(workers)-1 {
				workers = workers[:i]
			} else {
				workers = append(workers[:i], workers[i+1:]...)
			}
			break
		}
	}
	return append(workers, self)
}

//  The purge method looks for and kills expired workers. We hold workers
//  from oldest to most recent, so we stop at the first alive worker:
func s_workers_purge(workers []Worker_t) []Worker_t {
	now := time.Now()
	for i, worker := range workers {
		if now.Before(worker.expire) {
			return workers[i:] //  Worker is alive, we're done here
		}
	}
	return workers[0:0]
}

//  Pops frame off front of message and returns it as 'head'
//  If next frame is empty, pops that empty frame.
//  Return remaining frames of message as 'tail'
func unwrap(msg []string) (head string, tail []string) {
	head = msg[0]
	if len(msg) > 1 && msg[1] == "" {
		tail = msg[2:]
	} else {
		tail = msg[1:]
	}
	return
}
