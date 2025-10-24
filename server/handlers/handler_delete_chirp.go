package handlers

import (
	"net/http"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(idString)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "Couldn't parse chirp id", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't find bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	chirpData, err := cfg.Db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		server.RespondWithError(w, http.StatusNotFound, "Not found", err)
		return
	}

	if chirpData.UserID != userID {
		server.RespondWithError(w, http.StatusForbidden, "Users can only delete their own chirps", nil)
		return
	}

	err = cfg.Db.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
		UserID: userID,
		ID:     chirpID,
	})
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
