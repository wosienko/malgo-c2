package messages

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewRedisPublisher(rdb *redis.Client, watermillLogger watermill.LoggerAdapter) message.Publisher {
	var pub message.Publisher
	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client:        rdb,
		DefaultMaxlen: int64(1000),
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	return pub
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func NewRedisSubscriber(rdb *redis.Client, watermillLogger watermill.LoggerAdapter) message.Subscriber {
	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	return sub
}
