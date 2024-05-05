package ws

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/log"
	"github.com/gorilla/websocket"
	"time"
)

const (
	NewCommandsTopic             = "new-commands"
	CommandStatusUpdatedTopic    = "command-status-updated"
	CommandResultRetrievedTopic  = "command-result-retrieved"
	ResultChunkInsertedTopic     = "result-chunk-inserted"
	SessionKeyValueModifiedTopic = "session-key-value-modified"
	SessionKeyValueDeletedTopic  = "session-key-value-deleted"
	RenamedSessionTopic          = "renamed-session"
	SessionHeartbeatUpdatedTopic = "session-heartbeat-updated"
	NewSessionTopic              = "new-session"
)

func (h *Handler) WriteToWebsocket() {
	h.conn.EnableWriteCompression(true)

	ctx, cancel := context.WithCancel(context.Background())

	newCommandsChannel, err := h.pubSub.Subscribe(ctx, NewCommandsTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on new commands channel")
	}
	commandStatusChangedChannel, err := h.pubSub.Subscribe(ctx, CommandStatusUpdatedTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on command status changed channel")
	}
	commandResultRetrievedChannel, err := h.pubSub.Subscribe(ctx, CommandResultRetrievedTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on command result retrieved channel")
	}
	resultChunkInsertedChannel, err := h.pubSub.Subscribe(ctx, ResultChunkInsertedTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on result chunk inserted channel")
	}
	modifiedKeyValueChannel, err := h.pubSub.Subscribe(ctx, SessionKeyValueModifiedTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on modified key-value channel")
	}
	deletedKeyValueChannel, err := h.pubSub.Subscribe(ctx, SessionKeyValueDeletedTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on deleted key-value channel")
	}
	renamedSessionChannel, err := h.pubSub.Subscribe(ctx, RenamedSessionTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on renamed session channel")
	}
	sessionHeartbeatChannel, err := h.pubSub.Subscribe(ctx, SessionHeartbeatUpdatedTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on session heartbeat channel")
	}
	newSessionChannel, err := h.pubSub.Subscribe(ctx, NewSessionTopic)
	if err != nil {
		log.FromContext(context.Background()).Errorf("error listening on new session channel")
	}

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		cancel()
	}()
	for {
		select {
		case message := <-newCommandsChannel:
			if err := h.handleNewCommands(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling new command: %v", err)
			}
			message.Ack()
		case message := <-commandStatusChangedChannel:
			if err := h.handleCommandStatusChanged(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling command status changed: %v", err)
			}
			message.Ack()
		case message := <-commandResultRetrievedChannel:
			if err := h.handleCommandResultRetrieved(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling command result: %v", err)
			}
			message.Ack()
		case message := <-resultChunkInsertedChannel:
			if err := h.handleResultChunkInserted(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling result chunk: %v", err)
			}
			message.Ack()
		case message := <-modifiedKeyValueChannel:
			if err := h.handleModifiedSessionKeyValue(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling modified key-value: %v", err)
			}
			message.Ack()
		case message := <-deletedKeyValueChannel:
			if err := h.handleDeletedSessionKeyValue(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling deleted key-value: %v", err)
			}
			message.Ack()
		case message := <-renamedSessionChannel:
			if err := h.handleRenamedSession(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling renamed session: %v", err)
			}
			message.Ack()
		case message := <-sessionHeartbeatChannel:
			if err := h.handleUpdatedHeartbeat(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling updated heartbeat: %v", err)
			}
			message.Ack()
		case message := <-newSessionChannel:
			if err := h.handleNewSession(message.Payload); err != nil {
				log.FromContext(context.Background()).Errorf("Error handling new session: %v", err)
			}
			message.Ack()
		case <-ticker.C:
			if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.FromContext(context.Background()).Errorf("Error writing ping message: %v", err)
				return
			}
			if err := h.conn.WriteMessage(websocket.TextMessage, []byte("{\"type\": \"pong\"}")); err != nil {
				log.FromContext(context.Background()).Errorf("Error writing ping textMessage: %v", err)
				return
			}
		case <-h.cancel:
			log.FromContext(context.Background()).Infof("Closing writing channel\n")
			return
		}
	}
}
