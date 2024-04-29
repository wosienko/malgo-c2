package events

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/log"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/ws"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
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
