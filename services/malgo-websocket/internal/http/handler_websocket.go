package http

import (
	"context"
	"fmt"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/VipWW/malgo-c2/services/malgo-websocket/internal/ws"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Add origin check
		return true
	},
}

func (h *Handler) UpgradeToWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("error while upgrading: %v\n", err)
		return
	}

	log.FromContext(context.Background()).Infof("New websocket connection: %v\n", conn.RemoteAddr().String())

	handler := ws.NewHandler(
		conn,
		r.Context().Value("user_id").(string),
		h.userRepo,
		h.pubSub,
		h.eventBus,
		h.commandBus,
	)

	go handler.ReadFromWebsocket()
	go handler.WriteToWebsocket()
}
