package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// MakeRefreshToken generates a cryptographically secure 32-byte hex-encoded refresh token.
// Used by the login handler to create long-lived refresh tokens stored in the database.
func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)

	encodedStr := hex.EncodeToString(key)

	return encodedStr
}
