package handlers

import (
	"net/http"
	"time"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/server"
)

// HandlerValidateRefreshToken exchanges a valid refresh token for a new JWT access token.
// Verifies the refresh token against the database and generates a new 1-hour JWT.
func (cfg *ApiConfig) HandlerValidateRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "Couldn't find Refresh Token", err)
		return
	}

	userData, err := cfg.Db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't find user or token expired", err)
		return
	}

	JWTtoken, err := auth.MakeJWT(userData.ID, cfg.Secret, time.Duration(time.Hour*1))
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't create JWT token", err)
		return
	}

	responseStruct := struct {
		Token string `json:"token"`
	}{
		Token: JWTtoken,
	}

	server.RespondWithJSON(w, http.StatusOK, responseStruct)
}
