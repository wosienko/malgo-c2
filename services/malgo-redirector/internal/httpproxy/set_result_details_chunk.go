package httpproxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (h *Handler) setResultDetailsChunk(w http.ResponseWriter, r *http.Request) {
	commandId := httprouter.ParamsFromContext(r.Context()).ByName("command_id")
	offset, err := strconv.ParseInt(
		httprouter.ParamsFromContext(r.Context()).ByName("offset"),
		10,
		64,
	)
	if err != nil {
		http.Error(w, "Invalid int", http.StatusBadRequest)
		return
	}

	dataBase64 := make([]byte, r.ContentLength)
	fmt.Printf("ContentLength: %v\n", r.ContentLength)
	_, _ = r.Body.Read(dataBase64) // will return EOF

	data := make([]byte, base64.StdEncoding.DecodedLen(len(dataBase64)))
	_, err = base64.StdEncoding.Decode(data, dataBase64)
	if err != nil {
		http.Error(w, "Failed to decode base64", http.StatusInternalServerError)
		return
	}

	// trim null bytes
	data = bytes.Trim(data, "\x00")

	_, err = h.grpcClient.ResultDetailsChunk(r.Context(), &gateway.ResultDetailsChunkRequest{
		CommandId: commandId,
		Offset:    offset,
		Data:      data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
