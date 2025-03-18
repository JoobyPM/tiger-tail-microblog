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
	dbDSN := getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/tigertail?sslmode=disable")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0 // Default Redis DB
	port := getEnv("PORT", "8080")

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "Tiger-Tail Microblog API"}`))
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
