package httpproxy

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc("PUT", "/:project_id/:session_id", h.register)
	mux.HandlerFunc("GET", "/s/:session_id", h.getCommandInfo)
	mux.HandlerFunc("GET", "/c/:command_id/:offset", h.getCommandDetails)
	mux.HandlerFunc("PATCH", "/c/:command_id/:result_length", h.setResultInfo)
	mux.HandlerFunc("POST", "/c/:command_id/:offset", h.setResultDetailsChunk)

	return mux
}
