package httpproxy

import (
	"encoding/json"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (h *Handler) getCommandInfo(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the URL parameters
	sessionID := httprouter.ParamsFromContext(r.Context()).ByName("session_id")

	// Get the command info from the session ID
	commandInfo, err := h.grpcClient.CommandInfo(r.Context(),
		&gateway.CommandInfoRequest{SessionId: sessionID},
	)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// Session duplicate, it's ok
		} else {
			http.Error(w, "failed to get command info", http.StatusInternalServerError)
			return
		}
	}

	// Write the command info to the response
	jsonResponse, err := json.Marshal(commandInfo)
	if err != nil {
		http.Error(w, "failed to marshal command info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
