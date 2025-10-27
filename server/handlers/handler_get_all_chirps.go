package handlers

import (
	"net/http"
	"sort"

	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error
	urlQuery := r.URL.Query()

	if urlQuery.Has("author_id") {
		authorUUID, err := uuid.Parse(r.URL.Query().Get("author_id"))
		if err != nil {
			server.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse id from url", err)
			return
		}

		chirps, err = cfg.Db.GetAllChirpsByUserID(r.Context(), authorUUID)
		if err != nil {
			server.RespondWithError(w, http.StatusInternalServerError, "Couldn't get all Chirps for user", err)
			return
		}
	}

	sortOrder := urlQuery.Get("sort")
	chirps, err = cfg.Db.GetAllChirps(r.Context())
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps", err)
		return
	}

	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}

	server.RespondWithJSON(w, http.StatusOK, chirps)
}
