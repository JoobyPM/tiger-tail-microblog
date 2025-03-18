package domain

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrPostNotFound      = errors.New("post not found")
	ErrInvalidPostID     = errors.New("invalid post ID")
	ErrInvalidPostContent = errors.New("invalid post content")
)

// Post represents a microblog post
type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostWithUser represents a post with user information
type PostWithUser struct {
	Post
	Username string `json:"username"`
}

// PostRepository defines the interface for post data access
type PostRepository interface {
	// GetByID retrieves a post by ID
	GetByID(id string) (*Post, error)
	
	// Create creates a new post
	Create(post *Post) error
	
	// Update updates an existing post
	Update(post *Post) error
	
	// Delete deletes a post
	Delete(id string) error
	
	// ListByUser retrieves posts by a specific user with pagination
	ListByUser(userID string, offset, limit int) ([]*Post, error)
	
	// List retrieves a list of posts with pagination
	List(offset, limit int) ([]*PostWithUser, error)
	
	// CountByUser returns the total number of posts by a specific user
	CountByUser(userID string) (int, error)
	
	// Count returns the total number of posts
	Count() (int, error)
}

// PostService defines the interface for post business logic
type PostService interface {
	// GetByID retrieves a post by ID
	GetByID(id string) (*PostWithUser, error)
	
	// Create creates a new post
	Create(userID, content string) (*Post, error)
	
	// Update updates an existing post
	Update(id, userID, content string) (*Post, error)
	
	// Delete deletes a post
	Delete(id, userID string) error
	
	// ListByUser retrieves posts by a specific user with pagination
	ListByUser(userID string, page, limit int) ([]*Post, int, error)
	
	// List retrieves a list of posts with pagination
	List(page, limit int) ([]*PostWithUser, int, error)
}
