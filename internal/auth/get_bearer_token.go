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
	strippedToken := splitToken[1]

	return strippedToken, nil
}
