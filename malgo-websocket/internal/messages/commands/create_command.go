package commands

import (
	"context"
	"github.com/google/uuid"
	"malgo-websocket/internal/entities"
	"malgo-websocket/internal/log"
)

func (h *Handler) CreateCommand(ctx context.Context, command *entities.CreateCommand) error {
	log.FromContext(ctx).Info("Generating new command.")

	return h.commandRepo.AddCommand(ctx, entities.Command{
		ID:         uuid.MustParse(command.Header.IdempotencyKey),
		SessionId:  uuid.MustParse(command.SessionId),
		Command:    command.Command,
		OperatorId: uuid.MustParse(command.UserId),
	})
}
