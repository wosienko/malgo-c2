package commands

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/entities"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/log"
)

func (h *Handler) ModifySessionKeyValue(ctx context.Context, keyValue *entities.ModifySessionKeyValue) error {
	log.FromContext(ctx).Info("Modifying key value of session.")

	return h.sessionRepo.ModifyKeyValue(ctx, internalEntities.SessionKeyValue{
		SessionId: keyValue.SessionId,
		Key:       keyValue.Key,
		Value:     keyValue.Value,
	})
}
