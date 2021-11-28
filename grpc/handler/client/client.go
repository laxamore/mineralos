package client

import (
	"context"

	"github.com/laxamore/mineralos/grpc/grpc_type"
	pb "github.com/laxamore/mineralos/grpc/mineralos_proto"
	"github.com/laxamore/mineralos/utils/Log"
)

type ClientController struct{}

func (a ClientController) TryClient(c pb.MineralosClient, ctx context.Context, payload grpc_type.ClientPayload) {
	r, err := c.ReportStatus(ctx, &pb.Payload{
		RigId:        payload.RigID,
		AmdDriver:    payload.Drivers.AMD,
		NvidiaDriver: payload.Drivers.NVIDIA,
	})

	if err != nil {
		Log.Print("could not report status: cannot connect to server")
	}
	Log.Printf("Response: %s", r.GetMessage())
}
