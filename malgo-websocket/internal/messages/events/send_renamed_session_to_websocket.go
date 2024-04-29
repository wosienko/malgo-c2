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

func (h *Handler) SendRenamedSessionToWebsocket(ctx context.Context, command *entities.SessionNameModified) error {
	log.FromContext(ctx).Info("Sending renamed session through websockets")

	payload, err := proto.Marshal(command)
	if err != nil {
		return fmt.Errorf("could not marshal session name into protobuf: %v", err)
	}

	msg := message.NewMessage(
		uuid.NewString(),
		payload,
	)
	return h.pubSub.Publish(ws.RenamedSessionTopic, msg)
}
