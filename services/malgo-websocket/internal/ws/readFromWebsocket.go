package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type messageType struct {
	Type string `json:"type"`
}

func (h *Handler) ReadFromWebsocket() {
	defer func() {
		fmt.Printf("Closing reading channel\n")
		_ = h.conn.Close()
	}()

	h.conn.SetPongHandler(func(string) error { _ = h.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := h.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			h.cancel <- struct{}{}
			break
		}
		var msgType messageType
		err = json.Unmarshal(message, &msgType)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		switch msgType.Type {
		case "command":
			h.createNewCommand(message)
		case "cancel-command":
			h.cancelCommand(message)
		case "subscribeSession":
			h.subscribeSession(message)
		case "unsubscribeSession":
			h.unsubscribeSession()
		case "subscribeProject":
			h.subscribeProject(message)
		case "unsubscribeProject":
			h.unsubscribeProject()
		case "session-new-key":
			h.addSessionKeyValue(message)
		case "session-delete-key":
			h.deleteSessionKeyValue(message)
		case "session-update-key":
			h.updateSessionKeyValue(message)
		case "session-rename":
			h.renameSession(message)
		case "ping":
			continue
		}
	}
}
