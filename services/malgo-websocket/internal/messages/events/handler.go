package events

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
)

type Handler struct {
	pubSub *gochannel.GoChannel

	commandRepo CommandRepository
	resultRepo  ResultRepository
}

func NewHandler(
	pubSub *gochannel.GoChannel,
	commandRepo CommandRepository,
	resultRepo ResultRepository,
) Handler {
	if pubSub == nil {
		panic("pubSub is required")
	}
	if commandRepo == nil {
		panic("commandRepo is required")
	}
	if resultRepo == nil {
		panic("resultRepo is required")
	}
	return Handler{
		pubSub:      pubSub,
		commandRepo: commandRepo,
		resultRepo:  resultRepo,
	}
}

type CommandRepository interface {
	GetCommandByID(ctx context.Context, id string) (*entities.Command, error)
}

type ResultRepository interface {
	GetResultForCommand(ctx context.Context, commandId string) (string, error)
}
