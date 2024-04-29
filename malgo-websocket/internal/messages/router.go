package messages

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"malgo-websocket/internal/messages/commands"
	"malgo-websocket/internal/messages/events"
	"malgo-websocket/internal/messages/outbox"
)

func NewWatermillRouter(
	postgresSubscriber message.Subscriber,
	redisPublisher message.Publisher,
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

	useMiddlewares(router, watermillLogger)

	outbox.AddForwarderHandler(postgresSubscriber, redisPublisher, router, watermillLogger)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	err = eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"SendNewCommandToWebsocket",
			eventHandler.SendNewCommandsToWebsocket,
		),
		cqrs.NewEventHandler(
			"SendModifiedSessionKeyValueToWebsocket",
			eventHandler.SendModifiedSessionKeyValueToWebsocket,
		),
		cqrs.NewEventHandler(
			"SendDeletedSessionKeyValueToWebsocket",
			eventHandler.SendDeletedSessionKeyValueToWebsocket,
		),
		cqrs.NewEventHandler(
			"SendRenamedSessionToWebsocket",
			eventHandler.SendRenamedSessionToWebsocket,
		),
	)
	if err != nil {
		panic(err)
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(router, commandProcessorConfig)
	if err != nil {
		panic(err)
	}

	err = commandProcessor.AddHandlers(
		cqrs.NewCommandHandler(
			"CreateCommand",
			commandHandler.CreateCommand,
		),
		cqrs.NewCommandHandler(
			"AddSessionKeyValue",
			commandHandler.AddSessionKeyValue,
		),
		cqrs.NewCommandHandler(
			"DeleteSessionKeyValue",
			commandHandler.DeleteSessionKeyValue,
		),
		cqrs.NewCommandHandler(
			"ModifySessionKeyValue",
			commandHandler.ModifySessionKeyValue,
		),
		cqrs.NewCommandHandler(
			"RenameSession",
			commandHandler.RenameSession,
		),
	)
	if err != nil {
		panic(err)
	}

	return router
}
