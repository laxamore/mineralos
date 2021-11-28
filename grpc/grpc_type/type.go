package grpc_type

import (
	"github.com/laxamore/mineralos/utils/Linux"
)

type ClientPayload struct {
	Drivers Linux.GPUDriverVersion
	RigID   string
}

type Reply struct {
	Message string
}
