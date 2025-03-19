package e2e

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestBasicEndpoints(t *testing.T) {
	// Get the API URL from environment variable or use default
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	// Wait for the API to be ready
	fmt.Println("Waiting for API to be ready...")
	for i := 0; i < 10; i++ {
		resp, err := http.Get(apiURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("API is ready!")
			break
		}
		fmt.Println("API not ready yet, retrying...")
		time.Sleep(1 * time.Second)
	}

	// Test the root endpoint
	t.Run("Root", func(t *testing.T) {
		resp, err := http.Get(apiURL + "/")
		if err != nil {
			t.Fatalf("Failed to get root endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected content type %s, got %s", "application/json", contentType)
		}
	})

	// Test the livez endpoint
	t.Run("Livez", func(t *testing.T) {
		resp, err := http.Get(apiURL + "/livez")
		if err != nil {
			t.Fatalf("Failed to get livez endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "text/plain" {
			t.Errorf("Expected content type %s, got %s", "text/plain", contentType)
		}
	})

	// Test the posts endpoint
	t.Run("Posts", func(t *testing.T) {
		resp, err := http.Get(apiURL + "/api/posts")
		if err != nil {
			t.Fatalf("Failed to get posts endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected content type %s, got %s", "application/json", contentType)
		}
	})
}
