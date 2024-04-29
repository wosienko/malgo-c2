package main

import (
	"context"
	"fmt"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Printf("Starting test client\n")
	fmt.Printf("---------------------\n")
	fmt.Printf("Connecting to gRPC server\n")
	client, connClose, err := NewGrpcClient()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := connClose(); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("Connected to gRPC server\n")
	fmt.Printf("---------------------\n")

	//fmt.Printf("Sending session heartbeat\n")
	//_, err = client.SessionHeartbeat(context.Background(), &gateway.SessionHeartbeatRequest{
	//	SessionId: "d946c4b4-77df-4dc8-abbc-3935156f54d6",
	//})
	//if err != nil {
	//	fmt.Printf("error: %v\n", err)
	//}
	//fmt.Printf("Sent session heartbeat\n")

	//fmt.Printf("Registering new session\n")
	//_, err = client.RegisterNewSession(context.Background(), &gateway.RegisterNewSessionRequest{
	//	SessionId: uuid.NewString(),
	//	ProjectId: "4e960188-b535-40bf-99ff-567b8144e028",
	//})
	//if err != nil {
	//	fmt.Printf("error: %v\n", err)
	//}
	//fmt.Printf("Registered new session\n")

	fmt.Printf("Getting command info\n")
	resp, err := client.CommandInfo(context.Background(), &gateway.CommandInfoRequest{
		SessionId: "6b7efe93-fc33-4709-969f-4f3ad2d52a49",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Got command info: %v\n", resp)
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
