package events

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"malgo-websocket/internal/marshalers"
)

//var marshaler = cqrs.JSONMarshaler{
//	GenerateName: cqrs.StructName,
//}

var marshaler = marshalers.ProtobufMarshaler{GenerateName: cqrs.StructName}

func NewProcessorConfig(redisClient *redis.Client, watermillLogger watermill.LoggerAdapter) cqrs.EventProcessorConfig {
	return cqrs.EventProcessorConfig{
		GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			return fmt.Sprintf("events.%s", params.EventName), nil
		},
		SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        redisClient,
				ConsumerGroup: "svc-c2-websocket.events." + params.HandlerName,
			}, watermillLogger)
		},
		Marshaler: marshaler,
		Logger:    watermillLogger,
	}
}

func NewBusConfig(watermillLogger watermill.LoggerAdapter) cqrs.EventBusConfig {
	return cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return fmt.Sprintf("events.%s", params.EventName), nil
		},
		Marshaler: marshaler,
		Logger:    watermillLogger,
	}
}

func NewBusConfigWithoutLogger() cqrs.EventBusConfig {
	return cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return fmt.Sprintf("events.%s", params.EventName), nil
		},
		Marshaler: marshaler,
	}
}
