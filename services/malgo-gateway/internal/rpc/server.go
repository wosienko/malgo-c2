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
	resultRepo  ResultRepository
}

func NewGrpcServer(
	commandBus *cqrs.CommandBus,
	sessionRepo SessionRepository,
	commandRepo CommandRepository,
	resultRepo ResultRepository,
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
	if resultRepo == nil {
		panic("resultRepo is nil")
	}

	return &GrpcServer{
		commandBus: commandBus,

		sessionRepo: sessionRepo,
		commandRepo: commandRepo,
		resultRepo:  resultRepo,
	}
}

type SessionRepository interface {
	RegisterNewSession(ctx context.Context, session internalEntities.RegisterNewSession) error
}

type CommandRepository interface {
	GetCommandInfo(ctx context.Context, commandId string) (*internalEntities.CommandInfo, error)
	GetCommandChunk(ctx context.Context, query *internalEntities.CommandChunkQuery) (*internalEntities.CommandChunk, error)
}

type ResultRepository interface {
	SetResultLength(ctx context.Context, commandId string, length int) error
	AddResultChunk(ctx context.Context, chunk internalEntities.ResultChunk) error
}
