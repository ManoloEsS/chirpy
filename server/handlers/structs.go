package handlers

import (
	"time"

	"github.com/google/uuid"
)

// ResponseUser represents the JSON structure for user data in API responses.
// Excludes sensitive fields like hashed passwords for security.
type ResponseUser struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

// login represents user credentials for authentication requests.
type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
