package rpc

import (
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient(grpcServer string) (client gateway.GatewayServiceClient, close func() error, err error) {
	grpcDialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(grpcServer, grpcDialOpts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return gateway.NewGatewayServiceClient(conn), conn.Close, nil
}
