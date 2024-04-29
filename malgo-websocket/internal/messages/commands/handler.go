package commands

import (
	"context"
	"malgo-websocket/internal/entities"
)

type Handler struct {
	commandRepo CommandRepository
	sessionRepo SessionRepository
}

func NewHandler(
	commandRepo CommandRepository,
	sessionRepo SessionRepository,
) Handler {
	if commandRepo == nil {
		panic("missing commandRepo")
	}
	if sessionRepo == nil {
		panic("missing sessionRepo")
	}

	return Handler{
		commandRepo: commandRepo,
		sessionRepo: sessionRepo,
	}
}

type CommandRepository interface {
	AddCommand(ctx context.Context, command entities.Command) error
}

type SessionRepository interface {
	AddKeyValue(ctx context.Context, value entities.SessionKeyValue) error
	DeleteKeyValue(ctx context.Context, value entities.SessionKeyValue) error
	ModifyKeyValue(ctx context.Context, value entities.SessionKeyValue) error
	RenameSession(ctx context.Context, value entities.SessionName) error
}
