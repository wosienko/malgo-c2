package ws

import (
	"encoding/json"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) handleRenamedSession(input []byte) error {
	var renamedSession entities.SessionNameModified
	err := proto.Unmarshal(input, &renamedSession)
	if err != nil {
		return fmt.Errorf("error unmarshalling renamed session: %v", err)
	}

	if renamedSession.ProjectId != h.subscribedProject {
		return nil
	}

	response := internalEntities.SessionNameSentToOperator{
		MessageType: "session-renamed",
		SessionId:   renamedSession.SessionId,
		Name:        renamedSession.Name,
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
