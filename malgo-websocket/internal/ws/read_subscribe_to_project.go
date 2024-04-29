package ws

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"malgo-websocket/internal/entities"
)

func (h *Handler) subscribeProject(input []byte) {
	var inputObj entities.SubscribeToProject
	if err := json.Unmarshal(input, &inputObj); err != nil {
		log.Printf("Could not unmarshal project subscription: %v", err)
		return
	}

	projectId, err := uuid.Parse(inputObj.ProjectId)
	if err != nil {
		log.Printf("Could not parse project id: %v", err)
		return
	}

	h.subscribedProject = projectId.String()
}
