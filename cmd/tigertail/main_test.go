package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestRunServer tests the runServer function
func TestRunServer(t *testing.T) {
	// Save original HTTP DefaultServeMux
	originalServeMux := http.DefaultServeMux
	defer func() {
		http.DefaultServeMux = originalServeMux
	}()
	
	// Reset DefaultServeMux for this test
	http.DefaultServeMux = http.NewServeMux()
	
	// Run the server
	shutdown, err := runServer()
	if err != nil {
		t.Fatalf("runServer() error = %v", err)
	}
	
	// Make sure shutdown is not nil
	if shutdown == nil {
		t.Fatal("runServer() returned nil shutdown function")
	}
	
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Try to connect to the server on the default port (8080)
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()
	
	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	
	// Check the response body
	expectedBody := `{"status": "ok", "message": "Tiger-Tail Microblog API"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
	
	// Call the shutdown function
	shutdown()
}

// TestInitApp tests the initApp function
func TestInitApp(t *testing.T) {
	// Save original environment variables
	originalDBDSN := os.Getenv("DB_DSN")
	originalRedisAddr := os.Getenv("REDIS_ADDR")
	originalRedisPassword := os.Getenv("REDIS_PASSWORD")
	originalPort := os.Getenv("PORT")
	
	// Restore environment variables after test
	defer func() {
		os.Setenv("DB_DSN", originalDBDSN)
		os.Setenv("REDIS_ADDR", originalRedisAddr)
		os.Setenv("REDIS_PASSWORD", originalRedisPassword)
		os.Setenv("PORT", originalPort)
	}()
	
	// Set test environment variables
	testDBDSN := "postgres://user:password@testdb:5432/testdb?sslmode=disable"
	testRedisAddr := "testredis:6379"
	testRedisPassword := "testpassword"
	testPort := "9090"
	
	os.Setenv("DB_DSN", testDBDSN)
	os.Setenv("REDIS_ADDR", testRedisAddr)
	os.Setenv("REDIS_PASSWORD", testRedisPassword)
	os.Setenv("PORT", testPort)
	
	// Test initApp function
	port, err := initApp()
	if err != nil {
		t.Errorf("initApp() returned error: %v", err)
	}
	
	if port != testPort {
		t.Errorf("initApp() returned port = %q, want %q", port, testPort)
	}
}

// TestStartServer tests the startServer function
func TestStartServer(t *testing.T) {
	// Reset the default serve mux
	http.DefaultServeMux = http.NewServeMux()
	
	// Setup routes
	setupRoutes()
	
	// Create a temporary listener to find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	
	// Get the port
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	
	// Close the listener to free up the port
	listener.Close()
	
	// Start the server
	startServer(port)
	
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Try to connect to the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/", port))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()
	
	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	
	// Check the response body
	expectedBody := `{"status": "ok", "message": "Tiger-Tail Microblog API"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}
}

// TestSetupRoutes tests the setupRoutes function
func TestSetupRoutes(t *testing.T) {
	// Reset the default serve mux
	http.DefaultServeMux = http.NewServeMux()
	
	// Setup routes
	setupRoutes()
	
	// Create a test request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Serve the request
	http.DefaultServeMux.ServeHTTP(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check the content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
	
	// Check the response body
	expectedBody := `{"status": "ok", "message": "Tiger-Tail Microblog API"}`
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", body, expectedBody)
	}
}

func TestGetEnv(t *testing.T) {
	// Test cases
	testCases := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable is set",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "environment variable is not set",
			key:          "TEST_ENV_VAR_NONEXISTENT",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "environment variable is empty",
			key:          "TEST_ENV_VAR_EMPTY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			if tc.envValue != "" {
				os.Setenv(tc.key, tc.envValue)
				defer os.Unsetenv(tc.key)
			}

			// Test
			result := getEnv(tc.key, tc.defaultValue)

			// Assert
			if result != tc.expected {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tc.key, tc.defaultValue, result, tc.expected)
			}
		})
	}
}
