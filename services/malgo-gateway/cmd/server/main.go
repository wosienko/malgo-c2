package main

import (
	"context"
	"fmt"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"os"
	"os/signal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var networkSocket = fmt.Sprintf("127.0.0.1:%d", 8082) // TODO: read from env

	redisClient := messages.NewRedisClient(os.Getenv("REDIS_URL"))
	defer redisClient.Close()

	traceDB, err := otelsql.Open("postgres", os.Getenv("DATABASE_URL"),
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName("db"))
	if err != nil {
		panic(err)
	}

	db := sqlx.NewDb(traceDB, "postgres")
	defer db.Close()

	err = service.New(
		db,
		redisClient,
		networkSocket,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
