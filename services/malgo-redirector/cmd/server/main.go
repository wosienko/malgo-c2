package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/VipWW/malgo-c2/services/malgo-redirector/internal/rpc"
	"github.com/VipWW/malgo-c2/services/malgo-redirector/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.FromContext(context.Background()).Error(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Printf("Starting gRPC client\n")
	grpcClient, closeGrpc, err := rpc.NewGrpcClient(os.Getenv("GRPC_ADDR"))
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
		os.Getenv("DNS_ADDR"),
		grpcClient,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
