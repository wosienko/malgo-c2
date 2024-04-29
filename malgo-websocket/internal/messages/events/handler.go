package events

import "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

type Handler struct {
	pubSub *gochannel.GoChannel
}

func NewHandler(
	pubSub *gochannel.GoChannel,
) Handler {
	return Handler{
		pubSub: pubSub,
	}
}
