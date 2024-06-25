package httpproxy

import (
	gateway "github.com/VipWW/malgo-c2/services/common/service"
)

const maxHTTPMessageSize = 2048

type Handler struct {
	grpcClient gateway.GatewayServiceClient
}

func NewHandler(
	grpcClient gateway.GatewayServiceClient,
) *Handler {
	return &Handler{
		grpcClient: grpcClient,
	}
}
