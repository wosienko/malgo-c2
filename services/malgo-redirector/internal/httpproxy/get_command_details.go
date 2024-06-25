package httpproxy

import (
	"encoding/base64"
	"encoding/json"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (h *Handler) getCommandDetails(w http.ResponseWriter, r *http.Request) {
	// Get the command ID from the URL
	commandID := httprouter.ParamsFromContext(r.Context()).ByName("command_id")
	// offset is int64
	offset, err := strconv.ParseInt(
		httprouter.ParamsFromContext(r.Context()).ByName("offset"),
		10,
		64,
	)
	if err != nil {
		http.Error(w, "invalid number", http.StatusBadRequest)
		return
	}

	// Get the command details from the database
	command, err := h.grpcClient.CommandDetailsChunk(r.Context(), &gateway.CommandDetailsChunkRequest{
		CommandId: commandID,
		Offset:    offset,
		Length:    maxHTTPMessageSize,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(command)
	if err != nil {
		http.Error(w, "failed to marshal", http.StatusInternalServerError)
		return
	}

	b64encoded := base64.StdEncoding.EncodeToString(jsonResponse)
	w.Write([]byte(b64encoded))
}
