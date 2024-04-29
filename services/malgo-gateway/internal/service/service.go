package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	log2 "github.com/VipWW/malgo-c2/services/common/log"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/db"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/commands"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/events"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/observability"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func init() {
	log2.Init(logrus.InfoLevel)
}

type Service struct {
	db *sqlx.DB

	watermillRouter *message.Router

	httpServer *http.Server

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

	//eventBus := events.NewBusWithConfig(redisPublisher, events.NewBusConfig(watermillLogger))
	//commandBus := commands.NewBusWithConfig(redisPublisher, commands.NewBusConfig(watermillLogger))

	eventHandler := events.NewHandler()
	commandHandler := commands.NewHandler(
		commandRepo,
		sessionRepo,
	)

	eventProcessorConfig := events.NewProcessorConfig(redisClient, watermillLogger)
	commandProcessorConfig := commands.NewProcessorConfig(redisClient, watermillLogger)

	watermillRouter := messages.NewWatermillRouter(
		eventProcessorConfig,
		eventHandler,
		commandProcessorConfig,
		commandHandler,
		watermillLogger,
	)

	return Service{
		db:              dbConn,
		watermillRouter: watermillRouter,
		httpServer: &http.Server{
			Addr: httpNetworkAddress,
		},
		traceProvider: traceProvider,
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

		fmt.Printf("Starting HTTP server on %s\n", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
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
		fmt.Printf("Shutting down HTTP server\n")
		return s.httpServer.Shutdown(context.Background())
	})
	errgrp.Go(func() error {
		<-ctx.Done()
		fmt.Printf("Shutting down Watermill router\n")
		return s.watermillRouter.Close()
	})

	return errgrp.Wait()
}
