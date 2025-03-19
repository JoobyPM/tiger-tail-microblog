package server

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// LivezHandler handles liveness probe requests
func LivezHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK."))
	}
}

// ReadyzHandler handles readiness probe requests
func ReadyzHandler(db DBPinger, cache CachePinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		dbStatus := "up"
		if err := db.Ping(); err != nil {
			dbStatus = "down"
		}

		// Check cache connection
		cacheStatus := "up"
		if err := cache.Ping(); err != nil {
			cacheStatus = "down"
		}

		// Determine overall status
		status := http.StatusOK
		if dbStatus == "down" || cacheStatus == "down" {
			status = http.StatusServiceUnavailable
		}

		// Determine status message
		statusMsg := "ready"
		if status != http.StatusOK {
			statusMsg = "not ready"
		}

		// Respond with status
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": statusMsg,
			"checks": map[string]string{
				"database": dbStatus,
				"cache":    cacheStatus,
			},
		})
	}
}

// DBPinger defines the interface for database ping operations
type DBPinger interface {
	Ping() error
}

// CachePinger defines the interface for cache ping operations
type CachePinger interface {
	Ping() error
}

// PostHandler handles post-related requests
type PostHandler struct {
	postService domain.PostService
	postCache   PostCache
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService domain.PostService, postCache PostCache) *PostHandler {
	return &PostHandler{
		postService: postService,
		postCache:   postCache,
	}
}

// GetPostsHandler handles GET /posts requests
func (h *PostHandler) GetPostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET method
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Parse query parameters
		query := r.URL.Query()
		pageStr := query.Get("page")
		limitStr := query.Get("limit")

		// Set default values
		page := 1
		limit := 10

		// Parse page parameter
		if pageStr != "" {
			pageInt, err := strconv.Atoi(pageStr)
			if err != nil || pageInt < 1 {
				respondError(w, http.StatusBadRequest, "Invalid page parameter")
				return
			}
			page = pageInt
		}

		// Parse limit parameter
		if limitStr != "" {
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil || limitInt < 1 || limitInt > 100 {
				respondError(w, http.StatusBadRequest, "Invalid limit parameter")
				return
			}
			limit = limitInt
		}

		// Try to get posts from cache
		posts, err := h.postCache.GetPostsWithUser()
		if err == nil {
			// Cache hit
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"posts":  posts,
				"page":   page,
				"limit":  limit,
				"total":  len(posts),
				"source": "cache",
			})
			return
		}

		// Cache miss, get posts from service
		posts, total, err := h.postService.List(page, limit)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to get posts")
			return
		}

		// Set posts in cache
		go h.postCache.SetPostsWithUser(posts)

		// Respond with posts
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"posts":  posts,
			"page":   page,
			"limit":  limit,
			"total":  total,
			"source": "database",
		})
	}
}

// GetPostHandler handles GET /posts/:id requests
func (h *PostHandler) GetPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET method
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Extract post ID from URL
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 3 {
			respondError(w, http.StatusBadRequest, "Invalid URL")
			return
		}
		id := parts[len(parts)-1]

		// Try to get post from cache
		cachedPost, err := h.postCache.GetPost(id)
		if err == nil {
			// Cache hit
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"post":   cachedPost,
				"source": "cache",
			})
			return
		}

		// Cache miss, get post from service
		postWithUser, err := h.postService.GetByID(id)
		if err != nil {
			if err == domain.ErrPostNotFound {
				respondError(w, http.StatusNotFound, "Post not found")
			} else {
				respondError(w, http.StatusInternalServerError, "Failed to get post")
			}
			return
		}

		// Set post in cache (we only cache the Post part, not the PostWithUser)
		go h.postCache.SetPost(&postWithUser.Post)

		// Respond with post
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"post":   postWithUser,
			"source": "database",
		})
	}
}

// CreatePostHandler handles POST /posts requests
func (h *PostHandler) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST method
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Check authentication
		userID, err := authenticateRequest(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Parse request body
		var requestBody struct {
			Content string `json:"content"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate content
		if requestBody.Content == "" {
			respondError(w, http.StatusBadRequest, "Content is required")
			return
		}

		// Create post
		post, err := h.postService.Create(userID, requestBody.Content)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create post")
			return
		}

		// Invalidate cache
		go h.postCache.InvalidatePosts()

		// Respond with created post
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"post":    post,
			"message": "Post created successfully",
		})
	}
}

// PostCache defines the interface for post caching
type PostCache interface {
	GetPost(id string) (*domain.Post, error)
	SetPost(post *domain.Post) error
	GetPostsWithUser() ([]*domain.PostWithUser, error)
	SetPostsWithUser(posts []*domain.PostWithUser) error
	InvalidatePosts() error
}

// respondJSON responds with JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError responds with an error
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// authenticateRequest authenticates a request using Basic Auth
func authenticateRequest(r *http.Request) (string, error) {
	// Get username and password from Basic Auth
	username, password, ok := r.BasicAuth()
	if !ok {
		return "", domain.ErrUserNotFound
	}

	// Get expected username and password from environment variables
	expectedUsername := os.Getenv("AUTH_USERNAME")
	if expectedUsername == "" {
		expectedUsername = "admin" // Default if not set
	}
	
	expectedPassword := os.Getenv("AUTH_PASSWORD")
	if expectedPassword == "" {
		expectedPassword = "password" // Default if not set
	}

	// Check if username and password are valid
	if username == expectedUsername && password == expectedPassword {
		return "user_1", nil
	}

	return "", domain.ErrUserNotFound
}
