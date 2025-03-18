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

func main() {
	fmt.Println("Starting TigerTail...")

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

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Setup HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "Tiger-Tail Microblog API"}`))
	})

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting server on port %s...\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigChan
	fmt.Println("\nShutting down TigerTail...")
}

// getEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
