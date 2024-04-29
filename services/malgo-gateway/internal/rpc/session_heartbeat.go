package rpc

import (
	"context"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/google/uuid"
)

func (g *GrpcServer) SessionHeartbeat(ctx context.Context, req *gateway.SessionHeartbeatRequest) (*gateway.EmptyResponse, error) {
	err := uuid.Validate(req.SessionId)
	if err != nil {
		return nil, fmt.Errorf("invalid session id: %w", err)
	}

	err = g.commandBus.Send(
		ctx,
		&entities.UpdateSessionHeartbeat{
			Header:    entities.NewHeader(),
			SessionId: req.SessionId,
		},
	)
	if err != nil {
		return nil, err
	}

	return &gateway.EmptyResponse{}, nil
}
