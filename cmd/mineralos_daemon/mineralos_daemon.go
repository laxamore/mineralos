package main

import (
	"context"
	"fmt"
	"github.com/laxamore/mineralos/cmd/mineralos_daemon/handler/client"
	"github.com/laxamore/mineralos/internal/embeding"
	"github.com/laxamore/mineralos/internal/getgpu"
	"github.com/laxamore/mineralos/internal/logger"
	"github.com/laxamore/mineralos/pkg/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	pb "github.com/laxamore/mineralos/config/mineralos_proto"
	"github.com/sevlyar/go-daemon"
	"google.golang.org/grpc"

	"time"
)

const logFileName = "/mineralos/var/log/mineralos.log"

type rigConfig struct {
	SERVER_IP string
	RIG_ID    string
}

func readConf() rigConfig {
	file, err := os.Open("/mineralos/etc/rigs.conf")

	defer func() {
		if err = file.Close(); err != nil {
			logger.Panicf("%v", err)
		}
	}()

	b, _ := ioutil.ReadAll(file)
	rigConfString := string(b)

	SERVER_IP := strings.Split(rigConfString, "SERVER_IP=")[1]
	SERVER_IP = strings.Split(SERVER_IP, "\n")[0]

	RIG_ID := strings.Split(rigConfString, "RIG_ID=")[1]
	RIG_ID = strings.Split(RIG_ID, "\n")[0]

	return rigConfig{
		SERVER_IP: SERVER_IP,
		RIG_ID:    RIG_ID,
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

	go grpc_client()
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

func grpc_client() {
	drivers, err := getgpu.GetGPUDriverVersion()
	if err != nil {
		logger.Print(err)
	}

	gpus, err := getgpu.GetGPU()
	if err != nil {
		logger.Print(err)
	}
	pbGpus := getgpu.ArrGPUSToPBGPUS(gpus)

	RIG_CONF := readConf()
	creds, err := embeding.LoadClientTLSCert(false, "/mineralos/etc/ssl/certs/config-public-cert.pem")
	utils.CheckErr(err)

	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:9000", RIG_CONF.SERVER_IP), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMineralosClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cntrl := client.ClientController{}

LOOP:
	for {
		cntrl.TryClient(c, ctx, &pb.Payload{
			RigId: RIG_CONF.RIG_ID,
			Status: &pb.Status{
				Drivers: &pb.Drivers{
					AMD:    drivers.AMD,
					NVIDIA: drivers.NVIDIA,
				},
				Gpus: pbGpus,
			},
		})
		time.Sleep(time.Millisecond * 100)

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
