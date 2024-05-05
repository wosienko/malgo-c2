package ws

import (
	"encoding/json"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) handleDeletedSessionKeyValue(input []byte) error {
	var keyValue entities.SessionKeyValueDeleted
	err := proto.Unmarshal(input, &keyValue)
	if err != nil {
		return fmt.Errorf("error unmarshalling key-value: %v", err)
	}

	if keyValue.SessionId != h.subscribedSession {
		return nil
	}

	response := internalEntities.SessionKeyValueSentToOperator{
		MessageType: "session-key-value-delete",
		SessionId:   keyValue.SessionId,
		Key:         keyValue.Key,
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
