package handlers

import (
	"sync/atomic"

	"github.com/ManoloEsS/go_http_server/internal/database"
)

// api config to store a state
type ApiConfig struct {
	fileServerHits atomic.Int32
	Db             *database.Queries
	Platform       string
	Secret         string
}
