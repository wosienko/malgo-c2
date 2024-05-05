package ws

import (
	"encoding/json"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) handleCommandStatusChanged(input []byte) error {
	var command entities.CommandStatusModified
	err := proto.Unmarshal(input, &command)
	if err != nil {
		return fmt.Errorf("error unmarshalling key-value: %v", err)
	}

	if command.SessionId != h.subscribedSession {
		return nil
	}

	response := internalEntities.CommandStatusSentToOperator{
		MessageType: "command-status-updated",
		ID:          command.CommandId,
		Status:      command.Status,
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
