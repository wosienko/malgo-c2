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
}

func NewGrpcServer(
	commandBus *cqrs.CommandBus,
	sessionRepo SessionRepository,
) *GrpcServer {
	if commandBus == nil {
		panic("sessionRepository is nil")
	}
	if sessionRepo == nil {
		panic("sessionRepository is nil")
	}

	return &GrpcServer{
		commandBus: commandBus,

		sessionRepo: sessionRepo,
	}
}

type SessionRepository interface {
	RegisterNewSession(ctx context.Context, session internalEntities.RegisterNewSession) error
}
