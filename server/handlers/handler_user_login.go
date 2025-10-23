package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server"
)

func (cfg *ApiConfig) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		ResponseUser
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	expirationDuration := time.Hour
	if req.Expires_in_seconds > 0 && req.Expires_in_seconds < int(time.Hour) {
		expirationDuration = time.Duration(req.Expires_in_seconds) * time.Second
	}
	token, err := auth.MakeJWT(userData.ID, cfg.Secret, expirationDuration)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't create authentication token", err)
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()

	err = cfg.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID: userData.ID,
		Token:  refreshToken,
	})
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "couldn't create refresh token in database", err)
	}

	JSONresponse := response{
		ResponseUser: ResponseUser{
			ID:        userData.ID,
			CreatedAt: userData.CreatedAt,
			UpdatedAt: userData.UpdatedAt,
			Email:     userData.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}
	server.RespondWithJSON(w, http.StatusOK, JSONresponse)

}
