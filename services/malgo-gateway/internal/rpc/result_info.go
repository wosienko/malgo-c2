package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/google/uuid"
)

func (g *GrpcServer) ResultInfo(ctx context.Context, req *gateway.ResultInfoRequest) (*gateway.EmptyResponse, error) {
	err := uuid.Validate(req.CommandId)
	if err != nil {
		return nil, err
	}

	err = g.resultRepo.SetResultLength(ctx, req.CommandId, int(req.Length))
	if err != nil {
		return nil, err
	}

	return &gateway.EmptyResponse{}, nil
}
