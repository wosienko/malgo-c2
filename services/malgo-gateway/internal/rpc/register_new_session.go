package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
	"github.com/google/uuid"
)

func (g *GrpcServer) RegisterNewSession(ctx context.Context, req *gateway.RegisterNewSessionRequest) (*gateway.EmptyResponse, error) {
	err := uuid.Validate(req.SessionId)
	if err != nil {
		return nil, err
	}
	err = uuid.Validate(req.ProjectId)
	if err != nil {
		return nil, err
	}

	err = g.sessionRepo.RegisterNewSession(ctx, internalEntities.RegisterNewSession{
		SessionId: req.SessionId,
		ProjectId: req.ProjectId,
	})

	if err != nil {
		return nil, err
	}

	return &gateway.EmptyResponse{}, nil
}
