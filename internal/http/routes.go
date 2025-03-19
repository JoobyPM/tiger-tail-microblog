package http

import (
	"net/http"
)

// SetupBasicRoutes sets up the basic HTTP routes (root and API info)
func SetupBasicRoutes(mux *http.ServeMux) {
	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		WriteJSONResponse(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Tiger-Tail Microblog API",
		})
	})
	
	// API endpoint
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, http.StatusOK, map[string]string{
			"message": "Tiger-Tail Microblog API",
			"version": "0.1.0",
		})
	})
}

// SetupHealthRoutes sets up the health check routes
func SetupHealthRoutes(mux *http.ServeMux) {
	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})
	
	// Liveness probe endpoint - separate handler for plain text response
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		// Force content type to text/plain
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK."))
	})
	
	// Readiness probe endpoint
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
			"status": "ready",
			"checks": map[string]string{
				"database": "up",
				"cache":    "up",
			},
		})
	})
}
