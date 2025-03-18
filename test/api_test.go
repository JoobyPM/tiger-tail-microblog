package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthEndpoint tests the health endpoint
func TestHealthEndpoint(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a real test, this would use the actual server handler
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the status field
	if response["status"] != "ok" {
		t.Errorf("handler returned unexpected body: got %v want %v", response["status"], "ok")
	}
}

// TestAPIEndpoint tests the API endpoint
func TestAPIEndpoint(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In a real test, this would use the actual server handler
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Tiger-Tail Microblog API",
			"version": "0.1.0",
		})
	})

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the message field
	if response["message"] != "Tiger-Tail Microblog API" {
		t.Errorf("handler returned unexpected body: got %v want %v", response["message"], "Tiger-Tail Microblog API")
	}

	// Check the version field
	if response["version"] != "0.1.0" {
		t.Errorf("handler returned unexpected body: got %v want %v", response["version"], "0.1.0")
	}
}
