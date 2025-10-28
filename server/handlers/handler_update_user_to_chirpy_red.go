package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/server"
	"github.com/google/uuid"
)

// HandlerUpdateUserToChirpyRed handles webhook events from Polka API to upgrade users to Chirpy Red membership.
// It validates the API key, processes "user.upgraded" events, and updates the user's membership status in the database.
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

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't retrieve Authorization header", err)
		return
	}

	if key != cfg.PolkaAPI {
		server.RespondWithError(w, http.StatusUnauthorized, "Incorrect API key", err)
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
