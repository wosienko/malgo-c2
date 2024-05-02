package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) handleResultChunkInserted(input []byte) error {
	var chunk entities.ResultChunkInserted
	if err := proto.Unmarshal(input, &chunk); err != nil {
		return fmt.Errorf("could not unmarshal result chunk: %v", err)
	}

	if chunk.SessionId != h.subscribedSession {
		return nil
	}

	response := internalEntities.CommandResultChunkSentToOperator{
		MessageType: "result-chunk",
		CommandId:   chunk.CommandId,
		SessionId:   chunk.SessionId,
		CreatedAt:   chunk.CreatedAt.AsTime().UTC().String(),
		Progress:    chunk.Progress,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("could not marshal response: %v", err)
	}

	log.FromContext(context.Background()).Info("Sending result chunk through websockets")

	if err := h.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return fmt.Errorf("error sending message through websocket: %v", err)
	}
	return nil
}
