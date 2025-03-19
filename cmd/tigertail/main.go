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
	"github.com/JoobyPM/tiger-tail-microblog/internal/db"
	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// initApp initializes the application components
func initApp() (string, error) {
	// Read environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "tigertail")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	
	// Construct DB DSN from individual environment variables
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", 
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
	
	// Read Redis environment variables
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDBStr := getEnv("REDIS_DB", "0")
	redisDB := 0 // Default Redis DB
	fmt.Sscanf(redisDBStr, "%d", &redisDB)
	
	// Construct Redis address from individual environment variables
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	
	// Get server port
	port := getEnv("SERVER_PORT", getEnv("PORT", "8080"))

	// Log the connection details
	log.Printf("Connecting to PostgreSQL with DSN: %s", dbDSN)
	log.Printf("Connecting to Redis at %s (DB: %d)", redisAddr, redisDB)
	
	// Check if we should use real database
	useRealDB := getEnv("USE_REAL_DB", "false") == "true"
	var postgres *db.PostgresDB
	var err error
	
	if useRealDB {
		// Initialize database connection
		postgres, err = db.NewPostgresConnection(dbDSN)
		if err != nil {
			log.Printf("Error: Failed to connect to database: %v", err)
			return "", fmt.Errorf("failed to connect to database: %w", err)
		}
	} else {
		log.Printf("Stub: Would connect to PostgreSQL with DSN: %s", dbDSN)
		// Create a stub implementation
		postgres = db.NewPostgresStub()
	}

	// Check if we should use real Redis
	useRealRedis := getEnv("USE_REAL_REDIS", "false") == "true"
	var redisClient *cache.RedisClient
	
	if useRealRedis {
		// Initialize Redis connection
		redisClient, err = cache.NewRedisClient(redisAddr, redisPassword, redisDB)
		if err != nil {
			log.Printf("Error: Failed to connect to Redis: %v", err)
			return "", fmt.Errorf("failed to connect to Redis: %w", err)
		}
	} else {
		log.Printf("Stub: Would connect to Redis at %s (DB: %d)", redisAddr, redisDB)
		// Create a stub implementation
		redisClient = cache.NewRedisStub()
	}
	
	// Create repositories
	postRepo := db.NewPostRepository(postgres)
	
	// Create cache
	postCache := cache.NewPostCache(redisClient)
	
	// Setup routes with real implementations
	setupRoutes(postRepo, postCache)

	return port, nil
}

// setupRoutes sets up the HTTP routes
func setupRoutes(postRepo *db.PostRepository, postCache *cache.PostCache) {
	// Root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "Tiger-Tail Microblog API"}`))
	})
	
	// API endpoint
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Tiger-Tail Microblog API", "version": "0.1.0"}`))
	})
	
	// Posts endpoint - GET
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
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
				if err == nil && pageInt > 0 {
					page = pageInt
				}
			}

			// Parse limit parameter
			if limitStr != "" {
				limitInt, err := strconv.Atoi(limitStr)
				if err == nil && limitInt > 0 && limitInt <= 100 {
					limit = limitInt
				}
			}

			// Calculate offset
			offset := (page - 1) * limit

			// Try to get posts from cache
			posts, err := postCache.GetPostsWithUser()
			if err == nil {
				// Cache hit
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
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
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"posts":  posts,
				"page":   page,
				"limit":  limit,
				"total":  total,
				"source": "database",
			})
			return
		} else if r.Method == http.MethodPost {
			// Check authentication
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Unauthorized",
				})
				return
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
			if username != expectedUsername || password != expectedPassword {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Unauthorized",
				})
				return
			}

			// Parse request body
			var requestBody struct {
				Content string `json:"content"`
			}
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Invalid request body",
				})
				return
			}

			// Validate content
			if requestBody.Content == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Content is required",
				})
				return
			}

			// Create post
			post := &domain.Post{
				ID:        fmt.Sprintf("post_%d", time.Now().UnixNano()),
				UserID:    "user_1", // Use a fixed user ID for simplicity
				Content:   requestBody.Content,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Save post to database
			err = postRepo.Create(post)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Failed to create post",
				})
				return
			}

			// Invalidate cache
			go postCache.InvalidatePosts()

			// Return success
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"post":    post,
				"message": "Post created successfully",
			})
			return
		} else {
			// Method not allowed
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method not allowed",
			})
			return
		}
	})
	
	// Health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})
	
	// Liveness probe endpoint - separate handler for plain text response
	http.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		// Force content type to text/plain
		w.Header().Set("Content-Type", "text/plain")
		// Write the status code
		w.WriteHeader(http.StatusOK)
		// Write the response body
		w.Write([]byte("OK."))
	})
	
	// Readiness probe endpoint
	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ready", "checks": {"database": "up", "cache": "up"}}`))
	})
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
