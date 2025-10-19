package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/server"
)

func (cfg *ApiConfig) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Email == "" || req.Password == "" {
		server.RespondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
	}

	userData, err := cfg.Db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	matches, err := auth.CheckPasswordHash(req.Password, userData.HashedPassword)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Authentication failed", err)
		return
	}

	if !matches {
		server.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	JSONresponse := ResponseUser{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
	}
	server.RespondWithJSON(w, http.StatusOK, JSONresponse)

}
