package commands

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/entities"
)

func (h *Handler) UpdateSessionHeartbeat(ctx context.Context, session *entities.UpdateSessionHeartbeat) error {
	return h.sessionRepo.UpdateSessionHeartbeat(ctx, session.SessionId)
}
