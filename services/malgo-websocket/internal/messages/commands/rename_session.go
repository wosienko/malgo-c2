package commands

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
)

func (h *Handler) RenameSession(ctx context.Context, session *entities.ModifySessionName) error {
	log.FromContext(ctx).Info("Renaming the session.")

	return h.sessionRepo.RenameSession(ctx, internalEntities.SessionName{
		SessionId: session.SessionId,
		Name:      session.Name,
	})
}
