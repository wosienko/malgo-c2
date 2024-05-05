package http

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/ws"
)

type Handler struct {
	eventBus   *cqrs.EventBus
	commandBus *cqrs.CommandBus

	pubSub *gochannel.GoChannel

	userRepo ws.UserRepository
}
