package db

import (
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// MockPostgresDB is a mock implementation of the PostgreSQL database for testing
type MockPostgresDB struct {
	posts map[string]*domain.Post
}

func NewMockPostgresDB() *MockPostgresDB {
	return &MockPostgresDB{
		posts: make(map[string]*domain.Post),
	}
}

func TestNewPostRepository(t *testing.T) {
	mockDB := &PostgresDB{}
	repo := NewPostRepository(mockDB)
	
	if repo == nil {
		t.Fatal("NewPostRepository returned nil")
	}
}

func TestPostRepository_GetByID(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	// Test getting a post that exists
	post, err := repo.GetByID("post_1")
	if err != nil {
		t.Fatalf("Error getting post: %v", err)
	}
	
	if post.ID != "post_1" {
		t.Errorf("post.ID = %s, want %s", post.ID, "post_1")
	}
	
	// Test getting a post that doesn't exist
	_, err = repo.GetByID("nonexistent")
	if err != domain.ErrPostNotFound {
		t.Errorf("Expected ErrPostNotFound, got %v", err)
	}
}

func TestPostRepository_Create(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	post := &domain.Post{
		ID:        "post_test",
		UserID:    "user_1",
		Content:   "Test post content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := repo.Create(post)
	if err != nil {
		t.Fatalf("Error creating post: %v", err)
	}
}

func TestPostRepository_Update(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	post := &domain.Post{
		ID:        "post_test",
		UserID:    "user_1",
		Content:   "Updated content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := repo.Update(post)
	if err != nil {
		t.Fatalf("Error updating post: %v", err)
	}
}

func TestPostRepository_Delete(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	err := repo.Delete("post_test")
	if err != nil {
		t.Fatalf("Error deleting post: %v", err)
	}
}

func TestPostRepository_ListByUser(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	posts, err := repo.ListByUser("user_1", 0, 10)
	if err != nil {
		t.Fatalf("Error listing posts by user: %v", err)
	}
	
	if len(posts) != 2 {
		t.Errorf("len(posts) = %d, want %d", len(posts), 2)
	}
	
	for _, post := range posts {
		if post.UserID != "user_1" {
			t.Errorf("post.UserID = %s, want %s", post.UserID, "user_1")
		}
	}
}

func TestPostRepository_List(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	posts, err := repo.List(0, 10)
	if err != nil {
		t.Fatalf("Error listing posts: %v", err)
	}
	
	if len(posts) != 5 {
		t.Errorf("len(posts) = %d, want %d", len(posts), 5)
	}
}

func TestPostRepository_CountByUser(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	count, err := repo.CountByUser("user_1")
	if err != nil {
		t.Fatalf("Error counting posts by user: %v", err)
	}
	
	if count != 10 {
		t.Errorf("count = %d, want %d", count, 10)
	}
}

func TestPostRepository_Count(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	count, err := repo.Count()
	if err != nil {
		t.Fatalf("Error counting posts: %v", err)
	}
	
	if count != 25 {
		t.Errorf("count = %d, want %d", count, 25)
	}
}

func TestPostRepository_FetchAllPosts(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	posts, err := repo.FetchAllPosts()
	if err != nil {
		t.Fatalf("Error fetching all posts: %v", err)
	}
	
	if len(posts) != 10 {
		t.Errorf("len(posts) = %d, want %d", len(posts), 10)
	}
}

func TestPostRepository_CreatePost(t *testing.T) {
	db := &PostgresDB{}
	repo := NewPostRepository(db)
	
	post := &domain.Post{
		ID:        "post_test",
		UserID:    "user_1",
		Content:   "Test post content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := repo.CreatePost(post)
	if err != nil {
		t.Fatalf("Error creating post: %v", err)
	}
}
