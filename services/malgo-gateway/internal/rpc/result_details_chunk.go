package rpc

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/log"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
)

func (g *GrpcServer) ResultDetailsChunk(ctx context.Context, req *gateway.ResultDetailsChunkRequest) (*gateway.EmptyResponse, error) {
	log.FromContext(ctx).Infof("Adding result chunk %v", req.Offset)

	err := g.resultRepo.AddResultChunk(ctx, internalEntities.ResultChunk{
		CommandId: req.CommandId,
		Offset:    int(req.Offset),
		Chunk:     req.Data,
	})
	if err != nil {
		return nil, err
	}
	log.FromContext(ctx).Infof("Added result chunk %v", req.Offset)

	return &gateway.EmptyResponse{}, nil
}
