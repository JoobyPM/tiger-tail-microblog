package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// HandlerMockPostService is a mock implementation of domain.PostService for testing handlers
type HandlerMockPostService struct{}

func (m *HandlerMockPostService) GetByID(id string) (*domain.PostWithUser, error) {
	if id == "nonexistent" {
		return nil, domain.ErrPostNotFound
	}
	return &domain.PostWithUser{
		Post: domain.Post{
			ID:        id,
			UserID:    "user_1",
			Content:   "Test post content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
	}, nil
}

func (m *HandlerMockPostService) Create(userID, content string) (*domain.Post, error) {
	return &domain.Post{
		ID:        "post_123",
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *HandlerMockPostService) Update(id, userID, content string) (*domain.Post, error) {
	return &domain.Post{
		ID:        id,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *HandlerMockPostService) Delete(id, userID string) error {
	return nil
}

func (m *HandlerMockPostService) ListByUser(userID string, page, limit int) ([]*domain.Post, int, error) {
	posts := []*domain.Post{
		{
			ID:        "post_1",
			UserID:    userID,
			Content:   "Test post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "post_2",
			UserID:    userID,
			Content:   "Test post 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return posts, len(posts), nil
}

func (m *HandlerMockPostService) List(page, limit int) ([]*domain.PostWithUser, int, error) {
	posts := []*domain.PostWithUser{
		{
			Post: domain.Post{
				ID:        "post_1",
				UserID:    "user_1",
				Content:   "Test post 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username: "testuser1",
		},
		{
			Post: domain.Post{
				ID:        "post_2",
				UserID:    "user_2",
				Content:   "Test post 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username: "testuser2",
		},
	}
	return posts, len(posts), nil
}

// HandlerMockPostCache is a mock implementation of PostCache for testing handlers
type HandlerMockPostCache struct {
	posts       map[string]*domain.Post
	postsWithUser []*domain.PostWithUser
}

func NewHandlerMockPostCache() *HandlerMockPostCache {
	return &HandlerMockPostCache{
		posts:       make(map[string]*domain.Post),
		postsWithUser: nil,
	}
}

func (m *HandlerMockPostCache) GetPost(id string) (*domain.Post, error) {
	if post, ok := m.posts[id]; ok {
		return post, nil
	}
	return nil, domain.ErrPostNotFound
}

func (m *HandlerMockPostCache) SetPost(post *domain.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *HandlerMockPostCache) GetPostsWithUser() ([]*domain.PostWithUser, error) {
	if m.postsWithUser == nil {
		return nil, domain.ErrPostNotFound
	}
	return m.postsWithUser, nil
}

func (m *HandlerMockPostCache) SetPostsWithUser(posts []*domain.PostWithUser) error {
	m.postsWithUser = posts
	return nil
}

func (m *HandlerMockPostCache) InvalidatePosts() error {
	m.postsWithUser = nil
	return nil
}

func TestNewPostHandler(t *testing.T) {
	postService := &HandlerMockPostService{}
	postCache := NewHandlerMockPostCache()
	
	handler := NewPostHandler(postService, postCache)
	
	if handler == nil {
		t.Fatal("NewPostHandler returned nil")
	}
	
	if handler.postService != postService {
		t.Errorf("handler.postService = %v, want %v", handler.postService, postService)
	}
	
	if handler.postCache != postCache {
		t.Errorf("handler.postCache = %v, want %v", handler.postCache, postCache)
	}
}

func TestGetPostsHandler(t *testing.T) {
	postService := &HandlerMockPostService{}
	postCache := NewHandlerMockPostCache()
	
	handler := NewPostHandler(postService, postCache)
	
	// Test with GET method
	req, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetPostsHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check the response body
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	posts, ok := response["posts"].([]interface{})
	if !ok {
		t.Fatalf("Expected posts to be an array, got %T", response["posts"])
	}
	
	if len(posts) != 2 {
		t.Errorf("len(posts) = %d, want %d", len(posts), 2)
	}
	
	// Test with non-GET method
	req, err = http.NewRequest("POST", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr = httptest.NewRecorder()
	
	// Call the handler
	handler.GetPostsHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestGetPostHandler(t *testing.T) {
	postService := &HandlerMockPostService{}
	postCache := NewHandlerMockPostCache()
	
	handler := NewPostHandler(postService, postCache)
	
	// Test with GET method and valid post ID
	req, err := http.NewRequest("GET", "/api/posts/post_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.GetPostHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check the response body
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	_, ok := response["post"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected post to be an object, got %T", response["post"])
	}
	
	// Test with non-GET method
	req, err = http.NewRequest("POST", "/api/posts/post_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr = httptest.NewRecorder()
	
	// Call the handler
	handler.GetPostHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestCreatePostHandler(t *testing.T) {
	postService := &HandlerMockPostService{}
	postCache := NewHandlerMockPostCache()
	
	handler := NewPostHandler(postService, postCache)
	
	// Test with POST method and valid request body
	requestBody := map[string]string{
		"content": "Test post content",
	}
	
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	
	req, err := http.NewRequest("POST", "/api/posts/create", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set Basic Auth
	req.SetBasicAuth("admin", "password")
	
	rr := httptest.NewRecorder()
	
	// Call the handler
	handler.CreatePostHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	
	// Check the response body
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	
	_, ok := response["post"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected post to be an object, got %T", response["post"])
	}
	
	// Test with non-POST method
	req, err = http.NewRequest("GET", "/api/posts/create", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr = httptest.NewRecorder()
	
	// Call the handler
	handler.CreatePostHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
	
	// Test with missing authentication
	req, err = http.NewRequest("POST", "/api/posts/create", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	
	rr = httptest.NewRecorder()
	
	// Call the handler
	handler.CreatePostHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
	
	// Test with invalid request body
	req, err = http.NewRequest("POST", "/api/posts/create", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set Basic Auth
	req.SetBasicAuth("admin", "password")
	
	rr = httptest.NewRecorder()
	
	// Call the handler
	handler.CreatePostHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	// Test with empty content
	requestBody = map[string]string{
		"content": "",
	}
	
	body, err = json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	
	req, err = http.NewRequest("POST", "/api/posts/create", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	
	// Set Basic Auth
	req.SetBasicAuth("admin", "password")
	
	rr = httptest.NewRecorder()
	
	// Call the handler
	handler.CreatePostHandler()(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestAuthenticateRequest(t *testing.T) {
	// Test with valid credentials
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	req.SetBasicAuth("admin", "password")
	
	userID, err := authenticateRequest(req)
	if err != nil {
		t.Fatalf("Error authenticating request: %v", err)
	}
	
	if userID != "user_1" {
		t.Errorf("userID = %s, want %s", userID, "user_1")
	}
	
	// Test with invalid credentials
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	req.SetBasicAuth("admin", "wrong_password")
	
	_, err = authenticateRequest(req)
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
	
	// Test with missing credentials
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	_, err = authenticateRequest(req)
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}
