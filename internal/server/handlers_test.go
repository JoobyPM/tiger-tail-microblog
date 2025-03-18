package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLivezHandler tests the LivezHandler function
func TestLivezHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/livez", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := LivezHandler()
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "text/plain")
	}

	// Check the response body
	expected := "OK."
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestReadyzHandler tests the ReadyzHandler function
func TestReadyzHandler(t *testing.T) {
	testCases := []struct {
		name           string
		dbError        bool
		cacheError     bool
		expectedStatus int
		expectedReady  string
		expectedDB     string
		expectedCache  string
	}{
		{
			name:           "All services up",
			dbError:        false,
			cacheError:     false,
			expectedStatus: http.StatusOK,
			expectedReady:  "ready",
			expectedDB:     "up",
			expectedCache:  "up",
		},
		{
			name:           "Database down",
			dbError:        true,
			cacheError:     false,
			expectedStatus: http.StatusServiceUnavailable,
			expectedReady:  "not ready",
			expectedDB:     "down",
			expectedCache:  "up",
		},
		{
			name:           "Cache down",
			dbError:        false,
			cacheError:     true,
			expectedStatus: http.StatusServiceUnavailable,
			expectedReady:  "not ready",
			expectedDB:     "up",
			expectedCache:  "down",
		},
		{
			name:           "All services down",
			dbError:        true,
			cacheError:     true,
			expectedStatus: http.StatusServiceUnavailable,
			expectedReady:  "not ready",
			expectedDB:     "down",
			expectedCache:  "down",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock DB
			mockDB := &mockDBPinger{
				shouldError: tc.dbError,
			}

			// Create mock Cache
			mockCache := &mockCachePinger{
				shouldError: tc.cacheError,
			}

			// Create a request to pass to our handler
			req, err := http.NewRequest("GET", "/readyz", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the handler
			handler := ReadyzHandler(mockDB, mockCache)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			// Check the content type
			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
			}

			// Check the response body
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Error parsing response body: %v", err)
			}

			// Check status field
			if response["status"] != tc.expectedReady {
				t.Errorf("handler returned wrong status: got %v want %v", response["status"], tc.expectedReady)
			}

			// Check checks field
			checks, ok := response["checks"].(map[string]interface{})
			if !ok {
				t.Fatalf("checks field is not a map: %v", response["checks"])
			}

			// Check database status
			if checks["database"] != tc.expectedDB {
				t.Errorf("handler returned wrong database status: got %v want %v", checks["database"], tc.expectedDB)
			}

			// Check cache status
			if checks["cache"] != tc.expectedCache {
				t.Errorf("handler returned wrong cache status: got %v want %v", checks["cache"], tc.expectedCache)
			}
		})
	}
}

// mockDBPinger is a mock implementation of DBPinger for testing
type mockDBPinger struct {
	shouldError bool
}

func (m *mockDBPinger) Ping() error {
	if m.shouldError {
		return errors.New("database connection error")
	}
	return nil
}

// mockCachePinger is a mock implementation of CachePinger for testing
type mockCachePinger struct {
	shouldError bool
}

func (m *mockCachePinger) Ping() error {
	if m.shouldError {
		return errors.New("cache connection error")
	}
	return nil
}
