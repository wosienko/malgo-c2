package rpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-gateway/internal/entities"
)

type GrpcServer struct {
	gateway.UnimplementedGatewayServiceServer

	commandBus *cqrs.CommandBus

	sessionRepo SessionRepository
	commandRepo CommandRepository
}

func NewGrpcServer(
	commandBus *cqrs.CommandBus,
	sessionRepo SessionRepository,
	commandRepo CommandRepository,
) *GrpcServer {
	if commandBus == nil {
		panic("sessionRepository is nil")
	}
	if sessionRepo == nil {
		panic("sessionRepository is nil")
	}
	if commandRepo == nil {
		panic("commandRepo is nil")
	}

	return &GrpcServer{
		commandBus: commandBus,

		sessionRepo: sessionRepo,
		commandRepo: commandRepo,
	}
}

type SessionRepository interface {
	RegisterNewSession(ctx context.Context, session internalEntities.RegisterNewSession) error
}

type CommandRepository interface {
	GetCommandInfo(ctx context.Context, commandId string) (*internalEntities.CommandInfo, error)
	GetCommandChunk(ctx context.Context, query *internalEntities.CommandChunkQuery) (*internalEntities.CommandChunk, error)
}
