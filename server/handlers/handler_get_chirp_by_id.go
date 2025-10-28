package handlers

import (
	"net/http"

	"github.com/ManoloEsS/go_http_server/server"
	"github.com/google/uuid"
)

// HandlerGetChirpById retrieves a single chirp by its UUID from the database.
// Parses the chirpID from the URL path and returns the chirp data as JSON.
func (cfg *ApiConfig) HandlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	id, err := uuid.Parse(idString)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.Db.GetChirpByID(r.Context(), id)
	if err != nil {
		server.RespondWithError(w, http.StatusNotFound, "Not Found", err)
		return
	}

	server.RespondWithJSON(w, http.StatusOK, chirp)
}
