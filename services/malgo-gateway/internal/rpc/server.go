package rpc

import (
	gateway "github.com/VipWW/malgo-c2/services/common/service"
)

type GrpcServer struct {
	gateway.UnimplementedGatewayServiceServer
}
