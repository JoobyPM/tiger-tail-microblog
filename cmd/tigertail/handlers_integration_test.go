package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/cache"
	"github.com/JoobyPM/tiger-tail-microblog/internal/config"
	"github.com/JoobyPM/tiger-tail-microblog/internal/db"
	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// TestHandleGetPostsIntegration tests the handleGetPosts function with integration stubs
func TestHandleGetPostsIntegration(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Initialize stubs
	postgres := db.NewPostgresStub()
	redis := cache.NewRedisStub()
	
	// Create repositories
	postRepo := db.NewPostRepository(postgres)
	postCache := cache.NewPostCache(redis)

	// Call the handler
	handleGetPosts(rr, req, postRepo, postCache)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Logf("Response body: %s", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
		return
	}

	// Check that source is "database" (since cache will be empty)
	if source, ok := response["source"].(string); !ok || source != "database" {
		t.Errorf("handler returned unexpected source: got %v want %v", source, "database")
	}
}

// TestHandleCreatePostIntegration tests the handleCreatePost function with integration stubs
func TestHandleCreatePostIntegration(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		requestBody    string
		auth           bool
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid Request",
			requestBody:    `{"content":"Test post content"}`,
			auth:           true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Unauthorized",
			requestBody:    `{"content":"Test post content"}`,
			auth:           false,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Unauthorized",
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"content":}`,
			auth:           true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "Empty Content",
			requestBody:    `{"content":""}`,
			auth:           true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Content is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("POST", "/api/posts", strings.NewReader(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Add basic auth if needed
			if tc.auth {
				req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:password")))
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Initialize stubs
			postgres := db.NewPostgresStub()
			redis := cache.NewRedisStub()
			
			// Create repositories
			postRepo := db.NewPostRepository(postgres)
			postCache := cache.NewPostCache(redis)
			
			authCreds := &config.AuthCredentials{
				Username: config.SensitiveString("admin"),
				Password: config.SensitiveString("password"),
			}

			// Call the handler
			handleCreatePost(rr, req, postRepo, postCache, authCreds)

			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Logf("Response body: %s", rr.Body.String())
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
				return
			}

			// Parse the response body
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to parse response body: %v", err)
				return
			}

			// Check for expected error message
			if tc.expectedError != "" {
				if errMsg, ok := response["error"].(string); !ok || errMsg != tc.expectedError {
					t.Errorf("handler returned unexpected error: got %v want %v", errMsg, tc.expectedError)
				}
			} else if tc.expectedStatus == http.StatusCreated {
				// Check that post is returned for successful requests
				if post, ok := response["post"].(map[string]interface{}); !ok {
					t.Errorf("handler did not return post")
				} else {
					// Check post has content field
					if _, ok := post["content"]; !ok {
						t.Errorf("post does not have content field")
					}
				}
			}
		})
	}
}

// CustomPostgresStub is a custom stub for PostgresDB that returns predefined data
type CustomPostgresStub struct {
	posts []domain.Post
}

func NewCustomPostgresStub() *db.PostgresDB {
	// Use the standard stub as a base
	stub := db.NewPostgresStub()
	
	// We can't modify the stub directly, so we'll just return it
	return stub
}

func (s *CustomPostgresStub) Close() error {
	return nil
}

func (s *CustomPostgresStub) Ping() error {
	return nil
}

func (s *CustomPostgresStub) Exec(query string, args ...interface{}) error {
	return nil
}

func (s *CustomPostgresStub) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(s.posts))
	for _, post := range s.posts {
		result = append(result, map[string]interface{}{
			"id":         post.ID,
			"user_id":    post.UserID,
			"content":    post.Content,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
			"username":   "testuser",
		})
	}
	return result, nil
}

func (s *CustomPostgresStub) QueryRow(query string, args ...interface{}) (map[string]interface{}, error) {
	if strings.Contains(query, "COUNT") {
		return map[string]interface{}{"count": len(s.posts)}, nil
	}
	return nil, nil
}

// CustomRedisStub is a custom stub for RedisClient that can be configured to return cache hits or misses
type CustomRedisStub struct {
	cacheHit bool
	posts    []domain.Post
}

func NewCustomRedisStub(cacheHit bool) cache.RedisClientInterface {
	// Create a custom Redis stub
	now := time.Now()
	posts := []domain.Post{
		{
			ID:        "post_1",
			UserID:    "user_1",
			Content:   "Test post content 1",
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
		},
		{
			ID:        "post_2",
			UserID:    "user_1",
			Content:   "Test post content 2",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	
	return &CustomRedisStub{
		cacheHit: cacheHit,
		posts:    posts,
	}
}

func (s *CustomRedisStub) Get(key string) ([]byte, error) {
	if !s.cacheHit {
		return nil, cache.ErrCacheMiss
	}

	if key == "posts_with_user" {
		postsWithUser := make([]map[string]interface{}, 0, len(s.posts))
		for _, post := range s.posts {
			postsWithUser = append(postsWithUser, map[string]interface{}{
				"id":         post.ID,
				"user_id":    post.UserID,
				"content":    post.Content,
				"created_at": post.CreatedAt,
				"updated_at": post.UpdatedAt,
				"username":   "testuser",
			})
		}
		data, _ := json.Marshal(postsWithUser)
		return data, nil
	}

	return nil, cache.ErrCacheMiss
}

func (s *CustomRedisStub) Set(key string, value []byte, expiration time.Duration) error {
	return nil
}

func (s *CustomRedisStub) Delete(key string) error {
	return nil
}

func (s *CustomRedisStub) Exists(key string) (bool, error) {
	return s.cacheHit, nil
}

func (s *CustomRedisStub) Ping() error {
	return nil
}

func (s *CustomRedisStub) Close() error {
	return nil
}

func (s *CustomRedisStub) FlushDB() error {
	return nil
}

// TestHandleGetPostsWithCustomStubs tests the handleGetPosts function with custom stubs
func TestHandleGetPostsWithCustomStubs(t *testing.T) {
	// Test cases for cache hit and miss
	tests := []struct {
		name           string
		cacheHit       bool
		expectedSource string
	}{
		{
			name:           "Cache Hit",
			cacheHit:       true,
			expectedSource: "cache",
		},
		{
			name:           "Cache Miss",
			cacheHit:       false,
			expectedSource: "database",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/api/posts", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Initialize custom stubs
			postgres := NewCustomPostgresStub()
			redis := NewCustomRedisStub(tc.cacheHit)
			
			// Create repositories
			postRepo := db.NewPostRepository(postgres)
			postCache := cache.NewPostCache(redis)

			// Call the handler
			handleGetPosts(rr, req, postRepo, postCache)

			// Check the status code
			if status := rr.Code; status != http.StatusOK {
				t.Logf("Response body: %s", rr.Body.String())
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
				return
			}

			// Check the response body
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to parse response body: %v", err)
				return
			}

			// Check that source matches expected
			if source, ok := response["source"].(string); !ok || source != tc.expectedSource {
				t.Errorf("handler returned unexpected source: got %v want %v", source, tc.expectedSource)
			}
		})
	}
}
