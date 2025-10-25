package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ManoloEsS/go_http_server/internal/database"
	"github.com/ManoloEsS/go_http_server/server/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//load .env file
	godotenv.Load()
	//get database connection url
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbPlatform := os.Getenv("PLATFORM")
	if dbPlatform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secretAuth := os.Getenv("SECRET")
	if secretAuth == "" {
		log.Fatal("SECRET must be set")
	}

	//open
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't connect with database")
	}
	dbQueries := database.New(db)

	//Root path to serve files from
	const filepathRoot = "."
	//default port for serving
	const port = "8080"

	//initialize config to share state
	cfg := &handlers.ApiConfig{
		Db:       dbQueries,
		Platform: dbPlatform,
		Secret:   secretAuth,
	}

	//initialize file server
	fileServer := http.FileServer(http.Dir(filepathRoot))

	//initialize multiplexer to handle requests
	mux := http.NewServeMux()

	//initialize http server struct
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	//main page handler, mapped to app and stripped to be "/"
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	//handler for server health check
	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)
	//handler for hit metrics check
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerRequestMetrics)
	//handler to reset metrics
	mux.HandleFunc("POST /admin/reset", cfg.HandlerResetUsers)
	//handler to validate chirp length
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirp)
	//handler to create a user
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)
	//handler to create chirp
	mux.HandleFunc("POST /api/chirps", cfg.HandlerCreateChirp)
	//handler to retrieve all chirps in database
	mux.HandleFunc("GET /api/chirps", cfg.HandlerGetAllChirps)
	//handler that returns a single chirp using the id as a parameter
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.HandlerGetChirpById)
	//handler that takes a password and email for login
	mux.HandleFunc("POST /api/login", cfg.HandlerUserLogin)
	//handler that takes no body and takes a Header with Authorization: Bearer <token>
	mux.HandleFunc("POST /api/refresh", cfg.HandlerValidateRefreshToken)
	//handler that revokes the refresh token from the database and returns no body
	mux.HandleFunc("POST /api/revoke", cfg.HandlerRevokeRefreshToken)
	//handler that lets users update their passwords
	mux.HandleFunc("PUT /api/users", cfg.HandlerUserUpdate)
	//handler that deletes chirp from database from chirp id
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.HandlerDeleteChirp)
	//handler for webhook that updates user to chirpy red
	mux.HandleFunc("POST /api/polka/webhooks", cfg.HandlerUpdateUserToChirpyRed)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	//listen and serve that blocks the log.Fatal server shutdown
	log.Fatal(server.ListenAndServe())
}
