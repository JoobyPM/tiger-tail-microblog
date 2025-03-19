package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// Config represents the server configuration
type Config struct {
	Host    string
	Port    int
	BaseURL string
}

// Server represents an HTTP server
type Server struct {
	config      Config
	router      *http.ServeMux
	httpServer  *http.Server
	postService domain.PostService
	postCache   PostCache
	db          DBPinger
	cache       CachePinger
}

// New creates a new server
func New(config Config, postService domain.PostService, postCache PostCache, db DBPinger, cache CachePinger) *Server {
	router := http.NewServeMux()
	
	return &Server{
		config:      config,
		router:      router,
		postService: postService,
		postCache:   postCache,
		db:          db,
		cache:       cache,
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
	// Health checks
	s.router.HandleFunc("/health", s.handleHealth())
	s.router.HandleFunc("/livez", LivezHandler())
	s.router.HandleFunc("/readyz", ReadyzHandler(s.db, s.cache))
	
	// API routes
	s.router.HandleFunc("/api/", s.handleAPI())
	
	// Create post handler
	postHandler := NewPostHandler(s.postService, s.postCache)
	
	// Post routes
	s.router.HandleFunc("/api/posts", postHandler.GetPostsHandler())
	s.router.HandleFunc("/api/posts/create", postHandler.CreatePostHandler())
	
	// Individual post route - must be last to avoid conflicts
	s.router.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		// Extract post ID from URL
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 4 || parts[3] == "" || parts[3] == "create" {
			// Not a post ID request, let other handlers handle it
			http.NotFound(w, r)
			return
		}
		
		// Handle the post request
		postHandler.GetPostHandler()(w, r)
	})
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

// serverRespondJSON sends a JSON response (renamed to avoid conflict with handlers.go)
func serverRespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// serverRespondError sends an error response (renamed to avoid conflict with handlers.go)
func serverRespondError(w http.ResponseWriter, status int, message string) {
	serverRespondJSON(w, status, map[string]string{
		"error":   http.StatusText(status),
		"message": message,
	})
}
