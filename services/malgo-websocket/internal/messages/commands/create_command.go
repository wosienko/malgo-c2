package commands

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/entities"
	"github.com/VipWW/malgo-c2/services/common/log"
	internalEntities "github.com/VipWW/malgo-c2/services/malgo-websocket/internal/entities"
	"github.com/google/uuid"
)

func (h *Handler) CreateCommand(ctx context.Context, command *entities.CreateCommand) error {
	log.FromContext(ctx).Info("Generating new command.")

	return h.commandRepo.AddCommand(ctx, internalEntities.Command{
		ID:         uuid.MustParse(command.Header.IdempotencyKey),
		SessionId:  uuid.MustParse(command.SessionId),
		Command:    command.Command,
		OperatorId: uuid.MustParse(command.UserId),
	})
}
