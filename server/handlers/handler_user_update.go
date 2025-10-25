package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server"
)

func (cfg *ApiConfig) HandlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't find Token", err)
		return
	}

	validatedID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	updatedData := login{}
	err = json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "Couldn't decode request body", err)
		return
	}

	newHashedPassword, err := auth.HashPassword(updatedData.Password)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "Couldn't hash password", err)
		return
	}

	newParams := database.UpdateUserAuthenticationParams{
		Email:          updatedData.Email,
		HashedPassword: newHashedPassword,
		ID:             validatedID,
	}

	newUserData, err := cfg.Db.UpdateUserAuthentication(r.Context(), newParams)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user data in database", err)
		return
	}

	newResponse := ResponseUser{
		ID:          newUserData.ID,
		CreatedAt:   newUserData.CreatedAt,
		UpdatedAt:   newUserData.UpdatedAt,
		Email:       newUserData.Email,
		IsChirpyRed: newUserData.IsChirpyRed.Bool,
	}

	server.RespondWithJSON(w, http.StatusOK, newResponse)
}
