package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
	"github.com/google/uuid"
)

func (g *GrpcServer) CommandDetailsChunk(ctx context.Context, req *gateway.CommandDetailsChunkRequest) (*gateway.CommandDetailsChunkResponse, error) {
	err := uuid.Validate(req.CommandId)
	if err != nil {
		return nil, err
	}

	cmdChunk, err := g.commandRepo.GetCommandChunk(ctx, &internalEntities.CommandChunkQuery{
		CommandID: req.CommandId,
		Offset:    int(req.Offset),
		Length:    int(req.Length),
	})
	if err != nil {
		return nil, err
	}

	return &gateway.CommandDetailsChunkResponse{
		Data:        cmdChunk.Data,
		IsLastChunk: cmdChunk.IsLast,
	}, nil
}
