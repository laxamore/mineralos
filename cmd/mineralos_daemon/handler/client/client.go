package client

import (
	"context"
	pb "github.com/laxamore/mineralos/config/mineralos_proto"
	"github.com/laxamore/mineralos/internal/logger"
)

type ClientController struct{}

func (a ClientController) TryClient(c pb.MineralosClient, ctx context.Context, payload *pb.Payload) {
	r, err := c.ReportStatus(ctx, payload)
	if err != nil {
		logger.Printf("could not report status: cannot connect to server %v", err)
	}
	logger.Printf("Response: %s", r.GetMessage())
}
