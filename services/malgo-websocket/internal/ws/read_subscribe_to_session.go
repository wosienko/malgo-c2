package ws

import (
	"encoding/json"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/google/uuid"
	"log"
)

func (h *Handler) subscribeSession(input []byte) {
	var inputObj entities.SubscribeToSession
	if err := json.Unmarshal(input, &inputObj); err != nil {
		log.Printf("Could not unmarshal project subscription: %v", err)
		return
	}

	sessionId, err := uuid.Parse(inputObj.SessionId)
	if err != nil {
		log.Printf("Could not parse session id: %v", err)
		return
	}

	h.subscribedSession = sessionId.String()
}
