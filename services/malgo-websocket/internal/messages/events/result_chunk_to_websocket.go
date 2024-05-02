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

func (h *Handler) SendResultChunkToWebsocket(ctx context.Context, command *entities.ResultChunkInserted) error {
	log.FromContext(ctx).Info("Sending result chunk info through websockets")

	payload, err := proto.Marshal(command)
	if err != nil {
		return fmt.Errorf("could not marshal result chunk into protobuf: %v", err)
	}

	msg := message.NewMessage(
		uuid.NewString(),
		payload,
	)
	return h.pubSub.Publish(ws.ResultChunkInsertedTopic, msg)
}
