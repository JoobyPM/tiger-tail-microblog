package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/JoobyPM/tiger-tail-microblog/internal/cache"
	"github.com/JoobyPM/tiger-tail-microblog/internal/db"
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
	
	// Initialize database connection (stubbed)
	_, err := db.NewPostgresConnection(dbDSN)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		// Continue execution - this is just a stub
	}

	// Initialize Redis connection (stubbed)
	_, err = cache.NewRedisClient(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		// Continue execution - this is just a stub
	}

	return port, nil
}

// setupRoutes sets up the HTTP routes
func setupRoutes() {
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
	
	// Posts endpoint
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return some mock posts
		w.Write([]byte(`{
			"posts": [
				{
					"id": "post_1",
					"user_id": "user_1",
					"content": "This is a test post",
					"created_at": "2025-03-18T22:00:00Z",
					"updated_at": "2025-03-18T22:00:00Z",
					"username": "testuser1"
				},
				{
					"id": "post_2",
					"user_id": "user_2",
					"content": "This is another test post",
					"created_at": "2025-03-18T22:30:00Z",
					"updated_at": "2025-03-18T22:30:00Z",
					"username": "testuser2"
				}
			],
			"total": 2,
			"source": "database"
		}`))
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

	// Setup routes
	setupRoutes()

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
