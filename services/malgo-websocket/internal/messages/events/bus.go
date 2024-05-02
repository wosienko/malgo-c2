package events

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewBus(pub message.Publisher) *cqrs.EventBus {
	eventBus, err := cqrs.NewEventBusWithConfig(
		pub,
		NewBusConfigWithoutLogger(),
	)
	if err != nil {
		panic(err)
	}
	return eventBus
}

func NewBusWithConfig(pub message.Publisher, config cqrs.EventBusConfig) *cqrs.EventBus {
	eventBus, err := cqrs.NewEventBusWithConfig(
		pub,
		config,
	)
	if err != nil {
		panic(err)
	}
	return eventBus
}
