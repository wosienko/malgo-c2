package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	log2 "github.com/VipWW/malgo-c2/services/common/log"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/db"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/commands"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/observability"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/rpc"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
)

func init() {
	log2.Init(logrus.InfoLevel)
}

type Service struct {
	db *sqlx.DB

	watermillRouter *message.Router

	grpcServer *grpc.Server
	grpcAddr   string

	traceProvider *tracesdk.TracerProvider
}

func New(
	dbConn *sqlx.DB,
	redisClient *redis.Client,
	httpNetworkAddress string,
) Service {
	traceProvider := observability.ConfigureTraceProvider()

	watermillLogger := log2.NewWatermill(log2.FromContext(context.Background()))

	var redisPublisher message.Publisher
	redisPublisher = messages.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log2.CorrelationPublisherDecorator{Publisher: redisPublisher}
	redisPublisher = observability.TracingPublisherDecorator{Publisher: redisPublisher}

	//redisSubscriber := messages.NewRedisSubscriber(redisClient, watermillLogger)

	commandRepo := db.NewCommandRepository(dbConn)
	sessionRepo := db.NewSessionRepository(dbConn)

	commandBus := commands.NewBusWithConfig(redisPublisher, commands.NewBusConfig(watermillLogger))

	commandHandler := commands.NewHandler(
		commandRepo,
		sessionRepo,
	)

	commandProcessorConfig := commands.NewProcessorConfig(redisClient, watermillLogger)

	watermillRouter := messages.NewWatermillRouter(
		commandProcessorConfig,
		commandHandler,
		watermillLogger,
	)

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	gateway.RegisterGatewayServiceServer(s, rpc.NewGrpcServer(
		commandBus,
		sessionRepo,
	))

	return Service{
		db:              dbConn,
		watermillRouter: watermillRouter,
		grpcServer:      s,
		grpcAddr:        httpNetworkAddress,
		traceProvider:   traceProvider,
	}
}

func (s Service) Run(
	ctx context.Context,
) error {
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		fmt.Println("Starting Watermill router")
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		fmt.Printf("Starting GRPC server on %s\n", s.grpcAddr)
		lis, err := net.Listen("tcp", s.grpcAddr)
		if err != nil {
			return err
		}
		err = s.grpcServer.Serve(lis)

		if !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		fmt.Printf("Shutting down trace provider\n")
		return s.traceProvider.Shutdown(context.Background())
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		fmt.Printf("Shutting down GRPC server\n")
		s.grpcServer.GracefulStop()
		return nil
	})
	errgrp.Go(func() error {
		<-ctx.Done()
		fmt.Printf("Shutting down Watermill router\n")
		return s.watermillRouter.Close()
	})

	return errgrp.Wait()
}
