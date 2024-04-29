package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
)

func (g *GrpcServer) RegisterNewSession(context.Context, *gateway.RegisterNewSessionRequest) (*gateway.EmptyResponse, error) {
	panic("not implemented")
}
