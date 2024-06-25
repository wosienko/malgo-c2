package httpproxy

import (
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (h *Handler) setResultInfo(w http.ResponseWriter, r *http.Request) {
	commandId := httprouter.ParamsFromContext(r.Context()).ByName("command_id")
	resultLength, err := strconv.ParseInt(
		httprouter.ParamsFromContext(r.Context()).ByName("result_length"),
		10,
		64,
	)
	if err != nil {
		http.Error(w, "Invalid int", http.StatusBadRequest)
		return
	}

	_, err = h.grpcClient.ResultInfo(r.Context(), &gateway.ResultInfoRequest{
		CommandId: commandId,
		Length:    resultLength,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
