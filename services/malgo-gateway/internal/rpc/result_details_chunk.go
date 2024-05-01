package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
)

func (g *GrpcServer) ResultDetailsChunk(ctx context.Context, req *gateway.ResultDetailsChunkRequest) (*gateway.EmptyResponse, error) {
	err := g.resultRepo.AddResultChunk(ctx, internalEntities.ResultChunk{
		CommandId: req.CommandId,
		Offset:    int(req.Offset),
		Chunk:     req.Data,
	})
	if err != nil {
		return nil, err
	}

	return &gateway.EmptyResponse{}, nil
}
