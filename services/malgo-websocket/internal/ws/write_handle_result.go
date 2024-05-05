package ws

import (
	"encoding/json"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) handleCommandResultRetrieved(input []byte) error {
	var commandResult entities.CommandResultRetrieved
	err := proto.Unmarshal(input, &commandResult)
	if err != nil {
		return fmt.Errorf("error unmarshalling command result: %v", err)
	}

	if commandResult.SessionId != h.subscribedSession {
		return nil
	}

	response := internalEntities.CommandResultSentToOperator{
		MessageType: "command-result",
		CommandId:   commandResult.CommandId,
		SessionId:   commandResult.SessionId,
		Result:      commandResult.Result,
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
