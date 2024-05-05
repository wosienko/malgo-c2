package ws

import (
	"context"
	"encoding/json"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/google/uuid"
)

func (h *Handler) cancelCommand(input []byte) {
	var cmd entities.CancelCommand
	if err := json.Unmarshal(input, &cmd); err != nil {
		log.FromContext(context.Background()).Errorf("Could not unmarshal cancel command: %v", err)
		return
	}

	_, err := uuid.Parse(cmd.CommandId)
	if err != nil {
		log.FromContext(context.Background()).Errorf("Could not parse UUID: %v", err)
		return
	}

	cmd.Header = entities.NewHeader()

	if err = h.commandBus.Send(context.Background(), &cmd); err != nil {
		log.FromContext(context.Background()).Errorf("Could not send cancel command: %v", err)
		return
	}
}
