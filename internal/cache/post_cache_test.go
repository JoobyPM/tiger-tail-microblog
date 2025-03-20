package cache

import (
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// MockRedisClient is a mock implementation of the Redis client for testing
type MockRedisClient struct {
	data map[string][]byte
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string][]byte),
	}
}

func (m *MockRedisClient) Get(key string) ([]byte, error) {
	if value, ok := m.data[key]; ok {
		return value, nil
	}
	return nil, ErrCacheMiss
}

func (m *MockRedisClient) Set(key string, value []byte, expiration time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *MockRedisClient) Delete(key string) error {
	delete(m.data, key)
	return nil
}

func (m *MockRedisClient) Exists(key string) (bool, error) {
	_, ok := m.data[key]
	return ok, nil
}

func (m *MockRedisClient) Ping() error {
	return nil
}

func (m *MockRedisClient) Close() error {
	return nil
}

func (m *MockRedisClient) FlushDB() error {
	m.data = make(map[string][]byte)
	return nil
}

func TestNewPostCache(t *testing.T) {
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	if cache == nil {
		t.Fatal("NewPostCache returned nil")
	}
	
	if cache.client != client {
		t.Errorf("cache.client = %v, want %v", cache.client, client)
	}
}

func TestPostCache_GetPost(t *testing.T) {
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Test cache miss
	_, err := cache.GetPost("post_1")
	if err == nil {
		t.Error("Expected error for cache miss, got nil")
	}
	
	// Set a post in the cache
	post := &domain.Post{
		ID:        "post_1",
		UserID:    "user_1",
		Content:   "Test post content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err = cache.SetPost(post)
	if err != nil {
		t.Fatalf("Error setting post in cache: %v", err)
	}
	
	// Test cache hit
	cachedPost, err := cache.GetPost("post_1")
	if err != nil {
		t.Fatalf("Error getting post from cache: %v", err)
	}
	
	if cachedPost.ID != post.ID {
		t.Errorf("cachedPost.ID = %s, want %s", cachedPost.ID, post.ID)
	}
	if cachedPost.UserID != post.UserID {
		t.Errorf("cachedPost.UserID = %s, want %s", cachedPost.UserID, post.UserID)
	}
	if cachedPost.Content != post.Content {
		t.Errorf("cachedPost.Content = %s, want %s", cachedPost.Content, post.Content)
	}
}

func TestPostCache_GetPostsWithUser(t *testing.T) {
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Test cache miss
	_, _, err := cache.GetPostsWithUser()
	if err == nil {
		t.Error("Expected error for cache miss, got nil")
	}
	
	// Set posts in the cache
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
	
	err = cache.SetPostsWithUser(posts)
	if err != nil {
		t.Fatalf("Error setting posts in cache: %v", err)
	}
	
	// Test cache hit
	cachedPosts, total, err := cache.GetPostsWithUser()
	if err != nil {
		t.Fatalf("Error getting posts from cache: %v", err)
	}
	
	if len(cachedPosts) != len(posts) {
		t.Errorf("len(cachedPosts) = %d, want %d", len(cachedPosts), len(posts))
	}
	
	if total != len(posts) {
		t.Errorf("total = %d, want %d", total, len(posts))
	}
	
	for i, post := range cachedPosts {
		if post.ID != posts[i].ID {
			t.Errorf("post.ID = %s, want %s", post.ID, posts[i].ID)
		}
		if post.UserID != posts[i].UserID {
			t.Errorf("post.UserID = %s, want %s", post.UserID, posts[i].UserID)
		}
		if post.Content != posts[i].Content {
			t.Errorf("post.Content = %s, want %s", post.Content, posts[i].Content)
		}
		if post.Username != posts[i].Username {
			t.Errorf("post.Username = %s, want %s", post.Username, posts[i].Username)
		}
	}
}

func TestPostCache_GetPosts(t *testing.T) {
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Test cache miss
	_, err := cache.GetPosts()
	if err == nil {
		t.Error("Expected error for cache miss, got nil")
	}
	
	// Set posts in the cache
	posts := []*domain.Post{
		{
			ID:        "post_1",
			UserID:    "user_1",
			Content:   "Test post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "post_2",
			UserID:    "user_2",
			Content:   "Test post 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	
	err = cache.SetPosts(posts)
	if err != nil {
		t.Fatalf("Error setting posts in cache: %v", err)
	}
	
	// Test cache hit
	cachedPosts, err := cache.GetPosts()
	if err != nil {
		t.Fatalf("Error getting posts from cache: %v", err)
	}
	
	if len(cachedPosts) != len(posts) {
		t.Errorf("len(cachedPosts) = %d, want %d", len(cachedPosts), len(posts))
	}
	
	for i, post := range cachedPosts {
		if post.ID != posts[i].ID {
			t.Errorf("post.ID = %s, want %s", post.ID, posts[i].ID)
		}
		if post.UserID != posts[i].UserID {
			t.Errorf("post.UserID = %s, want %s", post.UserID, posts[i].UserID)
		}
		if post.Content != posts[i].Content {
			t.Errorf("post.Content = %s, want %s", post.Content, posts[i].Content)
		}
	}
}

func TestPostCache_InvalidatePosts(t *testing.T) {
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Set posts in the cache
	posts := []*domain.Post{
		{
			ID:        "post_1",
			UserID:    "user_1",
			Content:   "Test post 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	
	err := cache.SetPosts(posts)
	if err != nil {
		t.Fatalf("Error setting posts in cache: %v", err)
	}
	
	postsWithUser := []*domain.PostWithUser{
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
	}
	
	err = cache.SetPostsWithUser(postsWithUser)
	if err != nil {
		t.Fatalf("Error setting posts with user in cache: %v", err)
	}
	
	// Invalidate posts
	err = cache.InvalidatePosts()
	if err != nil {
		t.Fatalf("Error invalidating posts: %v", err)
	}
	
	// Verify posts are invalidated
	_, err = cache.GetPosts()
	if err == nil {
		t.Error("Expected error for cache miss after invalidation, got nil")
	}
	
	_, _, err = cache.GetPostsWithUser()
	if err == nil {
		t.Error("Expected error for cache miss after invalidation, got nil")
	}
}

func TestPostCache_InvalidatePost(t *testing.T) {
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Set a post in the cache
	post := &domain.Post{
		ID:        "post_1",
		UserID:    "user_1",
		Content:   "Test post content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := cache.SetPost(post)
	if err != nil {
		t.Fatalf("Error setting post in cache: %v", err)
	}
	
	// Invalidate the post
	err = cache.InvalidatePost("post_1")
	if err != nil {
		t.Fatalf("Error invalidating post: %v", err)
	}
	
	// Verify post is invalidated
	_, err = cache.GetPost("post_1")
	if err == nil {
		t.Error("Expected error for cache miss after invalidation, got nil")
	}
}
