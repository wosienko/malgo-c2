package main

import (
	"context"
	"fmt"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/google/uuid"
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

	//_, err = client.SessionHeartbeat(context.Background(), &gateway.SessionHeartbeatRequest{
	//	SessionId: "d946c4b4-77df-4dc8-abbc-3935156f54d6",
	//})
	//if err != nil {
	//	fmt.Printf("error: %v\n", err)
	//}

	_, err = client.RegisterNewSession(context.Background(), &gateway.RegisterNewSessionRequest{
		SessionId: uuid.NewString(),
		ProjectId: "4e960188-b535-40bf-99ff-567b8144e028",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
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
