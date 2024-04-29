package commands

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/log"
)

func (h *Handler) DeleteSessionKeyValue(ctx context.Context, keyValue *entities.DeleteSessionKeyValue) error {
	log.FromContext(ctx).Info("Deleting key value from session.")

	return h.sessionRepo.DeleteKeyValue(ctx, internalEntities.SessionKeyValue{
		SessionId: keyValue.SessionId,
		Key:       keyValue.Key,
	})
}
