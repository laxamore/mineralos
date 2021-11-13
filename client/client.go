//
//  Lazy Pirate client.
//  Use zmq_poll to do a safe request-reply
//  To run, start lpserver and then randomly kill/restart it
//

package main

import (
	"log"
	"runtime"

	zmq "github.com/pebbe/zmq4"

	"fmt"
	"strconv"
	"time"
)

const (
	REQUEST_TIMEOUT = 2500 * time.Millisecond //  msecs, (> 1000!)
	SERVER_ENDPOINT = "tcp://10.8.0.2:9000"
)

func main() {
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	defer zmq.AuthStop()

	//  Tell the authenticator to allow any CURVE requests for this domain
	zmq.AuthCurveAdd("*", zmq.CURVE_ALLOW_ANY)

	client_public, client_secret, err := zmq.NewCurveKeypair()
	checkErr(err)

	fmt.Println("I: connecting to server...")
	client, poller, err := newClientConnection(client_public, client_secret)
	if err != nil {
		panic(err)
	}

	sequence := 0
	for {
		//  We send a request, then we work to get a reply
		sequence++
		client.SendMessage(sequence)

		for expect_reply := true; expect_reply; {
			//  Poll socket for a reply, with timeout
			sockets, err := poller.Poll(REQUEST_TIMEOUT)
			if err != nil {
				break //  Interrupted
			}

			//  Here we process a server reply and exit our loop if the
			//  reply is valid. If we didn't a reply we close the client
			//  socket and resend the request. We try a number of times
			//  before finally abandoning:

			if len(sockets) > 0 {
				//  We got a reply from the server, must match sequence
				reply, err := client.RecvMessage(0)
				if err != nil {
					break //  Interrupted
				}
				seq, _ := strconv.Atoi(reply[0])
				if seq == sequence {
					fmt.Printf("I: server replied OK (%s)\n", reply[0])
					expect_reply = false
				} else {
					fmt.Printf("E: malformed reply from server: %s\n", reply)
				}
			} else {
				fmt.Println("W: no response from server, retrying...")
				//  Old socket is confused; close it and open a new one
				client.Close()
				client, poller, _ = newClientConnection(client_public, client_secret)

				//  Send request again, on new socket
				client.SendMessage(sequence)
			}
		}
	}
}

func newClientConnection(client_public string, client_secret string) (soc *zmq.Socket, pol *zmq.Poller, err error) {
	soc, err = zmq.NewSocket(zmq.REQ)
	soc.ClientAuthCurve("83<s>=wXS9RXKPR4wp:45?Pmo!y>R!qAy%^:^dDl", client_public, client_secret)
	soc.Connect(SERVER_ENDPOINT)

	// Recreate poller for new client
	pol = zmq.NewPoller()
	pol.Add(soc, zmq.POLLIN)

	return soc, pol, err
}

func checkErr(err error) {
	if err != nil {
		log.SetFlags(0)
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			log.Fatalf("%v:%v: %v", filename, lineno, err)
		} else {
			log.Fatalln(err)
		}
	}
}
