package ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const (
	NewCommandsTopic             = "new-commands"
	SessionKeyValueModifiedTopic = "session-key-value-modified"
	SessionKeyValueDeletedTopic  = "session-key-value-deleted"
	RenamedSessionTopic          = "renamed-session"
	SessionHeartbeatUpdatedTopic = "session-heartbeat-updated"
)

func (h *Handler) WriteToWebsocket() {
	h.conn.EnableWriteCompression(true)

	ctx, cancel := context.WithCancel(context.Background())

	newCommandsChannel, err := h.pubSub.Subscribe(ctx, NewCommandsTopic)
	if err != nil {
		fmt.Printf("error listening on new commands channel")
	}
	modifiedKeyValueChannel, err := h.pubSub.Subscribe(ctx, SessionKeyValueModifiedTopic)
	if err != nil {
		fmt.Printf("error listening on modified key-value channel")
	}
	deletedKeyValueChannel, err := h.pubSub.Subscribe(ctx, SessionKeyValueDeletedTopic)
	if err != nil {
		fmt.Printf("error listening on deleted key-value channel")
	}
	renamedSessionChannel, err := h.pubSub.Subscribe(ctx, RenamedSessionTopic)
	if err != nil {
		fmt.Printf("error listening on renamed session channel")
	}
	sessionHeartbeatChannel, err := h.pubSub.Subscribe(ctx, SessionHeartbeatUpdatedTopic)
	if err != nil {
		fmt.Printf("error listening on session heartbeat channel")
	}

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		cancel()
	}()
	for {
		select {
		case message := <-newCommandsChannel:
			fmt.Printf("Received new command to send through websocket\n")
			if err := h.handleNewCommands(message.Payload); err != nil {
				fmt.Printf("Error handling new command: %v", err)
				message.Nack()
			}
			message.Ack()
		case message := <-modifiedKeyValueChannel:
			fmt.Printf("Received modified key-value to send through websocket\n")
			if err := h.handleModifiedSessionKeyValue(message.Payload); err != nil {
				fmt.Printf("Error handling modified key-value: %v", err)
				message.Nack()
			}
			message.Ack()
		case message := <-deletedKeyValueChannel:
			fmt.Printf("Received deleted key-value to send through websocket\n")
			if err := h.handleDeletedSessionKeyValue(message.Payload); err != nil {
				fmt.Printf("Error handling deleted key-value: %v", err)
				message.Nack()
			}
			message.Ack()
		case message := <-renamedSessionChannel:
			fmt.Printf("Received renamed session to send through websocket\n")
			if err := h.handleRenamedSession(message.Payload); err != nil {
				fmt.Printf("Error handling renamed session: %v", err)
				message.Nack()
			}
			message.Ack()
		case message := <-sessionHeartbeatChannel:
			fmt.Printf("Received updated heartbeat to send through websocket\n")
			if err := h.handleUpdatedHeartbeat(message.Payload); err != nil {
				fmt.Printf("Error handling updated heartbeat: %v", err)
				message.Nack()
			}
			message.Ack()
		case <-ticker.C:
			if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("Error writing ping message: %v", err)
				return
			}
			if err := h.conn.WriteMessage(websocket.TextMessage, []byte("{\"type\": \"pong\"}")); err != nil {
				fmt.Printf("Error writing ping textMessage: %v", err)
				return
			}
		case <-h.cancel:
			fmt.Printf("Closing writing channel\n")
			return
		}
	}
}
