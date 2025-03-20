package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupBasicRoutes(t *testing.T) {
	// Create a new ServeMux for testing
	mux := http.NewServeMux()
	
	// Setup the routes
	SetupBasicRoutes(mux)
	
	// Test cases
	tests := []struct {
		name           string
		path           string
		wantStatus     int
		wantBodyKey    string
		wantBodyValue  string
	}{
		{
			name:          "Root path",
			path:          "/",
			wantStatus:    http.StatusOK,
			wantBodyKey:   "message",
			wantBodyValue: "Tiger-Tail Microblog API",
		},
		{
			name:          "API path",
			path:          "/api",
			wantStatus:    http.StatusOK,
			wantBodyKey:   "message",
			wantBodyValue: "Tiger-Tail Microblog API",
		},
		{
			name:          "Non-existent path",
			path:          "/nonexistent",
			wantStatus:    http.StatusNotFound,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Serve the request
			mux.ServeHTTP(rr, req)
			
			// Check status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			
			// For successful requests, check the response body
			if tt.wantStatus == http.StatusOK {
				var response map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse response body: %v", err)
				}
				
				// Check that the response contains the expected data
				if value, ok := response[tt.wantBodyKey]; !ok || value != tt.wantBodyValue {
					t.Errorf("response[%s] = %v, want %v", tt.wantBodyKey, value, tt.wantBodyValue)
				}
			}
		})
	}
}

func TestSetupHealthRoutes(t *testing.T) {
	// Create a new ServeMux for testing
	mux := http.NewServeMux()
	
	// Setup the routes
	SetupHealthRoutes(mux)
	
	// Test cases
	tests := []struct {
		name           string
		path           string
		wantStatus     int
		wantBodyKey    string
		wantBodyValue  string
		wantContentType string
	}{
		{
			name:          "Health endpoint",
			path:          "/health",
			wantStatus:    http.StatusOK,
			wantBodyKey:   "status",
			wantBodyValue: "ok",
			wantContentType: "application/json",
		},
		{
			name:          "Liveness probe",
			path:          "/livez",
			wantStatus:    http.StatusOK,
			wantContentType: "text/plain",
		},
		{
			name:          "Readiness probe",
			path:          "/readyz",
			wantStatus:    http.StatusOK,
			wantBodyKey:   "status",
			wantBodyValue: "ready",
			wantContentType: "application/json",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			
			// Create a response recorder
			rr := httptest.NewRecorder()
			
			// Serve the request
			mux.ServeHTTP(rr, req)
			
			// Check status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
			
			// Check content type
			if contentType := rr.Header().Get("Content-Type"); contentType != tt.wantContentType {
				t.Errorf("handler returned wrong content type: got %v want %v", contentType, tt.wantContentType)
			}
			
			// For JSON responses, check the response body
			if tt.wantContentType == "application/json" {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse response body: %v", err)
				}
				
				// Check that the response contains the expected data
				if value, ok := response[tt.wantBodyKey]; !ok || value != tt.wantBodyValue {
					t.Errorf("response[%s] = %v, want %v", tt.wantBodyKey, value, tt.wantBodyValue)
				}
				
				// For readiness probe, check that checks are included
				if tt.path == "/readyz" {
					if checks, ok := response["checks"].(map[string]interface{}); !ok {
						t.Errorf("response does not contain checks")
					} else {
						// Check that database and cache are included
						if _, ok := checks["database"]; !ok {
							t.Errorf("checks does not contain database")
						}
						if _, ok := checks["cache"]; !ok {
							t.Errorf("checks does not contain cache")
						}
					}
				}
			}
			
			// For text/plain responses, check the response body directly
			if tt.wantContentType == "text/plain" {
				if body := rr.Body.String(); body != "OK." {
					t.Errorf("response body = %v, want %v", body, "OK.")
				}
			}
		})
	}
}
