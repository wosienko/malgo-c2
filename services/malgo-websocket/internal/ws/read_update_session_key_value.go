package ws

import (
	"context"
	"encoding/json"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/google/uuid"
)

func (h *Handler) updateSessionKeyValue(input []byte) {
	var cmd entities.ModifySessionKeyValue
	if err := json.Unmarshal(input, &cmd); err != nil {
		log.FromContext(context.Background()).Errorf("Could not unmarshal command: %v", err)
		return
	}

	_, err := uuid.Parse(cmd.SessionId)
	if err != nil {
		log.FromContext(context.Background()).Errorf("Could not parse UUID: %v", err)
		return
	}

	cmd.Header = entities.NewHeader()

	if err = h.commandBus.Send(context.Background(), &cmd); err != nil {
		log.FromContext(context.Background()).Errorf("Could not send command: %v", err)
		return
	}
}
