package client

import (
	"context"

	pb "github.com/laxamore/mineralos/grpc/mineralos_proto"
	"github.com/laxamore/mineralos/utils/Log"
)

type ClientController struct{}

func (a ClientController) TryClient(c pb.MineralosClient, ctx context.Context, payload *pb.Payload) {
	r, err := c.ReportStatus(ctx, payload)
	if err != nil {
		Log.Print("could not report status: cannot connect to server")
	}
	Log.Printf("Response: %s", r.GetMessage())
}
