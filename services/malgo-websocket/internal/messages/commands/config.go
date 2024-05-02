package commands

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/VipWW/malgo-c2/services/common/marshalers"
	"github.com/redis/go-redis/v9"
)

//var marshaler = cqrs.JSONMarshaler{
//	GenerateName: cqrs.StructName,
//}

var marshaler = marshalers.ProtobufMarshaler{GenerateName: cqrs.StructName}

func NewProcessorConfig(redisClient *redis.Client, watermillLogger watermill.LoggerAdapter) cqrs.CommandProcessorConfig {
	return cqrs.CommandProcessorConfig{
		GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
			return fmt.Sprintf("commands.%s", params.CommandName), nil
		},
		SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return redisstream.NewSubscriber(
				redisstream.SubscriberConfig{
					Client:        redisClient,
					ConsumerGroup: "svc-c2-websocket.commands." + params.HandlerName,
				},
				watermillLogger,
			)
		},
		Marshaler: marshaler,
		Logger:    watermillLogger,
	}
}

func NewBusConfig(watermillLogger watermill.LoggerAdapter) cqrs.CommandBusConfig {
	return cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return fmt.Sprintf("commands.%s", params.CommandName), nil
		},
		Marshaler: marshaler,
		Logger:    watermillLogger,
	}
}
