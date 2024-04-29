package messages

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/commands"
	"github.com/VipWW/malgo-c2/services/malgo-gateway/internal/messages/events"
)

func NewWatermillRouter(
	eventProcessorConfig cqrs.EventProcessorConfig,
	eventHandler events.Handler,
	commandProcessorConfig cqrs.CommandProcessorConfig,
	commandHandler commands.Handler,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	UseMiddlewares(router, watermillLogger)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	err = eventProcessor.AddHandlers()
	if err != nil {
		panic(err)
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(router, commandProcessorConfig)
	if err != nil {
		panic(err)
	}

	err = commandProcessor.AddHandlers()
	if err != nil {
		panic(err)
	}

	return router
}
