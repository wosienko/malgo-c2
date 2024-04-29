package commands

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
}
