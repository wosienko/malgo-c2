package httpproxy

import (
	"context"
	gateway "github.com/VipWW/malgo-c2/services/common/service"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// get project_id and session_id from the URL
	projectID, err := uuid.Parse(params.ByName("project_id"))
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}
	sessionID, err := uuid.Parse(params.ByName("session_id"))
	if err != nil {
		http.Error(w, "Invalid session_id", http.StatusBadRequest)
		return
	}

	// register the project_id and session_id
	_, err = h.grpcClient.RegisterNewSession(context.Background(),
		&gateway.RegisterNewSessionRequest{
			ProjectId: projectID.String(),
			SessionId: sessionID.String(),
		})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// Session duplicate, it's ok
		} else {
			http.Error(w, "Failed to register session", http.StatusInternalServerError)
			return
		}
	}

	// return hello world
	w.Write([]byte("Hello, World!"))
}
