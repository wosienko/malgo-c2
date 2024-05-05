package ws

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 5) / 10
)

type Handler struct {
	conn *websocket.Conn

	userId string

	subscribedSession string
	subscribedProject string

	eventBus   *cqrs.EventBus
	commandBus *cqrs.CommandBus

	userRepo UserRepository

	pubSub *gochannel.GoChannel

	cancel chan struct{}
}

func NewHandler(
	conn *websocket.Conn,
	userId string,
	userRepo UserRepository,
	pubSub *gochannel.GoChannel,
	eventBus *cqrs.EventBus,
	commandBus *cqrs.CommandBus,
) *Handler {
	if conn == nil {
		panic("WS connection is nil")
	}

	if userRepo == nil {
		panic("UserRepo is nil")
	}

	if pubSub == nil {
		panic("PubSub is nil")
	}

	return &Handler{
		conn: conn,

		userId: userId,

		subscribedSession: "",
		subscribedProject: "",

		pubSub: pubSub,

		eventBus:   eventBus,
		commandBus: commandBus,

		userRepo: userRepo,

		cancel: make(chan struct{}),
	}
}

type UserRepository interface {
	GetUserIdIfLoggedInAndOperator(ctx context.Context, sessionId string) (string, error)
	IsUserAllowedToAccessProject(ctx context.Context, userId string, projectId string) (bool, error)
}
