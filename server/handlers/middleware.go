package handlers

import "net/http"

// MiddlewareMetricsInc increments the file server hit counter for each request.
// Wraps the provided handler and tracks the total number of requests received.
func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
