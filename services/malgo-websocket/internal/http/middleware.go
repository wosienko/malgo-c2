package http

import (
	"context"
	"github.com/VipWW/malgo-c2/services/common/log"
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: Add origin check
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionId string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "auth_session" {
				sessionId = cookie.Value
			}
		}
		userId, err := h.userRepo.GetUserIdIfLoggedInAndOperator(r.Context(), sessionId)
		if err != nil {
			log.FromContext(context.Background()).Errorf("Error: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if userId == "" {
			log.FromContext(context.Background()).Infof("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
