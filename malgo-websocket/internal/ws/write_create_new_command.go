package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"malgo-websocket/internal/entities"
)

func (h *Handler) handleNewCommands(input []byte) error {
	var command entities.CommandCreated
	err := proto.Unmarshal(input, &command)
	if err != nil {
		return fmt.Errorf("error unmarshalling command: %v", err)

	}
	if command.SessionId != h.subscribedSession {
		return nil
	}

	response := entities.CommandSentToOperator{
		MessageType: "new-command",
		ID:          command.CommandId,
		Type:        command.Type,
		Status:      command.Status,
		Command:     command.Command,
		Operator:    command.OperatorName,
		CreatedAt:   command.CreatedAt.AsTime().UTC().String(),
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
