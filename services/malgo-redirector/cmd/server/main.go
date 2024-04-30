package main

import (
	"context"
	"fmt"
	"github.com/VipWW/malgo-c2/services/malgo-redirector/internal/rpc"
	"github.com/VipWW/malgo-c2/services/malgo-redirector/internal/service"
	"os"
	"os/signal"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	panic(err)
	//}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Printf("Starting gRPC client\n")
	grpcClient, closeGrpc, err := rpc.NewGrpcClient("localhost:8082") // TODO: move to env
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Printf("Closing gRPC client\n")
		err := closeGrpc()
		if err != nil {
			panic(err)
		}
	}()

	err = service.New(
		"127.0.0.1:5653", // TODO: move to env
		grpcClient,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
