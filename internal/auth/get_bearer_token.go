package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authToken := headers.Get("Authorization")
	if authToken == "" {
		return "", errors.New("No authorization header in http.Header")
	}

	splitToken := strings.Split(strings.Trim(authToken, " "), " ")
	if len(splitToken) < 2 || splitToken[0] != "Bearer" {
		return "", errors.New("Malformed authorization header")
	}

	return splitToken[1], nil
}
