package main

import (
	"github.com/laxamore/mineralos/utils/Log"
	"github.com/laxamore/mineralos/zmq/client/client"
	zmq "github.com/pebbe/zmq4"

	"time"
)

func main() {
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	defer zmq.AuthStop()

	//  Tell the authenticator to allow any CURVE requests for this domain
	zmq.AuthCurveAdd("*", "*")

	cntrl := client.ClientController{
		REQUEST_TIMEOUT: 2500 * time.Millisecond, //  msecs, (> 1000!)
		SERVER_ENDPOINT: "tcp://127.0.0.1:9000",

		HEARTBEAT_INTERVAL: 1000 * time.Millisecond, //  msecs
		RIG_ID:             "b073badf-c10f-4a94-9d04-e0bb35b26d18",
		ClientKey:          "xOJ)9GNipQOkJWd7k^m:&fF9BSKlu0v73#JCyana",
		ClientPubKey:       "i!OLA?M*RxcyAwZ{#Kvn+ri^F3x-H-lrki=5n6xP",
	}

	Log.Print("Info: connecting to server...\n")
	client, poller, err := cntrl.NewClientConnection(cntrl.ClientPubKey, cntrl.ClientKey)
	if err != nil {
		panic(err)
	}

	for {
		lastPayload, _, err := cntrl.Client(client, poller)
		if err != nil {
			Log.Printf("waring: no response from server retrying...")

			//  Old socket is confused; close it and open a new one
			client.Close()
			client, poller, _ = cntrl.NewClientConnection(cntrl.ClientPubKey, cntrl.ClientKey)

			//  Send request again, on new socket
			client.SendMessage(lastPayload)
		}
	}
}
