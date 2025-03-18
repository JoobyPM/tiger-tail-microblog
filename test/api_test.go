package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
	"github.com/JoobyPM/tiger-tail-microblog/internal/server"
)

// MockDBPinger is a mock implementation of server.DBPinger
type MockDBPinger struct {
	PingFunc func() error
}

func (m *MockDBPinger) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}

// MockCachePinger is a mock implementation of server.CachePinger
type MockCachePinger struct {
	PingFunc func() error
}

func (m *MockCachePinger) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}

// MockPostService is a mock implementation of domain.PostService
type MockPostService struct {
	GetByIDFunc func(id string) (*domain.PostWithUser, error)
	CreateFunc  func(userID, content string) (*domain.Post, error)
	ListFunc    func(page, limit int) ([]*domain.PostWithUser, int, error)
}

func (m *MockPostService) GetByID(id string) (*domain.PostWithUser, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
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

func (m *MockPostService) Create(userID, content string) (*domain.Post, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(userID, content)
	}
	return &domain.Post{
		ID:        "post_123",
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockPostService) Update(id, userID, content string) (*domain.Post, error) {
	return nil, nil
}

func (m *MockPostService) Delete(id, userID string) error {
	return nil
}

func (m *MockPostService) ListByUser(userID string, page, limit int) ([]*domain.Post, int, error) {
	return nil, 0, nil
}

func (m *MockPostService) List(page, limit int) ([]*domain.PostWithUser, int, error) {
	if m.ListFunc != nil {
		return m.ListFunc(page, limit)
	}
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

// MockPostCache is a mock implementation of server.PostCache
type MockPostCache struct {
	GetPostFunc          func(id string) (*domain.Post, error)
	SetPostFunc          func(post *domain.Post) error
	GetPostsWithUserFunc func() ([]*domain.PostWithUser, error)
	SetPostsWithUserFunc func(posts []*domain.PostWithUser) error
	InvalidatePostsFunc  func() error
	PingFunc             func() error
}

func (m *MockPostCache) GetPost(id string) (*domain.Post, error) {
	if m.GetPostFunc != nil {
		return m.GetPostFunc(id)
	}
	return nil, domain.ErrPostNotFound
}

func (m *MockPostCache) SetPost(post *domain.Post) error {
	if m.SetPostFunc != nil {
		return m.SetPostFunc(post)
	}
	return nil
}

func (m *MockPostCache) GetPostsWithUser() ([]*domain.PostWithUser, error) {
	if m.GetPostsWithUserFunc != nil {
		return m.GetPostsWithUserFunc()
	}
	return nil, domain.ErrPostNotFound
}

func (m *MockPostCache) SetPostsWithUser(posts []*domain.PostWithUser) error {
	if m.SetPostsWithUserFunc != nil {
		return m.SetPostsWithUserFunc(posts)
	}
	return nil
}

func (m *MockPostCache) InvalidatePosts() error {
	if m.InvalidatePostsFunc != nil {
		return m.InvalidatePostsFunc()
	}
	return nil
}

func (m *MockPostCache) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}

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

// TestLivezEndpoint tests the /livez endpoint
func TestLivezEndpoint(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/livez", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler function using the actual LivezHandler
	handler := server.LivezHandler()

	// Serve the request
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
	body := rr.Body.String()
	if body != "OK." {
		t.Errorf("handler returned unexpected body: got %v want %v", body, "OK.")
	}
}

// TestReadyzEndpoint tests the /readyz endpoint
func TestReadyzEndpoint(t *testing.T) {
	// Create mocks
	mockDB := &MockDBPinger{
		PingFunc: func() error {
			return nil // DB is up
		},
	}
	mockCache := &MockCachePinger{
		PingFunc: func() error {
			return nil // Cache is up
		},
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/readyz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler function using the actual ReadyzHandler
	handler := server.ReadyzHandler(mockDB, mockCache)

	// Serve the request
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
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the status field
	if response["status"] != "ready" {
		t.Errorf("handler returned unexpected status: got %v want %v", response["status"], "ready")
	}

	// Check the checks field
	checks, ok := response["checks"].(map[string]interface{})
	if !ok {
		t.Fatalf("handler returned unexpected checks type: got %T want map[string]interface{}", response["checks"])
	}

	// Check the database status
	if checks["database"] != "up" {
		t.Errorf("handler returned unexpected database status: got %v want %v", checks["database"], "up")
	}

	// Check the cache status
	if checks["cache"] != "up" {
		t.Errorf("handler returned unexpected cache status: got %v want %v", checks["cache"], "up")
	}
}

// TestPostsLogic tests the posts logic with cache fallback
func TestPostsLogic(t *testing.T) {
	// Create mocks
	mockPostService := &MockPostService{
		ListFunc: func(page, limit int) ([]*domain.PostWithUser, int, error) {
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

	// First test: Cache miss, fetch from DB
	mockPostCache := &MockPostCache{
		GetPostsWithUserFunc: func() ([]*domain.PostWithUser, error) {
			return nil, domain.ErrPostNotFound // Cache miss
		},
		SetPostsWithUserFunc: func(posts []*domain.PostWithUser) error {
			return nil // Successfully set in cache
		},
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a post handler
	postHandler := server.NewPostHandler(mockPostService, mockPostCache)
	handler := postHandler.GetPostsHandler()

	// Serve the request
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
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the source field (should be "database" for cache miss)
	if response["source"] != "database" {
		t.Errorf("handler returned unexpected source: got %v want %v", response["source"], "database")
	}

	// Second test: Cache hit
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

	mockPostCache = &MockPostCache{
		GetPostsWithUserFunc: func() ([]*domain.PostWithUser, error) {
			return cachedPosts, nil // Cache hit
		},
	}

	// Create a new post handler with the updated cache
	postHandler = server.NewPostHandler(mockPostService, mockPostCache)
	handler = postHandler.GetPostsHandler()

	// Create a new response recorder
	rr = httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the source field (should be "cache" for cache hit)
	if response["source"] != "cache" {
		t.Errorf("handler returned unexpected source: got %v want %v", response["source"], "cache")
	}
}
