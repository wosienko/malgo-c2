package ws

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"malgo-websocket/internal/entities"
)

func (h *Handler) addSessionKeyValue(input []byte) {
	var cmd entities.AddSessionKeyValue
	if err := json.Unmarshal(input, &cmd); err != nil {
		log.Printf("Could not unmarshal command: %v", err)
		return
	}

	_, err := uuid.Parse(cmd.SessionId)
	if err != nil {
		log.Printf("Could not parse UUID: %v", err)
		return
	}

	cmd.Header = entities.NewHeader()

	if err = h.commandBus.Send(context.Background(), &cmd); err != nil {
		log.Printf("Could not send command: %v", err)
		return
	}
}
