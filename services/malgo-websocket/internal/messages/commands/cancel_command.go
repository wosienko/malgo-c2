package commands

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
)

func (h *Handler) CancelCommand(ctx context.Context, command *entities.CancelCommand) error {
	log.FromContext(ctx).Info("Canceling command.")

	return h.commandRepo.CancelCommand(ctx, command.CommandId)
}
