package ws

import (
	"encoding/json"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) handleUpdatedHeartbeat(input []byte) error {
	var updatedHeartbeat entities.SessionHeartbeatUpdated
	err := proto.Unmarshal(input, &updatedHeartbeat)
	if err != nil {
		return fmt.Errorf("error unmarshalling renamed session: %v", err)
	}

	if updatedHeartbeat.ProjectId != h.subscribedProject {
		return nil
	}

	response := internalEntities.SessionHeartbeatSentToOperator{
		MessageType: "session-heartbeat",
		SessionId:   updatedHeartbeat.SessionId,
		Heartbeat:   updatedHeartbeat.HeartbeatAt.AsTime().UTC().String(),
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshalling response: %v", err)
	}

	fmt.Printf("Sending message through websocket: %v\n", response)

	if err := h.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return fmt.Errorf("error sending message through websocket: %v", err)
	}
	return nil
}
