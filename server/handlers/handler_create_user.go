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
//
// Request body should contain:
//   - email: string    - User's email address
//   - password: string - User's password (will be hashed before storage)
//
// Returns:
//   - 201 Created with JSON response containing user ID and email on success
//   - 400 Bad Request if the request body is invalid or missing required fields
//   - 500 Internal Server Error if there's an issue with the database
//
// The function performs the following operations:
//  1. Decodes JSON request body
//  2. Hashes the password using bcrypt
//  3. Creates a new user in the database
//  4. Returns the created user information
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
