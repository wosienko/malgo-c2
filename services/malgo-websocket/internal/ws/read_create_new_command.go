package ws

import (
	"context"
	"encoding/json"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/google/uuid"
	"log"
)

func (h *Handler) createNewCommand(input []byte) {
	var cmd entities.CreateCommand
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
	cmd.UserId = h.userId

	if err = h.commandBus.Send(context.Background(), &cmd); err != nil {
		log.Printf("Could not send command: %v", err)
		return
	}
}
