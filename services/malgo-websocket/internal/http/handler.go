package http

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Handler struct {
	eventBus   *cqrs.EventBus
	commandBus *cqrs.CommandBus

	pubSub *gochannel.GoChannel

	userRepo UserRepo
}

type UserRepo interface {
	GetUserIdIfLoggedInAndOperator(ctx context.Context, sessionId string) (string, error)
}
