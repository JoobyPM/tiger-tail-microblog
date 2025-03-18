package service

import (
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// PostService implements the domain.PostService interface
type PostService struct {
	postRepo domain.PostRepository
	userRepo domain.UserRepository
}

// NewPostService creates a new post service
func NewPostService(postRepo domain.PostRepository, userRepo domain.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// GetByID retrieves a post by ID
func (s *PostService) GetByID(id string) (*domain.PostWithUser, error) {
	if id == "" {
		return nil, domain.ErrInvalidPostID
	}

	// Get post
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(post.UserID)
	if err != nil {
		return nil, err
	}

	// Create post with user
	postWithUser := &domain.PostWithUser{
		Post:     *post,
		Username: user.Username,
	}

	return postWithUser, nil
}

// Create creates a new post
func (s *PostService) Create(userID, content string) (*domain.Post, error) {
	// Validate input
	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}
	if content == "" {
		return nil, domain.ErrInvalidPostContent
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Create post
	now := time.Now()
	post := &domain.Post{
		ID:        generatePostID(),
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save post
	err = s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Update updates an existing post
func (s *PostService) Update(id, userID, content string) (*domain.Post, error) {
	// Validate input
	if id == "" {
		return nil, domain.ErrInvalidPostID
	}
	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}
	if content == "" {
		return nil, domain.ErrInvalidPostContent
	}

	// Get post
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the post
	if post.UserID != userID {
		return nil, domain.ErrPostNotFound
	}

	// Update post
	post.Content = content
	post.UpdatedAt = time.Now()

	// Save post
	err = s.postRepo.Update(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Delete deletes a post
func (s *PostService) Delete(id, userID string) error {
	// Validate input
	if id == "" {
		return domain.ErrInvalidPostID
	}
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	// Get post
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user owns the post
	if post.UserID != userID {
		return domain.ErrPostNotFound
	}

	// Delete post
	return s.postRepo.Delete(id)
}

// ListByUser retrieves posts by a specific user with pagination
func (s *PostService) ListByUser(userID string, page, limit int) ([]*domain.Post, int, error) {
	// Validate input
	if userID == "" {
		return nil, 0, domain.ErrInvalidUserID
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	// Get posts
	posts, err := s.postRepo.ListByUser(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.postRepo.CountByUser(userID)
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// List retrieves a list of posts with pagination
func (s *PostService) List(page, limit int) ([]*domain.PostWithUser, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get posts
	posts, err := s.postRepo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.postRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// generatePostID generates a unique post ID
// In a real application, this would use a proper ID generation method
func generatePostID() string {
	return "post_" + time.Now().Format("20060102150405")
}
