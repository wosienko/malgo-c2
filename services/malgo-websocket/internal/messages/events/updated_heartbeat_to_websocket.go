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

func (h *Handler) SendUpdatedHeartbeatToWebsocket(ctx context.Context, session *entities.SessionHeartbeatUpdated) error {
	log.FromContext(ctx).Info("Sending session heartbeat through websockets")

	payload, err := proto.Marshal(session)
	if err != nil {
		return fmt.Errorf("could not marshal session heartbeat into protobuf: %v", err)
	}

	msg := message.NewMessage(
		uuid.NewString(),
		payload,
	)
	return h.pubSub.Publish(ws.SessionHeartbeatUpdatedTopic, msg)
}
