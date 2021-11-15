package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/laxamore/mineralos/utils/Log"
	zmq "github.com/pebbe/zmq4"
)

type ClientController struct {
	REQUEST_TIMEOUT time.Duration //  msecs, (> 1000!)
	SERVER_ENDPOINT string

	HEARTBEAT_INTERVAL time.Duration //  msecs
	RIG_ID             string
	ClientKey          string
	ClientPubKey       string
}

type Payload struct {
	RigID  string
	Key    string
	PubKey string
	Status interface{}
}

func (a ClientController) Client(client *zmq.Socket, poller *zmq.Poller) ([]byte, []string, error) {
	var workerResponse []string
	//  We send a request, then we work to get a reply
	payload := Payload{
		RigID:  a.RIG_ID,
		Key:    a.ClientKey,
		PubKey: a.ClientPubKey,
		Status: []map[string]interface{}{
			{
				"gpu":      "RTX3060",
				"hashrate": "33mhs",
				"miner":    "t-rex",
				"algo":     "ethash",
			},
			{
				"gpu":      "RTX3060",
				"hashrate": "32mhs",
				"miner":    "t-rex",
				"algo":     "ethash",
			},
		},
	}

	payloadByte, err := json.Marshal(payload)

	if err != nil {
		Log.Panicf("error: payload json marshal %v", err)
	}

	client.SendMessage(payloadByte)

	for expect_reply := true; expect_reply; {
		//  Poll socket for a reply, with timeout
		sockets, err := poller.Poll(a.REQUEST_TIMEOUT)
		if err != nil {
			Log.Printf("%v", err)
			break //  Interrupted
		}

		//  Here we process a server reply and exit our loop if the
		//  reply is valid. If we didn't a reply we close the client
		//  socket and resend the request. We try a number of times
		//  before finally abandoning:

		if len(sockets) > 0 {
			type ServerReply struct {
				Status string
				Config interface{}
			}

			//  We got a reply from the server, must match sequence
			reply, err := client.RecvMessage(0)
			workerResponse = reply

			if err != nil {
				break //  Interrupted
			}

			var replyMsg ServerReply
			json.Unmarshal([]byte(reply[0]), &replyMsg)

			// seq, _ := strconv.Atoi(reply[0])
			if replyMsg.Status == "ok" {
				Log.Printf("info: server replied OK (%s)\n", replyMsg.Config)
				expect_reply = false
				time.Sleep(a.HEARTBEAT_INTERVAL)
			} else {
				Log.Printf("error: malformed reply from server: %s\n", reply)
			}
		} else {
			return payloadByte, workerResponse, fmt.Errorf("error: no response from server")
		}
	}

	return payloadByte, workerResponse, nil
}

func (a ClientController) NewClientConnection(client_public string, client_secret string) (soc *zmq.Socket, pol *zmq.Poller, err error) {
	soc, err = zmq.NewSocket(zmq.REQ)
	soc.ClientAuthCurve("83<s>=wXS9RXKPR4wp:45?Pmo!y>R!qAy%^:^dDl", client_public, client_secret)
	soc.Connect(a.SERVER_ENDPOINT)

	// Recreate poller for new client
	pol = zmq.NewPoller()
	pol.Add(soc, zmq.POLLIN)

	return soc, pol, err
}
