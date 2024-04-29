package rpc

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
)

func (g *GrpcServer) CommandInfo(context.Context, *gateway.CommandInfoRequest) (*gateway.CommandInfoResponse, error) {
	panic("not implemented")
}
