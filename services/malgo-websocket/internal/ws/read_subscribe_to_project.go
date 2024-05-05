package ws

import (
	"context"
	"encoding/json"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/google/uuid"
)

func (h *Handler) subscribeProject(input []byte) {
	var inputObj entities.SubscribeToProject
	if err := json.Unmarshal(input, &inputObj); err != nil {
		log.FromContext(context.Background()).Errorf("Could not unmarshal project subscription: %v", err)
		return
	}

	projectId, err := uuid.Parse(inputObj.ProjectId)
	if err != nil {
		log.FromContext(context.Background()).Errorf("Could not parse project id: %v", err)
		return
	}

	isUserAllowed, err := h.userRepo.IsUserAllowedToAccessProject(context.Background(), h.userId, projectId.String())
	if err != nil {
		log.FromContext(context.Background()).Errorf("Could not check if user is allowed to access project: %v", err)
		return
	}
	if !isUserAllowed {
		log.FromContext(context.Background()).Infof("User is not allowed to access project")
		return
	}

	h.subscribedProject = projectId.String()
}
