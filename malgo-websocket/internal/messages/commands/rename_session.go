package commands

import (
	"context"
	"malgo-websocket/internal/entities"
	"malgo-websocket/internal/log"
)

func (h *Handler) RenameSession(ctx context.Context, session *entities.ModifySessionName) error {
	log.FromContext(ctx).Info("Renaming the session.")

	return h.sessionRepo.RenameSession(ctx, entities.SessionName{
		SessionId: session.SessionId,
		Name:      session.Name,
	})
}
