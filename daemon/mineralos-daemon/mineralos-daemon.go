package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/laxamore/mineralos/daemon/utils"
	"github.com/laxamore/mineralos/utils/Log"
	"github.com/laxamore/mineralos/zmq/client"
	zmq "github.com/pebbe/zmq4"
	"github.com/sevlyar/go-daemon"

	"time"
)

const logFileName = "/mineralos/var/log/mineralos.log"

type rigConfig struct {
	SERVER_IP         string
	SERVER_PUBLIC_KEY string
	RIG_ID            string
	RIG_KEY           string
	RIG_PUBLIC_KEY    string
}

func readConf() rigConfig {
	file, err := os.Open("/mineralos/etc/rig.conf")

	defer func() {
		if err = file.Close(); err != nil {
			Log.Panicf("%v", err)
		}
	}()

	b, _ := ioutil.ReadAll(file)
	rigConfString := string(b)

	SERVER_IP := strings.Split(rigConfString, "SERVER_IP=")[1]
	SERVER_IP = strings.Split(SERVER_IP, "\n")[0]

	SERVER_PUBLIC_KEY := strings.Split(rigConfString, "SERVER_PUBLIC_KEY=")[1]
	SERVER_PUBLIC_KEY = strings.Split(SERVER_PUBLIC_KEY, "\n")[0]

	RIG_ID := strings.Split(rigConfString, "RIG_ID=")[1]
	RIG_ID = strings.Split(RIG_ID, "\n")[0]

	RIG_KEY := strings.Split(rigConfString, "RIG_KEY=")[1]
	RIG_KEY = strings.Split(RIG_KEY, "\n")[0]

	RIG_PUBLIC_KEY := strings.Split(rigConfString, "RIG_PUBLIC_KEY=")[1]
	RIG_PUBLIC_KEY = strings.Split(RIG_PUBLIC_KEY, "\n")[0]

	return rigConfig{
		SERVER_IP:         SERVER_IP,
		SERVER_PUBLIC_KEY: SERVER_PUBLIC_KEY,
		RIG_ID:            RIG_ID,
		RIG_KEY:           RIG_KEY,
		RIG_PUBLIC_KEY:    RIG_PUBLIC_KEY,
	}
}

func checkLogSize() int {
	duCommand := exec.Command("bash", "-c", fmt.Sprintf("du -b %s | awk '{printf \"%s\", $1}'", logFileName, "%s"))
	duCommandOutput, _ := duCommand.Output()

	logSize, err := strconv.Atoi(string(duCommandOutput))
	if err != nil {
		log.Fatalf("Unable to get log size: %s", err.Error())
	}

	return logSize
}

func main() {
	daemon.SetSigHandler(termHandler, syscall.SIGTERM)
	daemon.SetSigHandler(termHandler, syscall.SIGQUIT)
	daemon.SetSigHandler(reloadHandler, syscall.SIGHUP)

	cntxt := &daemon.Context{
		PidFileName: "/mineralos/var/tmp/mineralos.pid",
		PidFilePerm: 0644,
		LogFileName: logFileName,
		LogFilePerm: 0640,
		WorkDir:     "/mineralos",
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	lf, err := NewLogFile(logFileName, os.Stderr)
	if err != nil {
		log.Fatalf("Unable to create log file: %s", err.Error())
	}

	log.SetOutput(lf)

	go zmq_client()
	go logRotation(lf)

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}

var (
	stopZMQClient   = make(chan int)
	stopLogRotation = make(chan int)
	doneZMQClient   = make(chan int)
)

func zmq_client() {
	drivers, err := utils.GetGPUDriverVersion()
	if err != nil {
		Log.Print(err)
	}

	zmq.AuthSetVerbose(true)
	zmq.AuthStart()
	defer zmq.AuthStop()

	//  Tell the authenticator to allow any CURVE requests for this domain
	zmq.AuthCurveAdd("*", "*")

	RIG_CONF := readConf()

	cntrl := client.ClientController{
		REQUEST_TIMEOUT: 2500 * time.Millisecond, //  msecs, (> 1000!)
		SERVER_ENDPOINT: fmt.Sprintf("tcp://%s:9000", RIG_CONF.SERVER_IP),

		HEARTBEAT_INTERVAL: 100 * time.Millisecond, //  msecs
		RIG_ID:             RIG_CONF.RIG_ID,
		ClientKey:          RIG_CONF.RIG_KEY,
		ClientPubKey:       RIG_CONF.RIG_PUBLIC_KEY,
		ServerPubKey:       RIG_CONF.SERVER_PUBLIC_KEY,
		DisableLog:         false,
	}

	// Log.Print("Info: connecting to server...\n")
	client, poller, err := cntrl.NewClientConnection(cntrl.ClientPubKey, cntrl.ClientKey)
	if err != nil {
		panic(err)
	}

LOOP:
	for {
		cntrl.PayloadStatus.Drivers = drivers

		lastPayload, _, err := cntrl.Client(client, poller)
		if err != nil {
			// Log.Printf("waring: no response from server retrying...")

			//  Old socket is confused; close it and open a new one
			client.Close()
			client, poller, _ = cntrl.NewClientConnection(cntrl.ClientPubKey, cntrl.ClientKey)

			//  Send request again, on new socket
			client.SendMessage(lastPayload)
		}

		time.Sleep(time.Second)

		select {
		case <-stopZMQClient:
			break LOOP
		default:
		}
	}
	doneZMQClient <- 0
}

func logRotation(lf *LogFile) {
LOOP:
	for {
		if logSize := checkLogSize(); logSize > 5000*1000 {
			if err := lf.Rotate(); err != nil {
				log.Fatalf("Unable to rotate log: %s", err.Error())
			}
		}
		time.Sleep(time.Second)
		select {
		case <-stopLogRotation:
			break LOOP
		default:
		}
	}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")

	stopZMQClient <- 0
	stopLogRotation <- 0

	if sig == syscall.SIGQUIT {
		<-doneZMQClient
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("daemon reloaded")
	return nil
}
