package ws

import (
	"context"
	"encoding/json"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/google/uuid"
)

func (h *Handler) subscribeSession(input []byte) {
	var inputObj entities.SubscribeToSession
	if err := json.Unmarshal(input, &inputObj); err != nil {
		log.FromContext(context.Background()).Errorf("Could not unmarshal project subscription: %v", err)
		return
	}

	sessionId, err := uuid.Parse(inputObj.SessionId)
	if err != nil {
		log.FromContext(context.Background()).Errorf("Could not parse session id: %v", err)
		return
	}

	h.subscribedSession = sessionId.String()
}
