package main

import (
	"log"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

func main() {
	cntxt := &daemon.Context{
		PidFileName: "/mineralos/var/tmp/mineralos.pid",
		PidFilePerm: 0644,
		LogFileName: "/mineralos/var/log/mineralos.log",
		LogFilePerm: 0640,
		WorkDir:     "/mineralos",
	}

	d, err := cntxt.Search()

	if err != nil {
		log.Fatalf("Unable send signal to the daemon: %s", err.Error())
	}

	d.Signal(syscall.SIGQUIT)

	daemon.SendCommands(d)
}
