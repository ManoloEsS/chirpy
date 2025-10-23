package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server"
)

// Method of ApiConfig. takes a http.ResponseWriter and a http.Request and creates
// a new user in the database and responds with the user data in JSON
func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	server.RespondWithJSON(w, 201, newResponseUser)
}
