package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/cache"
	"github.com/JoobyPM/tiger-tail-microblog/internal/config"
	"github.com/JoobyPM/tiger-tail-microblog/internal/db"
	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
	httputil "github.com/JoobyPM/tiger-tail-microblog/internal/http"
	"github.com/JoobyPM/tiger-tail-microblog/internal/util"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Constants for pagination limits
const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// initApp initializes the application components
func initApp() (string, error) {
	// Load configuration from environment variables
	cfg := config.LoadConfigFromEnv()
	
	// Get server port
	port := fmt.Sprintf("%d", cfg.Server.Port)
	
	// Initialize database
	postgres, err := initDatabase(&cfg.Database, cfg.UseRealDB)
	if err != nil {
		return "", err
	}
	
	// Initialize Redis
	redisClient, err := initRedis(&cfg.Cache, cfg.UseRealCache)
	if err != nil {
		return "", err
	}
	
	// Create repositories and caches
	postRepo := db.NewPostRepository(postgres)
	postCache := cache.NewPostCache(redisClient)
	
	// Setup routes with real implementations
	setupRoutes(http.DefaultServeMux, postRepo, postCache, &cfg.Auth)

	return port, nil
}

// initDatabase initializes the database connection
func initDatabase(dbCreds *config.DatabaseCredentials, useRealDB bool) (*db.PostgresDB, error) {
	// Log the connection details (sanitized)
	log.Printf("Connecting to PostgreSQL with DSN: %s", dbCreds.GetSanitizedDSN())
	
	if !useRealDB {
		log.Printf("Stub: Would connect to PostgreSQL with DSN: %s", dbCreds.GetSanitizedDSN())
		return db.NewPostgresStub(), nil
	}
	
	// Initialize database connection
	postgres, err := db.NewPostgresConnection(dbCreds.GetDSN())
	if err != nil {
		log.Printf("Error: Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return postgres, nil
}

// initRedis initializes the Redis connection
func initRedis(redisCreds *config.RedisCredentials, useRealCache bool) (*cache.RedisClient, error) {
	// Log the connection details
	log.Printf("Connecting to Redis at %s:%s (DB: %d)", redisCreds.Host, redisCreds.Port, redisCreds.DB)
	
	if !useRealCache {
		log.Printf("Stub: Would connect to Redis at %s:%s (DB: %d)", redisCreds.Host, redisCreds.Port, redisCreds.DB)
		return cache.NewRedisStub(), nil
	}
	
	// Initialize Redis connection
	redisClient, err := cache.NewRedisClient(
		redisCreds.GetAddr(), 
		redisCreds.Password.Value(), 
		redisCreds.DB,
	)
	if err != nil {
		log.Printf("Error: Failed to connect to Redis: %v", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return redisClient, nil
}

// setupRoutes sets up the HTTP routes
func setupRoutes(mux *http.ServeMux, postRepo *db.PostRepository, postCache *cache.PostCache, authCreds *config.AuthCredentials) {
	// Setup basic routes
	httputil.SetupBasicRoutes(mux)
	
	// Setup API routes
	setupAPIRoutes(mux, postRepo, postCache, authCreds)
	
	// Setup health check routes
	httputil.SetupHealthRoutes(mux)
}

// setupAPIRoutes sets up the API routes
func setupAPIRoutes(mux *http.ServeMux, postRepo *db.PostRepository, postCache *cache.PostCache, authCreds *config.AuthCredentials) {
	// Posts endpoint
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetPosts(w, r, postRepo, postCache)
		case http.MethodPost:
			handleCreatePost(w, r, postRepo, postCache, authCreds)
		default:
			httputil.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
				"error": "Method not allowed",
			})
		}
	})
}

// handleGetPosts handles GET requests to the posts endpoint
func handleGetPosts(w http.ResponseWriter, r *http.Request, postRepo *db.PostRepository, postCache *cache.PostCache) {
	// Parse pagination parameters
	page, limit := parsePaginationParams(r)
	offset := (page - 1) * limit

	// Try to get posts from cache
	posts, err := postCache.GetPostsWithUser()
	if err == nil {
		// Cache hit
		httputil.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
			"posts":  posts,
			"page":   page,
			"limit":  limit,
			"total":  len(posts),
			"source": "cache",
		})
		return
	}

	// Cache miss, get posts from database
	posts, err = postRepo.List(offset, limit)
	if err != nil {
		httputil.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get posts",
		})
		return
	}

	// Get total count
	total, err := postRepo.Count()
	if err != nil {
		total = len(posts)
	}

	// Set posts in cache
	go postCache.SetPostsWithUser(posts)

	// Return posts
	httputil.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"posts":  posts,
		"page":   page,
		"limit":  limit,
		"total":  total,
		"source": "database",
	})
}

// handleCreatePost handles POST requests to the posts endpoint
func handleCreatePost(w http.ResponseWriter, r *http.Request, postRepo *db.PostRepository, postCache *cache.PostCache, authCreds *config.AuthCredentials) {
	// Check authentication
	if !authenticateRequest(w, r, authCreds) {
		return
	}

	// Parse and validate request body
	var requestBody struct {
		Content string `json:"content"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httputil.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if requestBody.Content == "" {
		httputil.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Content is required",
		})
		return
	}

	// Create post
	post := &domain.Post{
		ID:        util.GeneratePostID(),
		UserID:    "user_1", // Use a fixed user ID for simplicity
		Content:   requestBody.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save post to database
	if err := postRepo.Create(post); err != nil {
		httputil.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create post",
		})
		return
	}

	// Update cache instead of just invalidating it
	// This ensures the new post is immediately available in the cache
	posts, err := postRepo.List(0, DefaultPageSize)
	if err == nil {
		go postCache.SetPostsWithUser(posts)
	} else {
		// If we can't get the updated list, just invalidate the cache
		go postCache.InvalidatePosts()
	}

	// Return success
	httputil.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"post":    post,
		"message": "Post created successfully",
	})
}

// authenticateRequest authenticates a request using Basic Auth
func authenticateRequest(w http.ResponseWriter, r *http.Request, authCreds *config.AuthCredentials) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		httputil.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return false
	}

	// Get expected username and password from credentials
	expectedUsername := authCreds.Username.Value()
	expectedPassword := authCreds.Password.Value()

	// Check if username and password are valid
	if username != expectedUsername || password != expectedPassword {
		httputil.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return false
	}

	return true
}

// parsePaginationParams parses pagination parameters from a request
func parsePaginationParams(r *http.Request) (page, limit int) {
	// Set default values
	page = 1
	limit = DefaultPageSize

	// Parse query parameters
	query := r.URL.Query()
	
	// Parse page parameter
	if pageStr := query.Get("page"); pageStr != "" {
		if pageInt, err := strconv.Atoi(pageStr); err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	// Parse limit parameter
	if limitStr := query.Get("limit"); limitStr != "" {
		if limitInt, err := strconv.Atoi(limitStr); err == nil && limitInt > 0 && limitInt <= MaxPageSize {
			limit = limitInt
		}
	}

	return page, limit
}


// startServer starts the HTTP server
func startServer(port string) {
	go func() {
		fmt.Printf("Starting server on port %s...\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
}

// runServer initializes and runs the server, returning a shutdown function
func runServer() (shutdown func(), err error) {
	fmt.Println("Starting TigerTail...")

	// Initialize the application
	port, err := initApp()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize application: %w", err)
	}

	// Start server
	startServer(port)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Return a shutdown function that can be called to clean up
	return func() {
		// Clean up signal handler
		signal.Stop(sigChan)
		close(sigChan)
		fmt.Println("\nShutting down TigerTail...")
	}, nil
}

func main() {
	// Run the server
	shutdown, err := runServer()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan
	
	// Call shutdown function
	if shutdown != nil {
		shutdown()
	}
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
