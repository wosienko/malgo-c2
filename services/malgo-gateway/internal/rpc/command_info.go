package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/google/uuid"
)

func (g *GrpcServer) CommandInfo(ctx context.Context, req *gateway.CommandInfoRequest) (*gateway.CommandInfoResponse, error) {
	err := uuid.Validate(req.SessionId)
	if err != nil {
		return nil, err
	}

	info, err := g.commandRepo.GetCommandInfo(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	go func() {
		_, _ = g.SessionHeartbeat(ctx, &gateway.SessionHeartbeatRequest{
			SessionId: req.SessionId,
		})
	}()

	return &gateway.CommandInfoResponse{
		CommandId:     info.ID,
		Type:          info.Type,
		CommandLength: info.Length,
	}, nil
}
