package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// TestCreateAndReadPosts tests creating and reading posts with authentication
func TestCreateAndReadPosts(t *testing.T) {
	// Get the API URL from environment variable or use default
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	
	// Test root endpoint
	t.Run("Root", func(t *testing.T) {
		resp, err := http.Get(apiURL + "/")
		if err != nil {
			t.Fatalf("Failed to get root endpoint: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		
		var response map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		if response["status"] != "ok" {
			t.Errorf("status = %s, want %s", response["status"], "ok")
		}
		if response["message"] != "Tiger-Tail Microblog API" {
			t.Errorf("message = %s, want %s", response["message"], "Tiger-Tail Microblog API")
		}
	})
	
	// Test livez endpoint
	t.Run("Livez", func(t *testing.T) {
		resp, err := http.Get(apiURL + "/livez")
		if err != nil {
			t.Fatalf("Failed to get livez endpoint: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		
		contentType := resp.Header.Get("Content-Type")
		if contentType != "text/plain" {
			t.Errorf("Content-Type = %s, want %s", contentType, "text/plain")
		}
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		
		if string(body) != "OK." {
			t.Errorf("body = %s, want %s", string(body), "OK.")
		}
	})
	
	// Test posts endpoint - GET
	t.Run("Posts", func(t *testing.T) {
		resp, err := http.Get(apiURL + "/api/posts")
		if err != nil {
			t.Fatalf("Failed to get posts endpoint: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		
		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type = %s, want %s", contentType, "application/json")
		}
		
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		posts, ok := response["posts"].([]interface{})
		if !ok {
			t.Fatalf("posts is not an array")
		}
		
		fmt.Printf("Found %d posts\n", len(posts))
	})
	
	// Test posts endpoint - POST (unauthorized)
	t.Run("CreatePostUnauthorized", func(t *testing.T) {
		// Create request body
		requestBody := map[string]string{
			"content": "Test post content",
		}
		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		
		// Create request
		req, err := http.NewRequest("POST", apiURL + "/api/posts", bytes.NewBuffer(requestJSON))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		
		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		
		// Check response
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
		}
		
		var response map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		if response["error"] != "Unauthorized" {
			t.Errorf("error = %s, want %s", response["error"], "Unauthorized")
		}
	})
	
	// Test posts endpoint - POST (authorized)
	t.Run("CreatePostAuthorized", func(t *testing.T) {
		// Create request body
		requestBody := map[string]string{
			"content": "Test post content created at " + time.Now().Format(time.RFC3339),
		}
		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		
		// Create request
		req, err := http.NewRequest("POST", apiURL + "/api/posts", bytes.NewBuffer(requestJSON))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth("admin", "password") // Use the default credentials
		
		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()
		
		// Check response
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusCreated)
		}
		
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		if response["message"] != "Post created successfully" {
			t.Errorf("message = %s, want %s", response["message"], "Post created successfully")
		}
		
		post, ok := response["post"].(map[string]interface{})
		if !ok {
			t.Fatalf("post is not an object")
		}
		
		if post["content"] != requestBody["content"] {
			t.Errorf("post.content = %s, want %s", post["content"], requestBody["content"])
		}
		
		// Verify the post was created by getting all posts
		resp, err = http.Get(apiURL + "/api/posts")
		if err != nil {
			t.Fatalf("Failed to get posts endpoint: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
		}
		
		var postsResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&postsResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		posts, ok := postsResponse["posts"].([]interface{})
		if !ok {
			t.Fatalf("posts is not an array")
		}
		
		// Check if the post is in the list
		found := false
		for _, p := range posts {
			post, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			
			if post["content"] == requestBody["content"] {
				found = true
				break
			}
		}
		
		if !found {
			t.Errorf("Created post not found in the list of posts")
		}
	})
}
