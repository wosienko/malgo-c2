package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
)

func (g *GrpcServer) CommandInfo(ctx context.Context, req *gateway.CommandInfoRequest) (*gateway.CommandInfoResponse, error) {
	info, err := g.commandRepo.GetCommandInfo(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &gateway.CommandInfoResponse{
		CommandId:     info.ID,
		Type:          info.Type,
		CommandLength: info.Length,
	}, nil
}
