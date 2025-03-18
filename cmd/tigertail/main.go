package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize configuration
	fmt.Println("Initializing Tiger-Tail Microblog...")

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
		port := "8080"
		fmt.Printf("Starting server on port %s...\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigChan
	fmt.Println("\nShutting down Tiger-Tail Microblog...")
}
