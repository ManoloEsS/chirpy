package handlers

import (
	"sync/atomic"

	"github.com/ManoloEsS/go_http_server/internal/database"
)

// ApiConfig holds shared application state and configuration.
// Provides database access, secrets, and metrics across all handlers.
type ApiConfig struct {
	fileServerHits atomic.Int32
	Db             *database.Queries
	Platform       string
	Secret         string
	PolkaAPI       string
}
