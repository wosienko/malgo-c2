package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
)

func (g *GrpcServer) SessionHeartbeat(context.Context, *gateway.SessionHeartbeatRequest) (*gateway.EmptyResponse, error) {
	panic("not implemented")
}
