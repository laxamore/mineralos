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

		HEARTBEAT_INTERVAL: 100 * time.Millisecond, //  msecs
		RIG_ID:             "56b1c9df-867f-4b17-b899-fcca4ab68232",
		ClientKey:          "-GC{6aVX0Bw{ryfY804K!F>gWe{)1#ML@3j=ib[4",
		ClientPubKey:       "IIV?kfy73WCH4+Tf(<N9HxM?Ken[ro(xTTt{C6F@",
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
