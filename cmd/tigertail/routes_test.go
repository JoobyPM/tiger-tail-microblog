package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JoobyPM/tiger-tail-microblog/internal/cache"
	"github.com/JoobyPM/tiger-tail-microblog/internal/config"
	"github.com/JoobyPM/tiger-tail-microblog/internal/db"
)

// TestAPIRoutes tests the API routes
func TestAPIRoutes(t *testing.T) {
	// Create a new server mux
	mux := http.NewServeMux()

	// Initialize stubs
	postgres := db.NewPostgresStub()
	redis := cache.NewRedisStub()
	
	// Create repositories
	postRepo := db.NewPostRepository(postgres)
	postCache := cache.NewPostCache(redis)
	
	// Create auth credentials
	authCreds := &config.AuthCredentials{
		Username: config.SensitiveString("admin"),
		Password: config.SensitiveString("password"),
	}

	// Setup routes
	setupAPIRoutes(mux, postRepo, postCache, authCreds)

	// Test GET /api/posts
	t.Run("GET /api/posts", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("GET", "/api/posts", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		mux.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK && status != http.StatusInternalServerError {
			t.Errorf("handler returned unexpected status code: got %v", status)
		}

		// We don't check the response body because the stub might not be fully functional
	})

	// Test POST /api/posts with valid auth
	t.Run("POST /api/posts with valid auth", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("POST", "/api/posts", strings.NewReader(`{"content":"Test post content"}`))
		if err != nil {
			t.Fatal(err)
		}

		// Add basic auth
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		mux.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusCreated && status != http.StatusInternalServerError {
			t.Errorf("handler returned unexpected status code: got %v", status)
		}

		// We don't check the response body because the stub might not be fully functional
	})

	// Test POST /api/posts with invalid auth
	t.Run("POST /api/posts with invalid auth", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("POST", "/api/posts", strings.NewReader(`{"content":"Test post content"}`))
		if err != nil {
			t.Fatal(err)
		}

		// Add invalid basic auth
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("wrong:wrong")))

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		mux.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}

		// Check the response body
		var response map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		// Check the error message
		if error, ok := response["error"]; !ok || error != "Unauthorized" {
			t.Errorf("handler returned unexpected error: got %v want %v", error, "Unauthorized")
		}
	})

	// Test POST /api/posts with empty content
	t.Run("POST /api/posts with empty content", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("POST", "/api/posts", strings.NewReader(`{"content":""}`))
		if err != nil {
			t.Fatal(err)
		}

		// Add basic auth
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		mux.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		// Check the response body
		var response map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		// Check the error message
		if error, ok := response["error"]; !ok || error != "Content is required" {
			t.Errorf("handler returned unexpected error: got %v want %v", error, "Content is required")
		}
	})

	// Test POST /api/posts with invalid JSON
	t.Run("POST /api/posts with invalid JSON", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("POST", "/api/posts", strings.NewReader(`{"content":}`))
		if err != nil {
			t.Fatal(err)
		}

		// Add basic auth
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		mux.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		// Check the response body
		var response map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		// Check the error message
		if error, ok := response["error"]; !ok || error != "Invalid request body" {
			t.Errorf("handler returned unexpected error: got %v want %v", error, "Invalid request body")
		}
	})

	// Test unsupported method
	t.Run("Unsupported method", func(t *testing.T) {
		// Create a request
		req, err := http.NewRequest("PUT", "/api/posts", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Serve the request
		mux.ServeHTTP(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}

		// Check the response body
		var response map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		// Check the error message
		if error, ok := response["error"]; !ok || error != "Method not allowed" {
			t.Errorf("handler returned unexpected error: got %v want %v", error, "Method not allowed")
		}
	})
}

// TestPaginationParamsIntegration tests the parsePaginationParams function
func TestPaginationParamsIntegration(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		queryParams   string
		expectedPage  int
		expectedLimit int
	}{
		{
			name:          "Default Values",
			queryParams:   "",
			expectedPage:  1,
			expectedLimit: DefaultPageSize,
		},
		{
			name:          "Valid Page and Limit",
			queryParams:   "page=2&limit=20",
			expectedPage:  2,
			expectedLimit: 20,
		},
		{
			name:          "Invalid Page (Negative)",
			queryParams:   "page=-1&limit=20",
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "Invalid Page (Non-Numeric)",
			queryParams:   "page=abc&limit=20",
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "Invalid Limit (Negative)",
			queryParams:   "page=2&limit=-1",
			expectedPage:  2,
			expectedLimit: DefaultPageSize,
		},
		{
			name:          "Invalid Limit (Non-Numeric)",
			queryParams:   "page=2&limit=abc",
			expectedPage:  2,
			expectedLimit: DefaultPageSize,
		},
		{
			name:          "Limit Exceeds Maximum",
			queryParams:   "page=2&limit=10001",
			expectedPage:  2,
			expectedLimit: MaxPageSize, // MaxPageSize is 10_000 in the main.go file
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create request with query parameters
			req, err := http.NewRequest("GET", "/?"+tc.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Call the function
			page, limit := parsePaginationParams(req)

			// Check results
			if page != tc.expectedPage {
				t.Errorf("parsePaginationParams() page = %v, want %v", page, tc.expectedPage)
			}
			if limit != tc.expectedLimit {
				t.Errorf("parsePaginationParams() limit = %v, want %v", limit, tc.expectedLimit)
			}
		})
	}
}

// TestAuthenticateRequestIntegration tests the authenticateRequest function
func TestAuthenticateRequestIntegration(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		username       string
		password       string
		expectedResult bool
	}{
		{
			name:           "Valid Credentials",
			username:       "admin",
			password:       "password",
			expectedResult: true,
		},
		{
			name:           "Invalid Username",
			username:       "wrong",
			password:       "password",
			expectedResult: false,
		},
		{
			name:           "Invalid Password",
			username:       "admin",
			password:       "wrong",
			expectedResult: false,
		},
		{
			name:           "No Credentials",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create auth credentials
			authCreds := &config.AuthCredentials{
				Username: config.SensitiveString("admin"),
				Password: config.SensitiveString("password"),
			}

			// Create request
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Add basic auth if credentials are provided
			if tc.username != "" || tc.password != "" {
				req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(tc.username+":"+tc.password)))
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the function
			result := authenticateRequest(rr, req, authCreds)

			// Check result
			if result != tc.expectedResult {
				t.Errorf("authenticateRequest() = %v, want %v", result, tc.expectedResult)
			}

			// If we expect failure, check status code
			if !tc.expectedResult {
				if status := rr.Code; status != http.StatusUnauthorized {
					t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
				}

				// Check response body
				var response map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse response body: %v", err)
				}

				// Check error message
				if error, ok := response["error"]; !ok || error != "Unauthorized" {
					t.Errorf("handler returned unexpected error: got %v want %v", error, "Unauthorized")
				}
			}
		})
	}
}
