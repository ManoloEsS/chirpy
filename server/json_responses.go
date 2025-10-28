package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithJSON marshals the payload to JSON and writes it to the response with the given status code.
// Sets Content-Type to application/json and logs errors if marshaling fails.
func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError responds with a JSON error message and logs the error if provided.
// Logs 5XX errors with additional context and wraps the message in a standard error response structure.
func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
