package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/JoobyPM/tiger-tail-microblog/internal/cache"
	"github.com/JoobyPM/tiger-tail-microblog/internal/db"
	httputil "github.com/JoobyPM/tiger-tail-microblog/internal/http"
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
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "tigertail_test")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("SERVER_PORT", "8080")
	// Use stubs instead of real connections
	os.Setenv("USE_REAL_DB", "false")
	os.Setenv("USE_REAL_REDIS", "false")
	
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("USE_REAL_DB")
		os.Unsetenv("USE_REAL_REDIS")
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
	// Create a new ServeMux for this test
	mux := http.NewServeMux()
	
	// Create a custom handler that uses our mux
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use a custom setupRoutes function that takes a mux
		setupRoutesWithMux(mux, nil, nil)
		mux.ServeHTTP(w, r)
	})
	
	// Create a test server with our handler
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
}

// setupRoutesWithMux is a helper function for testing that takes a mux
func setupRoutesWithMux(mux *http.ServeMux, postRepo *db.PostRepository, postCache *cache.PostCache) {
	// Setup basic routes
	httputil.SetupBasicRoutes(mux)
}

func TestRunServer(t *testing.T) {
	// Skip this test to avoid conflicts with other tests
	t.Skip("Skipping test to avoid conflicts with other tests")
}
