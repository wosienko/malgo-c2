package ws

import (
	"context"
	"encoding/json"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/google/uuid"
	"log"
)

func (h *Handler) cancelCommand(input []byte) {
	var cmd entities.CancelCommand
	if err := json.Unmarshal(input, &cmd); err != nil {
		log.Printf("Could not unmarshal cancel command: %v", err)
		return
	}

	_, err := uuid.Parse(cmd.CommandId)
	if err != nil {
		log.Printf("Could not parse UUID: %v", err)
		return
	}

	cmd.Header = entities.NewHeader()

	if err = h.commandBus.Send(context.Background(), &cmd); err != nil {
		log.Printf("Could not send cancel command: %v", err)
		return
	}
}
