package events

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"malgo-websocket/internal/entities"
	"malgo-websocket/internal/log"
	"malgo-websocket/internal/ws"
)

func (h *Handler) SendNewCommandsToWebsocket(ctx context.Context, command *entities.CommandCreated) error {
	log.FromContext(ctx).Info("Sending new command through websockets")

	payload, err := proto.Marshal(command)
	if err != nil {
		return fmt.Errorf("could not marshal command into protobuf: %v", err)
	}

	msg := message.NewMessage(
		uuid.NewString(),
		payload,
	)
	return h.pubSub.Publish(ws.NewCommandsTopic, msg)
}
