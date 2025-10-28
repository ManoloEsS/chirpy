package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ManoloEsS/go_http_server/internal/auth"
	"github.com/ManoloEsS/go_http_server/internal/config"
	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server"
	"github.com/google/uuid"
)

// HandlerCreateChirp creates a new chirp message for the authenticated user.
// Validates JWT, filters profanity, enforces length limits, and saves to database.
func (cfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	//chirp struct to use for decoding request
	type ChirpParams struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	validatedID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	newRequestChirpParams := ChirpParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newRequestChirpParams)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode chirp", err)
		return
	}

	//validate chirp body length and filter profanity
	filteredChirpBody, err := validateChirp(newRequestChirpParams.Body)
	if err != nil {
		server.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	validatedChirpData := database.CreateChirpParams{
		Body:   filteredChirpBody,
		UserID: validatedID,
	}
	//Add chirp to database and return the struct
	validatedChirp, err := cfg.Db.CreateChirp(context.Background(), validatedChirpData)
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "couldn't save chirp to database", err)
		return
	}

	resp := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}{
		ID:        validatedChirp.ID,
		CreatedAt: validatedChirp.CreatedAt,
		UpdatedAt: validatedChirp.CreatedAt,
		Body:      validatedChirp.Body,
		UserID:    validatedID,
	}

	//respond with success code and response instance
	server.RespondWithJSON(w, 201, resp)
}

// validateChirp checks chirp length and filters profane words.
func validateChirp(body string) (string, error) {
	if len(body) > config.MaxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	filteredBody := profaneFilter(body, badWords)

	return filteredBody, nil

}

// profaneFilter replaces profane words with asterisks.
func profaneFilter(prefilter string, profanity map[string]struct{}) string {
	splitString := strings.Split(prefilter, " ")
	for i, word := range splitString {
		loweredWord := strings.ToLower(word)
		if _, ok := profanity[loweredWord]; ok {
			splitString[i] = "****"
		}
	}
	return strings.Join(splitString, " ")
}

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}
