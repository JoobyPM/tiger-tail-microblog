package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// TestGetPostsHandler tests the GetPostsHandler method
func TestGetPostsHandler(t *testing.T) {
	testCases := []struct {
		name           string
		cacheHit       bool
		expectedSource string
		expectedStatus int
	}{
		{
			name:           "Cache hit",
			cacheHit:       true,
			expectedSource: "cache",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Cache miss",
			cacheHit:       false,
			expectedSource: "database",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock post service
			mockPostService := &mockPostService{
				listFunc: func(page, limit int) ([]*domain.PostWithUser, int, error) {
					posts := []*domain.PostWithUser{
						{
							Post: domain.Post{
								ID:        "post_1",
								UserID:    "user_1",
								Content:   "Test post 1 from DB",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
							Username: "testuser1",
						},
						{
							Post: domain.Post{
								ID:        "post_2",
								UserID:    "user_2",
								Content:   "Test post 2 from DB",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
							Username: "testuser2",
						},
					}
					return posts, len(posts), nil
				},
			}

			// Create cached posts
			cachedPosts := []*domain.PostWithUser{
				{
					Post: domain.Post{
						ID:        "post_1",
						UserID:    "user_1",
						Content:   "Test post 1 from cache",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Username: "testuser1",
				},
				{
					Post: domain.Post{
						ID:        "post_2",
						UserID:    "user_2",
						Content:   "Test post 2 from cache",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Username: "testuser2",
				},
			}

			// Create mock post cache
			mockPostCache := &mockPostCache{
				getPostsWithUserFunc: func() ([]*domain.PostWithUser, error) {
					if tc.cacheHit {
						return cachedPosts, nil
					}
					return nil, errors.New("cache miss")
				},
				setPostsWithUserFunc: func(posts []*domain.PostWithUser) error {
					return nil
				},
			}

			// Create post handler
			postHandler := NewPostHandler(mockPostService, mockPostCache)

			// Create a request
			req, err := http.NewRequest("GET", "/api/posts", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler := postHandler.GetPostsHandler()
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

			// Check the source field
			if response["source"] != tc.expectedSource {
				t.Errorf("handler returned unexpected source: got %v want %v", response["source"], tc.expectedSource)
			}
		})
	}
}

// TestGetPostHandler tests the GetPostHandler method
func TestGetPostHandler(t *testing.T) {
	testCases := []struct {
		name           string
		postID         string
		cacheHit       bool
		serviceError   error
		expectedStatus int
		expectedSource string
	}{
		{
			name:           "Cache hit",
			postID:         "post_1",
			cacheHit:       true,
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			expectedSource: "cache",
		},
		{
			name:           "Cache miss, service success",
			postID:         "post_1",
			cacheHit:       false,
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			expectedSource: "database",
		},
		{
			name:           "Cache miss, post not found",
			postID:         "nonexistent",
			cacheHit:       false,
			serviceError:   domain.ErrPostNotFound,
			expectedStatus: http.StatusNotFound,
			expectedSource: "",
		},
		{
			name:           "Cache miss, service error",
			postID:         "post_1",
			cacheHit:       false,
			serviceError:   errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedSource: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock post service
			mockPostService := &mockPostService{
				getByIDFunc: func(id string) (*domain.PostWithUser, error) {
					if tc.serviceError != nil {
						return nil, tc.serviceError
					}
					return &domain.PostWithUser{
						Post: domain.Post{
							ID:        id,
							UserID:    "user_1",
							Content:   "Test post content from DB",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Username: "testuser",
					}, nil
				},
			}

			// Create mock post cache
			mockPostCache := &mockPostCache{
				getPostFunc: func(id string) (*domain.Post, error) {
					if tc.cacheHit {
						return &domain.Post{
							ID:        id,
							UserID:    "user_1",
							Content:   "Test post content from cache",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					}
					return nil, errors.New("cache miss")
				},
				setPostFunc: func(post *domain.Post) error {
					return nil
				},
			}

			// Create post handler
			postHandler := NewPostHandler(mockPostService, mockPostCache)

			// Create a request
			req, err := http.NewRequest("GET", "/api/posts/"+tc.postID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler := postHandler.GetPostHandler()
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			// If we expect an error, we don't need to check the response body
			if tc.expectedStatus != http.StatusOK {
				return
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

			// Check the source field
			if response["source"] != tc.expectedSource {
				t.Errorf("handler returned unexpected source: got %v want %v", response["source"], tc.expectedSource)
			}
		})
	}
}

// TestCreatePostHandler tests the CreatePostHandler method
func TestCreatePostHandler(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		auth           bool
		content        string
		serviceError   error
		expectedStatus int
	}{
		{
			name:           "Valid post creation",
			method:         "POST",
			auth:           true,
			content:        "Test post content",
			serviceError:   nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Method not allowed",
			method:         "GET",
			auth:           true,
			content:        "Test post content",
			serviceError:   nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Unauthorized",
			method:         "POST",
			auth:           false,
			content:        "Test post content",
			serviceError:   nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty content",
			method:         "POST",
			auth:           true,
			content:        "",
			serviceError:   nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Service error",
			method:         "POST",
			auth:           true,
			content:        "Test post content",
			serviceError:   errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock post service
			mockPostService := &mockPostService{
				createFunc: func(userID, content string) (*domain.Post, error) {
					if tc.serviceError != nil {
						return nil, tc.serviceError
					}
					return &domain.Post{
						ID:        "post_123",
						UserID:    userID,
						Content:   content,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				},
			}

			// Create mock post cache
			mockPostCache := &mockPostCache{
				invalidatePostsFunc: func() error {
					return nil
				},
			}

			// Create post handler
			postHandler := NewPostHandler(mockPostService, mockPostCache)

			// Create request body
			requestBody := map[string]string{
				"content": tc.content,
			}
			bodyBytes, _ := json.Marshal(requestBody)

			// Create a request
			req, err := http.NewRequest(tc.method, "/api/posts", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatal(err)
			}

			// Add authentication if needed
			if tc.auth {
				req.SetBasicAuth("admin", "password")
			}

			// Set content type
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler := postHandler.CreatePostHandler()
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			// If we expect an error, we don't need to check the response body
			if tc.expectedStatus != http.StatusCreated {
				return
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

			// Check the message field
			if response["message"] != "Post created successfully" {
				t.Errorf("handler returned unexpected message: got %v want %v", response["message"], "Post created successfully")
			}
		})
	}
}

// TestAuthenticateRequest tests the authenticateRequest function
func TestAuthenticateRequest(t *testing.T) {
	testCases := []struct {
		name          string
		username      string
		password      string
		expectedError error
		expectedID    string
	}{
		{
			name:          "Valid credentials",
			username:      "admin",
			password:      "password",
			expectedError: nil,
			expectedID:    "user_1",
		},
		{
			name:          "Invalid credentials",
			username:      "admin",
			password:      "wrong",
			expectedError: domain.ErrUserNotFound,
			expectedID:    "",
		},
		{
			name:          "No credentials",
			username:      "",
			password:      "",
			expectedError: domain.ErrUserNotFound,
			expectedID:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Add authentication if needed
			if tc.username != "" || tc.password != "" {
				req.SetBasicAuth(tc.username, tc.password)
			}

			// Call the function
			userID, err := authenticateRequest(req)

			// Check the error
			if (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError == nil) {
				t.Errorf("authenticateRequest returned unexpected error: got %v want %v", err, tc.expectedError)
			}

			// Check the user ID
			if userID != tc.expectedID {
				t.Errorf("authenticateRequest returned unexpected user ID: got %v want %v", userID, tc.expectedID)
			}
		})
	}
}

// mockPostService is a mock implementation of domain.PostService for testing
type mockPostService struct {
	getByIDFunc func(id string) (*domain.PostWithUser, error)
	createFunc  func(userID, content string) (*domain.Post, error)
	listFunc    func(page, limit int) ([]*domain.PostWithUser, int, error)
}

func (m *mockPostService) GetByID(id string) (*domain.PostWithUser, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockPostService) Create(userID, content string) (*domain.Post, error) {
	if m.createFunc != nil {
		return m.createFunc(userID, content)
	}
	return nil, nil
}

func (m *mockPostService) Update(id, userID, content string) (*domain.Post, error) {
	return nil, nil
}

func (m *mockPostService) Delete(id, userID string) error {
	return nil
}

func (m *mockPostService) ListByUser(userID string, page, limit int) ([]*domain.Post, int, error) {
	return nil, 0, nil
}

func (m *mockPostService) List(page, limit int) ([]*domain.PostWithUser, int, error) {
	if m.listFunc != nil {
		return m.listFunc(page, limit)
	}
	return nil, 0, nil
}

// mockPostCache is a mock implementation of PostCache for testing
type mockPostCache struct {
	getPostFunc          func(id string) (*domain.Post, error)
	setPostFunc          func(post *domain.Post) error
	getPostsWithUserFunc func() ([]*domain.PostWithUser, error)
	setPostsWithUserFunc func(posts []*domain.PostWithUser) error
	invalidatePostsFunc  func() error
}

func (m *mockPostCache) GetPost(id string) (*domain.Post, error) {
	if m.getPostFunc != nil {
		return m.getPostFunc(id)
	}
	return nil, errors.New("cache miss")
}

func (m *mockPostCache) SetPost(post *domain.Post) error {
	if m.setPostFunc != nil {
		return m.setPostFunc(post)
	}
	return nil
}

func (m *mockPostCache) GetPostsWithUser() ([]*domain.PostWithUser, error) {
	if m.getPostsWithUserFunc != nil {
		return m.getPostsWithUserFunc()
	}
	return nil, errors.New("cache miss")
}

func (m *mockPostCache) SetPostsWithUser(posts []*domain.PostWithUser) error {
	if m.setPostsWithUserFunc != nil {
		return m.setPostsWithUserFunc(posts)
	}
	return nil
}

func (m *mockPostCache) InvalidatePosts() error {
	if m.invalidatePostsFunc != nil {
		return m.invalidatePostsFunc()
	}
	return nil
}

func (m *mockPostCache) Ping() error {
	return nil
}
