package commands

import "context"

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
}

type SessionRepository interface {
	UpdateSessionHeartbeat(ctx context.Context, sessionId string) error
}
