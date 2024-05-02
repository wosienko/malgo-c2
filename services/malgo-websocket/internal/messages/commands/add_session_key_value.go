package commands

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
)

func (h *Handler) AddSessionKeyValue(ctx context.Context, keyValue *entities.AddSessionKeyValue) error {
	log.FromContext(ctx).Info("Adding new key value to session.")

	return h.sessionRepo.AddKeyValue(ctx, internalEntities.SessionKeyValue{
		SessionId: keyValue.SessionId,
		Key:       keyValue.Key,
		Value:     keyValue.Value,
	})
}
