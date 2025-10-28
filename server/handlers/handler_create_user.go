package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server"
)

// HandlerCreateUser handles the creation of new user accounts.
// It processes POST requests with JSON payload containing user credentials.
func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	userRequest := login{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		server.RespondWithError(w, 500, "Couldn't decode user email", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "couldn't encrypt password", err)
		return
	}

	newUser, err := cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userRequest.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		server.RespondWithError(w, 500, "Couldn't create user", err)
		return
	}

	newResponseUser := ResponseUser{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed.Bool,
	}

	server.RespondWithJSON(w, 201, newResponseUser)
}
