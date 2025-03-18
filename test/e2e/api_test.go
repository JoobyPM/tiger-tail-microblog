package e2e

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestConfig holds the test configuration
type TestConfig struct {
	APIURL string
}

// LoadTestConfig loads the test configuration from environment variables
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		APIURL: getEnv("API_URL", "http://localhost:8080"),
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

// TestHealthEndpoint tests the health endpoint
func TestHealthEndpoint(t *testing.T) {
	config := LoadTestConfig()
	url := fmt.Sprintf("%s/health", config.APIURL)

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type %s, got %s", "application/json", contentType)
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the status field
	status, ok := response["status"].(string)
	if !ok {
		t.Fatalf("Expected status field to be a string")
	}
	if status != "ok" {
		t.Errorf("Expected status to be 'ok', got '%s'", status)
	}
}

// TestLivezEndpoint tests the /livez endpoint
func TestLivezEndpoint(t *testing.T) {
	config := LoadTestConfig()
	url := fmt.Sprintf("%s/livez", config.APIURL)

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Expected content type %s, got %s", "text/plain", contentType)
	}

	// Read the response body
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check the response body
	if string(body) != "OK." {
		t.Errorf("Expected body to be 'OK.', got '%s'", string(body))
	}
}

// TestReadyzEndpoint tests the /readyz endpoint
func TestReadyzEndpoint(t *testing.T) {
	config := LoadTestConfig()
	url := fmt.Sprintf("%s/readyz", config.APIURL)

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type %s, got %s", "application/json", contentType)
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the status field
	status, ok := response["status"].(string)
	if !ok {
		t.Fatalf("Expected status field to be a string")
	}
	if status != "ready" {
		t.Errorf("Expected status to be 'ready', got '%s'", status)
	}

	// Check the checks field
	checks, ok := response["checks"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected checks field to be a map")
	}

	// Check the database status
	dbStatus, ok := checks["database"].(string)
	if !ok {
		t.Fatalf("Expected database status to be a string")
	}
	if dbStatus != "up" {
		t.Errorf("Expected database status to be 'up', got '%s'", dbStatus)
	}

	// Check the cache status
	cacheStatus, ok := checks["cache"].(string)
	if !ok {
		t.Fatalf("Expected cache status to be a string")
	}
	if cacheStatus != "up" {
		t.Errorf("Expected cache status to be 'up', got '%s'", cacheStatus)
	}
}

// TestAPIEndpoint tests the API endpoint
func TestAPIEndpoint(t *testing.T) {
	config := LoadTestConfig()
	url := fmt.Sprintf("%s/api", config.APIURL)

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type %s, got %s", "application/json", contentType)
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the message field
	message, ok := response["message"].(string)
	if !ok {
		t.Fatalf("Expected message field to be a string")
	}
	if message != "Tiger-Tail Microblog API" {
		t.Errorf("Expected message to be 'Tiger-Tail Microblog API', got '%s'", message)
	}

	// Check the version field
	version, ok := response["version"].(string)
	if !ok {
		t.Fatalf("Expected version field to be a string")
	}
	if version == "" {
		t.Errorf("Expected version to be non-empty")
	}
}

// TestPostsEndpoint tests the posts endpoint
func TestPostsEndpoint(t *testing.T) {
	config := LoadTestConfig()
	url := fmt.Sprintf("%s/api/posts", config.APIURL)

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type %s, got %s", "application/json", contentType)
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the posts field
	posts, ok := response["posts"].([]interface{})
	if !ok {
		t.Fatalf("Expected posts field to be an array")
	}

	// Check that we have some posts
	if len(posts) == 0 {
		t.Errorf("Expected at least one post")
	}

	// Check the source field (should be "database" or "cache")
	source, ok := response["source"].(string)
	if !ok {
		t.Fatalf("Expected source field to be a string")
	}
	if source != "database" && source != "cache" {
		t.Errorf("Expected source to be 'database' or 'cache', got '%s'", source)
	}
}

// TestMain is the entry point for the test suite
func TestMain(m *testing.M) {
	// Wait for services to be ready
	config := LoadTestConfig()
	waitForService(config.APIURL)

	// Run the tests
	code := m.Run()

	// Exit with the test status code
	os.Exit(code)
}

// waitForService waits for the service to be ready
func waitForService(url string) {
	maxRetries := 30
	retryInterval := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(fmt.Sprintf("%s/livez", url))
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			fmt.Println("Service is ready")
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		fmt.Printf("Waiting for service to be ready (attempt %d/%d)...\n", i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	fmt.Println("Service is not ready after maximum retries")
}
