package commands

import (
	"context"
	"malgo-websocket/internal/entities"
	"malgo-websocket/internal/log"
)

func (h *Handler) DeleteSessionKeyValue(ctx context.Context, keyValue *entities.DeleteSessionKeyValue) error {
	log.FromContext(ctx).Info("Deleting key value from session.")

	return h.sessionRepo.DeleteKeyValue(ctx, entities.SessionKeyValue{
		SessionId: keyValue.SessionId,
		Key:       keyValue.Key,
	})
}
