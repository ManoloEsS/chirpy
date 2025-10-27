package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetBearerToken extracts the Bearer token from the Authorization header in the provided HTTP headers.
//
// Parameters:
// - headers (http.Header): The HTTP headers from which the Authorization header will be retrieved.
//
// Returns:
// - string: The extracted Bearer token if the Authorization header is valid.
// - error: An error if the Authorization header is missing, empty, or malformed.
//
// The function expects the Authorization header to be in the format:
// "Bearer <token>". If the header is not present, or if it does not follow
// this format, an appropriate error is returned.
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
