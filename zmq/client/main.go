package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/laxamore/mineralos/utils/Log"
	"github.com/laxamore/mineralos/zmq/client/client"
	zmq "github.com/pebbe/zmq4"

	"time"
)

type rigConfig struct {
	SERVER_PUBLIC_KEY string
	RIG_ID            string
	RIG_KEY           string
	RIG_PUBLIC_KEY    string
}

func readConf() rigConfig {
	file, err := os.Open(os.Args[1])

	defer func() {
		if err = file.Close(); err != nil {
			Log.Panicf("%v", err)
		}
	}()

	b, _ := ioutil.ReadAll(file)
	rigConfString := string(b)

	SERVER_PUBLIC_KEY := strings.Split(rigConfString, "SERVER_PUBLIC_KEY=")[1]
	SERVER_PUBLIC_KEY = strings.Split(SERVER_PUBLIC_KEY, "\n")[0]

	RIG_ID := strings.Split(rigConfString, "RIG_ID=")[1]
	RIG_ID = strings.Split(RIG_ID, "\n")[0]

	RIG_KEY := strings.Split(rigConfString, "RIG_KEY=")[1]
	RIG_KEY = strings.Split(RIG_KEY, "\n")[0]

	RIG_PUBLIC_KEY := strings.Split(rigConfString, "RIG_PUBLIC_KEY=")[1]
	RIG_PUBLIC_KEY = strings.Split(RIG_PUBLIC_KEY, "\n")[0]

	return rigConfig{
		SERVER_PUBLIC_KEY: SERVER_PUBLIC_KEY,
		RIG_ID:            RIG_ID,
		RIG_KEY:           RIG_KEY,
		RIG_PUBLIC_KEY:    RIG_PUBLIC_KEY,
	}
}

func main() {
	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	defer zmq.AuthStop()

	//  Tell the authenticator to allow any CURVE requests for this domain
	zmq.AuthCurveAdd("*", "*")

	RIG_CONF := readConf()

	cntrl := client.ClientController{
		REQUEST_TIMEOUT: 2500 * time.Millisecond, //  msecs, (> 1000!)
		SERVER_ENDPOINT: fmt.Sprintf("tcp://%s:9000", os.Args[2]),

		HEARTBEAT_INTERVAL: 100 * time.Millisecond, //  msecs
		RIG_ID:             RIG_CONF.RIG_ID,
		ClientKey:          RIG_CONF.RIG_KEY,
		ClientPubKey:       RIG_CONF.RIG_PUBLIC_KEY,
		ServerPubKey:       RIG_CONF.SERVER_PUBLIC_KEY,
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
