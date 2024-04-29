package commands

import (
	"context"
	"malgo-websocket/internal/entities"
	"malgo-websocket/internal/log"
)

func (h *Handler) ModifySessionKeyValue(ctx context.Context, keyValue *entities.ModifySessionKeyValue) error {
	log.FromContext(ctx).Info("Modifying key value of session.")

	return h.sessionRepo.ModifyKeyValue(ctx, entities.SessionKeyValue{
		SessionId: keyValue.SessionId,
		Key:       keyValue.Key,
		Value:     keyValue.Value,
	})
}
