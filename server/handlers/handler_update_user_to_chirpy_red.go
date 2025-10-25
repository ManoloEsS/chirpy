package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ManoloEsS/go_http_server/server"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerUpdateUserToChirpyRed(w http.ResponseWriter, r *http.Request) {
	type UpgradeEvent struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	eventData := UpgradeEvent{}
	err := json.NewDecoder(r.Body).Decode(&eventData)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}

	if eventData.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.Db.UpdateUserIsChirpyRedTrue(r.Context(), eventData.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			server.RespondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
