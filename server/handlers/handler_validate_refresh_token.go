package handlers

import (
	"net/http"
	"time"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/server"
)

func (cfg *ApiConfig) HandlerValidateRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, "Couldn't find Refresh Token", err)
		return
	}

	id, err := cfg.Db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't find user or token expired", err)
		return
	}

	JWTtoken, err := auth.MakeJWT(id, cfg.Secret, time.Duration(time.Hour*1))
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
