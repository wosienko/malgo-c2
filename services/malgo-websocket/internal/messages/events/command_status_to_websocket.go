package events

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/ws"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) SendCommandStatusToWebsocket(ctx context.Context, keyValue *entities.CommandStatusModified) error {
	log.FromContext(ctx).Info("Sending command status through websockets")

	payload, err := proto.Marshal(keyValue)
	if err != nil {
		return fmt.Errorf("could not marshal command status into protobuf: %v", err)
	}

	msg := message.NewMessage(
		uuid.NewString(),
		payload,
	)
	return h.pubSub.Publish(ws.CommandStatusUpdatedTopic, msg)
}
