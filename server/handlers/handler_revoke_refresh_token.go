package handlers

import (
	"net/http"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/server"
)

func (cfg *ApiConfig) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "Couldn't find Refresh Token", err)
		return
	}

	err = cfg.Db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
