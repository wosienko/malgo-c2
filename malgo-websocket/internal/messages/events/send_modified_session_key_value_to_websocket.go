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

func (h *Handler) SendModifiedSessionKeyValueToWebsocket(ctx context.Context, keyValue *entities.SessionKeyValueModified) error {
	log.FromContext(ctx).Info("Sending modified key-value through websockets")

	payload, err := proto.Marshal(keyValue)
	if err != nil {
		return fmt.Errorf("could not marshal key-value into protobuf: %v", err)
	}

	msg := message.NewMessage(
		uuid.NewString(),
		payload,
	)
	return h.pubSub.Publish(ws.SessionKeyValueModifiedTopic, msg)
}
