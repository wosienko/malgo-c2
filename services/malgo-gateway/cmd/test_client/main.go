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
	//	ProjectId: "6ed44c13-1a72-4652-936e-d6799a487736",
	//})
	//if err != nil {
	//	fmt.Printf("error: %v\n", err)
	//}
	//fmt.Printf("Registered new session\n")

	fmt.Printf("Getting command info\n")
	resp1, err := client.CommandInfo(context.Background(), &gateway.CommandInfoRequest{
		SessionId: "a264d459-4ba4-453b-8bd6-f40ba39eb087",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Got command info: %v\n", resp1)

	fmt.Printf("Getting command details chunk\n")
	resp2, err := client.CommandDetailsChunk(context.Background(), &gateway.CommandDetailsChunkRequest{
		CommandId: resp1.CommandId,
		Offset:    0,
		Length:    9,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Got command details chunk: %v\n", resp2)

	fmt.Printf("Setting result info\n")
	_, err = client.ResultInfo(context.Background(), &gateway.ResultInfoRequest{
		CommandId: resp1.CommandId,
		Length:    6,
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Set result info\n")

	fmt.Printf("Adding result chunk\n")
	_, err = client.ResultDetailsChunk(context.Background(), &gateway.ResultDetailsChunkRequest{
		CommandId: resp1.CommandId,
		Offset:    0,
		Data:      []byte("whoami"),
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Added result chunk\n")
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
