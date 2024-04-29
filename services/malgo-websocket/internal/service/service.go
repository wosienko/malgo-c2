package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/db"
	httpSrv "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/http"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/log"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/messages"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/messages/commands"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/messages/events"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/messages/outbox"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/observability"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func init() {
	log.Init(logrus.InfoLevel)
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

	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	var redisPublisher message.Publisher
	redisPublisher = messages.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}
	redisPublisher = observability.TracingPublisherDecorator{Publisher: redisPublisher}

	//redisSubscriber := messages.NewRedisSubscriber(redisClient, watermillLogger)

	postgresSubscriber := outbox.NewPostgresSubscriber(dbConn.DB, watermillLogger)

	userRepo := db.NewUserRepository(dbConn)
	commandRepo := db.NewCommandRepository(dbConn)
	sessionRepo := db.NewSessionRepository(dbConn)

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, watermillLogger) // TODO: check if buffer size is ok

	eventBus := events.NewBusWithConfig(redisPublisher, events.NewBusConfig(watermillLogger))
	commandBus := commands.NewBusWithConfig(redisPublisher, commands.NewBusConfig(watermillLogger))

	eventHandler := events.NewHandler(
		pubSub,
	)
	commandHandler := commands.NewHandler(
		commandRepo,
		sessionRepo,
	)

	eventProcessorConfig := events.NewProcessorConfig(redisClient, watermillLogger)
	commandProcessorConfig := commands.NewProcessorConfig(redisClient, watermillLogger)

	watermillRouter := messages.NewWatermillRouter(
		postgresSubscriber,
		redisPublisher,
		eventProcessorConfig,
		eventHandler,
		commandProcessorConfig,
		commandHandler,
		watermillLogger,
	)

	router := httpSrv.NewRouter(
		eventBus,
		commandBus,
		userRepo,
		pubSub,
	)

	return Service{
		db:              dbConn,
		watermillRouter: watermillRouter,
		httpServer: &http.Server{
			Addr:    httpNetworkAddress,
			Handler: router,
		},
		traceProvider: traceProvider,
	}
}

func (s Service) Run(
	ctx context.Context,
) error {
	err := outbox.InitializeSchema(s.db.DB)
	if err != nil {
		return fmt.Errorf("failed to initialize outbox schema: %w", err)
	}

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
