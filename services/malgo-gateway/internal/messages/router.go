package messages

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/commands"
)

func NewWatermillRouter(
	commandProcessorConfig cqrs.CommandProcessorConfig,
	commandHandler commands.Handler,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	UseMiddlewares(router, watermillLogger)

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(router, commandProcessorConfig)
	if err != nil {
		panic(err)
	}

	err = commandProcessor.AddHandlers(
		cqrs.NewCommandHandler(
			"UpdateSessionHeartbeat",
			commandHandler.UpdateSessionHeartbeat,
		),
	)
	if err != nil {
		panic(err)
	}

	return router
}
