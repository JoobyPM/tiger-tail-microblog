package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")
	
	value := getEnv("TEST_ENV_VAR", "default_value")
	if value != "test_value" {
		t.Errorf("getEnv() = %s, want %s", value, "test_value")
	}
	
	// Test with non-existing environment variable
	value = getEnv("NON_EXISTING_ENV_VAR", "default_value")
	if value != "default_value" {
		t.Errorf("getEnv() = %s, want %s", value, "default_value")
	}
}

func TestInitApp(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/tigertail_test?sslmode=disable")
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("PORT", "8080")
	defer func() {
		os.Unsetenv("DB_DSN")
		os.Unsetenv("REDIS_ADDR")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("PORT")
	}()
	
	// Test initApp
	port, err := initApp()
	if err != nil {
		t.Fatalf("initApp() error = %v", err)
	}
	
	// Check that port is set correctly
	if port != "8080" {
		t.Errorf("port = %s, want %s", port, "8080")
	}
}

func TestSetupRoutes(t *testing.T) {
	// Setup test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupRoutes()
		http.DefaultServeMux.ServeHTTP(w, r)
	})
	
	server := httptest.NewServer(handler)
	defer server.Close()
	
	// Test root route
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	
	// Reset default serve mux for other tests
	http.DefaultServeMux = http.NewServeMux()
}

func TestRunServer(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/tigertail_test?sslmode=disable")
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("PORT", "8081")
	defer func() {
		os.Unsetenv("DB_DSN")
		os.Unsetenv("REDIS_ADDR")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("PORT")
	}()
	
	// Test runServer
	shutdown, err := runServer()
	if err != nil {
		t.Fatalf("runServer() error = %v", err)
	}
	
	// Call shutdown function
	if shutdown != nil {
		shutdown()
	}
}
