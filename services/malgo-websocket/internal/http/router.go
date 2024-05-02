package http

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewRouter(
	eventBus *cqrs.EventBus,
	commandBus *cqrs.CommandBus,
	userRepo UserRepo,
	pubSub *gochannel.GoChannel,
) http.Handler {
	mux := httprouter.New()

	handler := Handler{
		eventBus:   eventBus,
		commandBus: commandBus,

		pubSub: pubSub,

		userRepo: userRepo,
	}

	dynamic := alice.New(CORS, handler.Auth)

	mux.HandlerFunc(http.MethodGet, "/health", handler.HealthCheck)

	mux.Handler(http.MethodGet, "/ws", dynamic.ThenFunc(handler.UpgradeToWebsocket))

	mux.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	return mux
}
