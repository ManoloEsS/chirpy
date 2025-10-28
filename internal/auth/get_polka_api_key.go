package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts the API key from the Authorization header.
// Expected format: "ApiKey <key>". Used by HandlerUpdateUserToChirpyRed for webhook authentication.
func GetAPIKey(headers http.Header) (string, error) {
	authData := headers.Get("Authorization")
	if authData == "" {
		return "", errors.New("No Authorization header in http header")
	}

	splitData := strings.Split(authData, " ")
	if len(splitData) < 2 || splitData[0] != "ApiKey" {
		return "", errors.New("Malformed Authorization header")
	}

	return splitData[1], nil
}
