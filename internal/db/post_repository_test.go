package db

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// TestNewPostRepository tests the NewPostRepository function
func TestNewPostRepository(t *testing.T) {
	// Create a stub PostgresDB
	stubDB := NewPostgresStub()
	
	// Call the function
	repo := NewPostRepository(stubDB)
	
	// Check that the repository was created
	if repo == nil {
		t.Fatal("NewPostRepository returned nil")
	}
	
	// Check that the repository has the correct DB
	if repo.db != stubDB {
		t.Errorf("repo.db = %v, want %v", repo.db, stubDB)
	}
}

// TestPostRepository_GetByID tests the GetByID method
func TestPostRepository_GetByID(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		id          string
		wantPost    *domain.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success with stub",
			id:   "post_1",
			wantPost: &domain.Post{
				ID:      "post_1",
				UserID:  "user_1",
				Content: "Test post content",
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			post, err := repo.GetByID(tt.id)
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("GetByID() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			
			// Check post fields
			if post != nil && tt.wantPost != nil {
				if post.ID != tt.wantPost.ID {
					t.Errorf("post.ID = %v, want %v", post.ID, tt.wantPost.ID)
				}
				if post.UserID != tt.wantPost.UserID {
					t.Errorf("post.UserID = %v, want %v", post.UserID, tt.wantPost.UserID)
				}
				if post.Content != tt.wantPost.Content {
					t.Errorf("post.Content = %v, want %v", post.Content, tt.wantPost.Content)
				}
			}
		})
	}
}

// TestPostRepository_Create tests the Create method
func TestPostRepository_Create(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		post        *domain.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success with stub",
			post: &domain.Post{
				ID:        "post_1",
				UserID:    "user_1",
				Content:   "Test post content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			err := repo.Create(tt.post)
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
				t.Errorf("Create() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

// TestPostRepository_Update tests the Update method
func TestPostRepository_Update(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		post        *domain.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Success with stub",
			post: &domain.Post{
				ID:        "post_1",
				UserID:    "user_1",
				Content:   "Updated content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			err := repo.Update(tt.post)
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil {
				if tt.expectedErr == domain.ErrPostNotFound && err != domain.ErrPostNotFound {
					t.Errorf("Update() error = %v, expectedErr %v", err, tt.expectedErr)
				} else if tt.expectedErr != domain.ErrPostNotFound && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
					t.Errorf("Update() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			}
		})
	}
}

// TestPostRepository_Delete tests the Delete method
func TestPostRepository_Delete(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		id          string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "Success with stub",
			id:      "post_1",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			err := repo.Delete(tt.id)
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil {
				if tt.expectedErr == domain.ErrPostNotFound && err != domain.ErrPostNotFound {
					t.Errorf("Delete() error = %v, expectedErr %v", err, tt.expectedErr)
				} else if tt.expectedErr != domain.ErrPostNotFound && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
					t.Errorf("Delete() error = %v, expectedErr %v", err, tt.expectedErr)
				}
			}
		})
	}
}

// TestPostRepository_Count tests the Count method
func TestPostRepository_Count(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		wantCount   int
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "Success with stub",
			wantCount: 16, // Stub returns 16
			wantErr:   false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			count, err := repo.Count()
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
				t.Errorf("Count() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			
			// Check count
			if count != tt.wantCount {
				t.Errorf("Count() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

// TestPostRepository_CountByUser tests the CountByUser method
func TestPostRepository_CountByUser(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		userID      string
		wantCount   int
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "Success with stub",
			userID:    "user_1",
			wantCount: 2, // Stub returns 2
			wantErr:   false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			count, err := repo.CountByUser(tt.userID)
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("CountByUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
				t.Errorf("CountByUser() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			
			// Check count
			if count != tt.wantCount {
				t.Errorf("CountByUser() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

// TestPostRepository_ListByUser tests the ListByUser method
func TestPostRepository_ListByUser(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		userID      string
		offset      int
		limit       int
		wantLen     int
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "Success with stub",
			userID:  "user_1",
			offset:  0,
			limit:   10,
			wantLen: 2, // Stub returns 2 posts
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			posts, err := repo.ListByUser(tt.userID, tt.offset, tt.limit)
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("ListByUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
				t.Errorf("ListByUser() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			
			// Check posts length
			if len(posts) != tt.wantLen {
				t.Errorf("len(posts) = %v, want %v", len(posts), tt.wantLen)
			}
			
			// Check that all posts have the correct user ID
			for _, post := range posts {
				if post.UserID != tt.userID {
					t.Errorf("post.UserID = %v, want %v", post.UserID, tt.userID)
				}
			}
		})
	}
}

// TestPostRepository_List tests the List method
func TestPostRepository_List(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		offset      int
		limit       int
		wantLen     int
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "Success with stub",
			offset:  0,
			limit:   10,
			wantLen: 10, // Stub returns min(limit, totalPosts) posts
			wantErr: false,
		},
		{
			name:    "Success with stub - offset beyond total",
			offset:  20,
			limit:   10,
			wantLen: 0, // No posts when offset > totalPosts
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			posts, err := repo.List(tt.offset, tt.limit)
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
				t.Errorf("List() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			
			// Check posts length
			if len(posts) != tt.wantLen {
				t.Errorf("len(posts) = %v, want %v", len(posts), tt.wantLen)
			}
		})
	}
}

// TestPostRepository_FetchAllPosts tests the FetchAllPosts method
func TestPostRepository_FetchAllPosts(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		wantLen     int
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "Success with stub",
			wantLen: 16, // Stub returns 16 posts
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create repository with stub DB
			stubDB := NewPostgresStub()
			repo := NewPostRepository(stubDB)
			
			// Call the method
			posts, err := repo.FetchAllPosts()
			
			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAllPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Check specific error
			if tt.wantErr && tt.expectedErr != nil && err != nil && !strings.Contains(err.Error(), tt.expectedErr.Error()) {
				t.Errorf("FetchAllPosts() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			
			// Check posts length
			if len(posts) != tt.wantLen {
				t.Errorf("len(posts) = %v, want %v", len(posts), tt.wantLen)
			}
		})
	}
}
