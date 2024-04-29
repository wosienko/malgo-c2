package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"malgo-websocket/internal/entities"
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

	response := entities.SessionKeyValueSentToOperator{
		MessageType: "session-key-value-delete",
		SessionId:   keyValue.SessionId,
		Key:         keyValue.Key,
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
