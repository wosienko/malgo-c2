package ws

import (
	"encoding/json"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) handleNewSession(input []byte) error {
	var newSession entities.SessionRegistered
	err := proto.Unmarshal(input, &newSession)
	if err != nil {
		return fmt.Errorf("error unmarshalling renamed session: %v", err)
	}

	if newSession.ProjectId != h.subscribedProject {
		return nil
	}

	response := internalEntities.SessionSentToOperator{
		MessageType: "session-new",
		SessionId:   newSession.SessionId,
		Name:        newSession.Name,
		CreatedAt:   newSession.CreatedAt.AsTime().UTC().String(),
		Heartbeat:   newSession.HeartbeatAt.AsTime().UTC().String(),
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshalling response: %v", err)
	}

	if err := h.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return fmt.Errorf("error sending message through websocket: %v", err)
	}
	return nil
}
