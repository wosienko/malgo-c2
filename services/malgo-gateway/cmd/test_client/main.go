package main

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, connClose, err := NewGrpcClient()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := connClose(); err != nil {
			panic(err)
		}
	}()

	_, err = client.SessionHeartbeat(context.Background(), &gateway.SessionHeartbeatRequest{
		SessionId: "d946c4b4-77df-4dc8-abbc-3935156f54d6",
	})
	if err != nil {
		panic(err)
	}
}

func NewGrpcClient() (client gateway.GatewayServiceClient, close func() error, err error) {
	grpcAddr := "localhost:8082"

	grpcDialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(grpcAddr, grpcDialOpts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return gateway.NewGatewayServiceClient(conn), conn.Close, nil
}
