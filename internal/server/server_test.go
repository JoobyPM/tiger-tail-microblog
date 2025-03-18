package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Test cases
	testCases := []struct {
		name   string
		config Config
	}{
		{
			name: "default config",
			config: Config{
				Host:    "localhost",
				Port:    8080,
				BaseURL: "http://localhost:8080",
			},
		},
		{
			name: "custom config",
			config: Config{
				Host:    "127.0.0.1",
				Port:    9090,
				BaseURL: "http://example.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			server := New(tc.config)

			// Assert
			if server == nil {
				t.Fatal("New returned nil")
			}
			if server.config != tc.config {
				t.Errorf("config = %v, want %v", server.config, tc.config)
			}
			if server.router == nil {
				t.Error("router is nil")
			}
			if server.httpServer == nil {
				t.Error("httpServer is nil")
			}
			if server.httpServer.Addr != fmt.Sprintf("%s:%d", tc.config.Host, tc.config.Port) {
				t.Errorf("httpServer.Addr = %s, want %s", server.httpServer.Addr, fmt.Sprintf("%s:%d", tc.config.Host, tc.config.Port))
			}
			if server.httpServer.ReadTimeout != 15*time.Second {
				t.Errorf("httpServer.ReadTimeout = %v, want %v", server.httpServer.ReadTimeout, 15*time.Second)
			}
			if server.httpServer.WriteTimeout != 15*time.Second {
				t.Errorf("httpServer.WriteTimeout = %v, want %v", server.httpServer.WriteTimeout, 15*time.Second)
			}
			if server.httpServer.IdleTimeout != 60*time.Second {
				t.Errorf("httpServer.IdleTimeout = %v, want %v", server.httpServer.IdleTimeout, 60*time.Second)
			}
		})
	}
}

func TestRegisterRoutes(t *testing.T) {
	// Setup
	server := New(Config{
		Host:    "localhost",
		Port:    8080,
		BaseURL: "http://localhost:8080",
	})

	// Test
	server.registerRoutes()

	// Assert - we can't directly test the routes, but we can test that the handlers work
	// by making requests to them
	testServer := httptest.NewServer(server.router)
	defer testServer.Close()

	// Test health endpoint
	resp, err := http.Get(testServer.URL + "/health")
	if err != nil {
		t.Fatalf("Error making request to /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s, want %s", contentType, "application/json")
	}

	var healthResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if healthResponse["status"] != "ok" {
		t.Errorf("status = %s, want %s", healthResponse["status"], "ok")
	}

	// Test API endpoint
	resp, err = http.Get(testServer.URL + "/api/")
	if err != nil {
		t.Fatalf("Error making request to /api/: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	contentType = resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s, want %s", contentType, "application/json")
	}

	var apiResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if apiResponse["message"] != "Tiger-Tail Microblog API" {
		t.Errorf("message = %s, want %s", apiResponse["message"], "Tiger-Tail Microblog API")
	}
	if apiResponse["version"] != "0.1.0" {
		t.Errorf("version = %s, want %s", apiResponse["version"], "0.1.0")
	}
}

func TestHandleHealth(t *testing.T) {
	// Setup
	server := New(Config{
		Host:    "localhost",
		Port:    8080,
		BaseURL: "http://localhost:8080",
	})
	handler := server.handleHealth()

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("status = %s, want %s", response["status"], "ok")
	}

	// Check that time is a valid RFC3339 timestamp
	_, err = time.Parse(time.RFC3339, response["time"])
	if err != nil {
		t.Errorf("time is not a valid RFC3339 timestamp: %s", response["time"])
	}
}

func TestHandleAPI(t *testing.T) {
	// Setup
	server := New(Config{
		Host:    "localhost",
		Port:    8080,
		BaseURL: "http://localhost:8080",
	})
	handler := server.handleAPI()

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/api/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response["message"] != "Tiger-Tail Microblog API" {
		t.Errorf("message = %s, want %s", response["message"], "Tiger-Tail Microblog API")
	}
	if response["version"] != "0.1.0" {
		t.Errorf("version = %s, want %s", response["version"], "0.1.0")
	}
}

func TestRespondJSON(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Test data
	data := map[string]string{
		"key": "value",
	}

	// Call the function
	respondJSON(rr, http.StatusOK, data)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("respondJSON returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("respondJSON returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response["key"] != "value" {
		t.Errorf("key = %s, want %s", response["key"], "value")
	}

	// Test with nil data
	rr = httptest.NewRecorder()
	respondJSON(rr, http.StatusNoContent, nil)

	// Check the status code
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("respondJSON returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Check the content type
	contentType = rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("respondJSON returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check the response body is empty
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if len(body) != 0 {
		t.Errorf("body = %s, want empty", body)
	}
}

func TestRespondError(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function
	respondError(rr, http.StatusBadRequest, "Invalid request")

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("respondError returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("respondError returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response["error"] != http.StatusText(http.StatusBadRequest) {
		t.Errorf("error = %s, want %s", response["error"], http.StatusText(http.StatusBadRequest))
	}
	if response["message"] != "Invalid request" {
		t.Errorf("message = %s, want %s", response["message"], "Invalid request")
	}
}

func TestStartAndStop(t *testing.T) {
	// This is a more complex test that starts and stops the server
	// We'll use a goroutine to start the server and then stop it after a short delay

	// Setup
	server := New(Config{
		Host:    "localhost",
		Port:    0, // Use port 0 to let the OS choose a free port
		BaseURL: "http://localhost",
	})

	// Create a custom HTTP server with a random port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Error creating listener: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close() // Close the listener, we just needed it to get a free port

	server.httpServer.Addr = fmt.Sprintf("localhost:%d", port)

	// Start the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	// Wait a moment for the server to start
	time.Sleep(100 * time.Millisecond)

	// Make a request to the server to verify it's running
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
	if err != nil {
		t.Fatalf("Error making request to server: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Stop the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		t.Fatalf("Error stopping server: %v", err)
	}

	// Check that the server has stopped
	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Unexpected error from server: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for server to stop")
	}

	// Verify the server is no longer accepting connections
	_, err = http.Get(fmt.Sprintf("http://localhost:%d/health", port))
	if err == nil {
		t.Error("Expected error making request to stopped server, got nil")
	}
}
