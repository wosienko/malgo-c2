package commands

import (
	"context"
	"malgo-websocket/internal/entities"
	"malgo-websocket/internal/log"
)

func (h *Handler) AddSessionKeyValue(ctx context.Context, keyValue *entities.AddSessionKeyValue) error {
	log.FromContext(ctx).Info("Adding new key value to session.")

	return h.sessionRepo.AddKeyValue(ctx, entities.SessionKeyValue{
		SessionId: keyValue.SessionId,
		Key:       keyValue.Key,
		Value:     keyValue.Value,
	})
}
