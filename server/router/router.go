// //
// //  Simple Pirate broker.
// //  This is identical to load-balancing pattern, with no reliability
// //  mechanisms. It depends on the client for recovery. Runs forever.
// //

// package main

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// 	"github.com/laxamore/mineralos/db"
// 	zmq "github.com/pebbe/zmq4"
// )

// type ServerRepositoryInterface interface {
// 	FindOne(string, string, interface{}) map[string]interface{}
// }

// type ServerController struct{}

// const (
// 	WORKER_READY = "\001" //  Signals worker is ready
// )

// func (a ServerController) Server(socRtrClient *zmq.Socket, socRtrWorker *zmq.Socket,
// 	repositoryInterface ServerRepositoryInterface) {

// 	//  Queue of available workers
// 	workers := make([]string, 0)

// 	poller1 := zmq.NewPoller()
// 	poller1.Add(socRtrWorker, zmq.POLLIN)
// 	poller2 := zmq.NewPoller()
// 	poller2.Add(socRtrWorker, zmq.POLLIN)
// 	poller2.Add(socRtrClient, zmq.POLLIN)

// 	//  The body of this example is exactly the same as lbbroker2.
// LOOP:
// 	for {
// 		//  Poll frontend only if we have available workers
// 		var sockets []zmq.Polled
// 		var err error
// 		if len(workers) > 0 {
// 			sockets, err = poller2.Poll(-1)
// 		} else {
// 			sockets, err = poller1.Poll(-1)
// 		}
// 		if err != nil {
// 			break //  Interrupted
// 		}
// 		for _, socket := range sockets {
// 			switch s := socket.Socket; s {
// 			case socRtrWorker: //  Handle worker activity on backend
// 				//  Use worker identity for load-balancing
// 				msg, err := s.RecvMessage(0)
// 				if err != nil {
// 					break LOOP //  Interrupted
// 				}
// 				var identity string
// 				identity, msg = unwrap(msg)
// 				workers = append(workers, identity)

// 				//  Forward message to client if it's not a READY
// 				if msg[0] != WORKER_READY {
// 					socRtrClient.SendMessage(msg)
// 				}

// 			case socRtrClient:
// 				//  Get client request, route to first available worker
// 				msg, err := s.RecvMessage(0)
// 				if err == nil {
// 					socRtrWorker.SendMessage(workers[0], "", msg)
// 					workers = workers[1:]
// 				}
// 			}
// 		}
// 	}
// }

// //  Pops frame off front of message and returns it as 'head'
// //  If next frame is empty, pops that empty frame.
// //  Return remaining frames of message as 'tail'
// func unwrap(msg []string) (head string, tail []string) {
// 	head = msg[0]
// 	if len(msg) > 1 && msg[1] == "" {
// 		tail = msg[2:]
// 	} else {
// 		tail = msg[1:]
// 	}
// 	return
// }

// func main() {
// 	// load .env file
// 	err := godotenv.Load()

// 	if err != nil {
// 		log.Panicf("Error loading .env file")
// 	}

// 	//  Start authentication engine
// 	zmq.AuthSetVerbose(true)
// 	zmq.AuthStart()
// 	defer zmq.AuthStop()
// 	zmq.AuthAllow("*", "0.0.0.0/0")

// 	//  Tell the authenticator to allow any CURVE requests for this domain
// 	zmq.AuthCurveAdd("*", zmq.CURVE_ALLOW_ANY)
// 	server_secret := os.Getenv("SERVER_SECRET")

// 	clientRtr, _ := zmq.NewSocket(zmq.ROUTER)
// 	clientRtr.SetRouterHandover(true)
// 	clientRtr.ServerAuthCurve("*", server_secret)
// 	defer clientRtr.Close()

// 	workerRtr, _ := zmq.NewSocket(zmq.ROUTER)
// 	defer workerRtr.Close()

// 	clientRtr.Bind("tcp://*:9000") //  For clients
// 	workerRtr.Bind("tcp://*:9001") //  For workers

// 	repo := db.MongoDB{}
// 	cntrl := ServerController{}

// 	cntrl.Server(clientRtr, workerRtr, repo)
// }

//
//  Paranoid Pirate queue.
//

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"time"
)

const (
	HEARTBEAT_LIVENESS = 3                       //  3-5 is reasonable
	HEARTBEAT_INTERVAL = 1000 * time.Millisecond //  msecs

	PPP_READY     = "\001" //  Signals worker is ready
	PPP_HEARTBEAT = "\002" //  Signals worker heartbeat
)

//  Here we define the worker class; a structure and a set of functions that
//  as constructor, destructor, and methods on worker objects:

type worker_t struct {
	identity  string    //  Identity of worker
	id_string string    //  Printable identity
	expire    time.Time //  Expires at this time
}

//  Construct new worker
func s_worker_new(identity string) worker_t {
	return worker_t{
		identity:  identity,
		id_string: identity,
		expire:    time.Now().Add(HEARTBEAT_INTERVAL * HEARTBEAT_LIVENESS),
	}
}

//  The ready method puts a worker to the end of the ready list:

func s_worker_ready(self worker_t, workers []worker_t) []worker_t {
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

func s_workers_purge(workers []worker_t) []worker_t {
	now := time.Now()
	for i, worker := range workers {
		if now.Before(worker.expire) {
			return workers[i:] //  Worker is alive, we're done here
		}
	}
	return workers[0:0]
}

//  The main task is a load-balancer with heartbeating on workers so we
//  can detect crashed or blocked worker tasks:

func main() {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Panicf("Error loading .env file")
	}

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
	defer workerRtr.Close()

	clientRtr.Bind("tcp://*:9000") //  For clients
	workerRtr.Bind("tcp://*:9001") //  For workers

	//  List of available workers
	workers := make([]worker_t, 0)

	//  Send out heartbeats at regular intervals
	heartbeat_at := time.Tick(HEARTBEAT_INTERVAL)

	poller1 := zmq.NewPoller()
	poller1.Add(workerRtr, zmq.POLLIN)
	poller2 := zmq.NewPoller()
	poller2.Add(workerRtr, zmq.POLLIN)
	poller2.Add(clientRtr, zmq.POLLIN)

	for {
		//  Poll clientRtr only if we have available workers
		var sockets []zmq.Polled
		var err error
		if len(workers) > 0 {
			sockets, err = poller2.Poll(HEARTBEAT_INTERVAL)
		} else {
			sockets, err = poller1.Poll(HEARTBEAT_INTERVAL)
		}
		if err != nil {
			break //  Interrupted
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
				workers = s_worker_ready(s_worker_new(identity), workers)

				//  Validate control message, or return reply to client
				if len(msg) == 1 {
					if msg[0] != PPP_READY && msg[0] != PPP_HEARTBEAT {
						fmt.Println("E: invalid message from worker", msg)
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
				workerRtr.SendMessage(workers[0].identity, msg)
				workers = workers[1:]
			}
		}

		//  We handle heartbeating after any socket activity. First we send
		//  heartbeats to any idle workers if it's time. Then we purge any
		//  dead workers:

		select {
		case <-heartbeat_at:
			for _, worker := range workers {
				workerRtr.SendMessage(worker.identity, PPP_HEARTBEAT)
			}
		default:
		}
		workers = s_workers_purge(workers)
	}
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
