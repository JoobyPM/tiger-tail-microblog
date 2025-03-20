package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		wantStatus int
		wantHeader string
	}{
		{
			name:       "Success with data",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusOK,
			wantHeader: "application/json",
		},
		{
			name:       "Error status with data",
			statusCode: http.StatusBadRequest,
			data:       map[string]string{"error": "bad request"},
			wantStatus: http.StatusBadRequest,
			wantHeader: "application/json",
		},
		{
			name:       "Success with nil data",
			statusCode: http.StatusNoContent,
			data:       nil,
			wantStatus: http.StatusNoContent,
			wantHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the function
			WriteJSONResponse(rr, tt.statusCode, tt.data)

			// Check status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			// Check content type
			if contentType := rr.Header().Get("Content-Type"); contentType != tt.wantHeader {
				t.Errorf("handler returned wrong content type: got %v want %v", contentType, tt.wantHeader)
			}

			// If data is not nil, check the response body
			if tt.data != nil {
				var response map[string]string
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse response body: %v", err)
				}

				// Check that the response contains the expected data
				for key, expectedValue := range tt.data.(map[string]string) {
					if actualValue, ok := response[key]; !ok || actualValue != expectedValue {
						t.Errorf("response[%s] = %v, want %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}

// Test with invalid JSON data that would cause an encoding error
func TestWriteJSONResponse_EncodingError(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a struct with a channel which cannot be JSON encoded
	type BadData struct {
		Ch chan int
	}
	badData := BadData{
		Ch: make(chan int),
	}

	// Call the function
	WriteJSONResponse(rr, http.StatusOK, badData)

	// Check status code - should still be 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check that the response contains the error message
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if errorMsg, ok := response["error"]; !ok || errorMsg != "Internal server error" {
		t.Errorf("response[error] = %v, want %v", errorMsg, "Internal server error")
	}
}
