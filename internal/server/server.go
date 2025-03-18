package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Config represents the server configuration
type Config struct {
	Host    string
	Port    int
	BaseURL string
}

// Server represents an HTTP server
type Server struct {
	config     Config
	router     *http.ServeMux
	httpServer *http.Server
}

// New creates a new server
func New(config Config) *Server {
	router := http.NewServeMux()
	
	return &Server{
		config: config,
		router: router,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Register routes
	s.registerRoutes()
	
	// Start server
	log.Printf("Starting server on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping server...")
	return s.httpServer.Shutdown(ctx)
}

// registerRoutes registers the server routes
func (s *Server) registerRoutes() {
	// Health check
	s.router.HandleFunc("/health", s.handleHealth())
	
	// API routes
	s.router.HandleFunc("/api/", s.handleAPI())
}

// handleHealth returns a handler for health check requests
func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	}
}

// handleAPI returns a handler for API requests
func (s *Server) handleAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Tiger-Tail Microblog API",
			"version": "0.1.0",
		})
	}
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	})
}
