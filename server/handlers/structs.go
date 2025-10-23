package handlers

import (
	"time"

	"github.com/google/uuid"
)

// struct used to create the json response with appropriate fields from user created by cfg.Db.CreateUser
type ResponseUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
