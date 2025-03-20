package cache

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// TestPostCache_SetPostsWithUserAndTotal tests the SetPostsWithUserAndTotal method
func TestPostCache_SetPostsWithUserAndTotal(t *testing.T) {
	// Create a mock Redis client
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Create test posts
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
	
	// Set posts with total
	total := 10 // Total is different from len(posts) to test that it's stored correctly
	err := cache.SetPostsWithUserAndTotal(posts, total)
	if err != nil {
		t.Fatalf("Error setting posts with user and total: %v", err)
	}
	
	// Verify that the posts were stored in Redis
	data, err := client.Get("posts_with_user")
	if err != nil {
		t.Fatalf("Error getting posts from Redis: %v", err)
	}
	
	// Unmarshal the data
	var postsWithTotal PostsWithTotal
	err = json.Unmarshal(data, &postsWithTotal)
	if err != nil {
		t.Fatalf("Error unmarshaling posts with total: %v", err)
	}
	
	// Check that the posts and total were stored correctly
	if len(postsWithTotal.Posts) != len(posts) {
		t.Errorf("len(postsWithTotal.Posts) = %d, want %d", len(postsWithTotal.Posts), len(posts))
	}
	
	if postsWithTotal.Total != total {
		t.Errorf("postsWithTotal.Total = %d, want %d", postsWithTotal.Total, total)
	}
	
	// Check that the posts were stored correctly
	for i, post := range postsWithTotal.Posts {
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

// TestPostCache_GetPostsWithUser_BackwardCompatibility tests backward compatibility
// with the old format of posts_with_user (without total)
func TestPostCache_GetPostsWithUser_BackwardCompatibility(t *testing.T) {
	// Create a mock Redis client
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Create test posts
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
	
	// Marshal posts directly (old format)
	data, err := json.Marshal(posts)
	if err != nil {
		t.Fatalf("Error marshaling posts: %v", err)
	}
	
	// Set posts in Redis directly
	err = client.Set("posts_with_user", data, 5*time.Minute)
	if err != nil {
		t.Fatalf("Error setting posts in Redis: %v", err)
	}
	
	// Get posts with user
	cachedPosts, total, err := cache.GetPostsWithUser()
	if err != nil {
		t.Fatalf("Error getting posts from cache: %v", err)
	}
	
	// Check that the posts were retrieved correctly
	if len(cachedPosts) != len(posts) {
		t.Errorf("len(cachedPosts) = %d, want %d", len(cachedPosts), len(posts))
	}
	
	// Check that the total is the same as the number of posts (backward compatibility)
	if total != len(posts) {
		t.Errorf("total = %d, want %d", total, len(posts))
	}
	
	// Check that the posts were retrieved correctly
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

// TestPostCache_Ping tests the Ping method
func TestPostCache_Ping(t *testing.T) {
	// Create a mock Redis client
	client := NewMockRedisClient()
	cache := NewPostCache(client)
	
	// Test Ping
	err := cache.Ping()
	if err != nil {
		t.Errorf("cache.Ping() error = %v, want nil", err)
	}
}

// TestPostsWithTotal tests the PostsWithTotal struct
func TestPostsWithTotal(t *testing.T) {
	// Create test posts
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
	}
	
	// Create PostsWithTotal
	postsWithTotal := PostsWithTotal{
		Posts: posts,
		Total: 10,
	}
	
	// Check that the posts and total were set correctly
	if len(postsWithTotal.Posts) != len(posts) {
		t.Errorf("len(postsWithTotal.Posts) = %d, want %d", len(postsWithTotal.Posts), len(posts))
	}
	
	if postsWithTotal.Total != 10 {
		t.Errorf("postsWithTotal.Total = %d, want %d", postsWithTotal.Total, 10)
	}
	
	// Marshal and unmarshal to test JSON serialization
	data, err := json.Marshal(postsWithTotal)
	if err != nil {
		t.Fatalf("Error marshaling posts with total: %v", err)
	}
	
	var unmarshaledPostsWithTotal PostsWithTotal
	err = json.Unmarshal(data, &unmarshaledPostsWithTotal)
	if err != nil {
		t.Fatalf("Error unmarshaling posts with total: %v", err)
	}
	
	// Check that the posts and total were unmarshaled correctly
	if len(unmarshaledPostsWithTotal.Posts) != len(posts) {
		t.Errorf("len(unmarshaledPostsWithTotal.Posts) = %d, want %d", len(unmarshaledPostsWithTotal.Posts), len(posts))
	}
	
	if unmarshaledPostsWithTotal.Total != 10 {
		t.Errorf("unmarshaledPostsWithTotal.Total = %d, want %d", unmarshaledPostsWithTotal.Total, 10)
	}
}
