package rpc

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
)

type GrpcServer struct {
	gateway.UnimplementedGatewayServiceServer

	commandBus *cqrs.CommandBus
}

func NewGrpcServer(
	commandBus *cqrs.CommandBus,
) *GrpcServer {
	if commandBus == nil {
		panic("sessionRepository is nil")
	}

	return &GrpcServer{
		commandBus: commandBus,
	}
}
