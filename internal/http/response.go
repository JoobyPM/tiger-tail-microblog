package http

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSONResponse writes a JSON response with the given status code and data
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// If encoding fails, log the error and write a simple error message
			log.Printf("Error encoding JSON response: %v", err)
			w.Write([]byte(`{"error":"Internal server error"}`))
		}
	}
}
