package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// MockPostRepository is a mock implementation of domain.PostRepository
type MockPostRepository struct {
	posts map[string]*domain.Post
	// Track method calls for verification
	getByIDCalled      bool
	createCalled       bool
	updateCalled       bool
	deleteCalled       bool
	listByUserCalled   bool
	listCalled         bool
	countByUserCalled  bool
	countCalled        bool
}

// NewMockPostRepository creates a new mock post repository
func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts: make(map[string]*domain.Post),
	}
}

// GetByID retrieves a post by ID
func (m *MockPostRepository) GetByID(id string) (*domain.Post, error) {
	m.getByIDCalled = true
	post, ok := m.posts[id]
	if !ok {
		return nil, domain.ErrPostNotFound
	}
	return post, nil
}

// Create creates a new post
func (m *MockPostRepository) Create(post *domain.Post) error {
	m.createCalled = true
	if _, ok := m.posts[post.ID]; ok {
		return errors.New("post already exists")
	}
	m.posts[post.ID] = post
	return nil
}

// Update updates an existing post
func (m *MockPostRepository) Update(post *domain.Post) error {
	m.updateCalled = true
	if _, ok := m.posts[post.ID]; !ok {
		return domain.ErrPostNotFound
	}
	m.posts[post.ID] = post
	return nil
}

// Delete deletes a post
func (m *MockPostRepository) Delete(id string) error {
	m.deleteCalled = true
	if _, ok := m.posts[id]; !ok {
		return domain.ErrPostNotFound
	}
	delete(m.posts, id)
	return nil
}

// ListByUser retrieves posts by a specific user with pagination
func (m *MockPostRepository) ListByUser(userID string, offset, limit int) ([]*domain.Post, error) {
	m.listByUserCalled = true
	posts := make([]*domain.Post, 0)
	for _, post := range m.posts {
		if post.UserID == userID {
			posts = append(posts, post)
		}
	}
	
	// Apply pagination
	if offset >= len(posts) {
		return []*domain.Post{}, nil
	}
	
	end := offset + limit
	if end > len(posts) {
		end = len(posts)
	}
	
	return posts[offset:end], nil
}

// List retrieves a list of posts with pagination
func (m *MockPostRepository) List(offset, limit int) ([]*domain.PostWithUser, error) {
	m.listCalled = true
	posts := make([]*domain.PostWithUser, 0)
	for _, post := range m.posts {
		posts = append(posts, &domain.PostWithUser{
			Post:     *post,
			Username: "user_" + post.UserID, // Mock username
		})
	}
	
	// Apply pagination
	if offset >= len(posts) {
		return []*domain.PostWithUser{}, nil
	}
	
	end := offset + limit
	if end > len(posts) {
		end = len(posts)
	}
	
	return posts[offset:end], nil
}

// CountByUser returns the total number of posts by a specific user
func (m *MockPostRepository) CountByUser(userID string) (int, error) {
	m.countByUserCalled = true
	count := 0
	for _, post := range m.posts {
		if post.UserID == userID {
			count++
		}
	}
	return count, nil
}

// Count returns the total number of posts
func (m *MockPostRepository) Count() (int, error) {
	m.countCalled = true
	return len(m.posts), nil
}

// TestNewPostService tests the NewPostService function
func TestNewPostService(t *testing.T) {
	// Setup
	postRepo := NewMockPostRepository()
	userRepo := NewMockUserRepository()
	
	// Test
	service := NewPostService(postRepo, userRepo)
	
	// Assert
	if service == nil {
		t.Fatal("NewPostService returned nil")
	}
	if service.postRepo != postRepo {
		t.Errorf("postRepo = %v, want %v", service.postRepo, postRepo)
	}
	if service.userRepo != userRepo {
		t.Errorf("userRepo = %v, want %v", service.userRepo, userRepo)
	}
}

// TestPostGetByID tests the GetByID method for posts
func TestPostGetByID(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		id          string
		setupRepos  func(*MockPostRepository, *MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name: "valid post ID",
			id:   "post_123",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123",
					Content: "Test post",
				}
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: false,
		},
		{
			name:       "empty post ID",
			id:         "",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidPostID,
		},
		{
			name: "non-existent post ID",
			id:   "post_456",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123",
					Content: "Test post",
				}
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: true,
			errorType:   domain.ErrPostNotFound,
		},
		{
			name: "post with non-existent user",
			id:   "post_123",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_456", // Non-existent user
					Content: "Test post",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			postRepo := NewMockPostRepository()
			userRepo := NewMockUserRepository()
			tc.setupRepos(postRepo, userRepo)
			service := NewPostService(postRepo, userRepo)
			
			// Test
			post, err := service.GetByID(tc.id)
			
			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tc.errorType != nil && !errors.Is(err, tc.errorType) {
					t.Errorf("Expected error type %v, got %v", tc.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if post == nil {
					t.Errorf("Expected post, got nil")
				}
				if post.ID != tc.id {
					t.Errorf("post.ID = %q, want %q", post.ID, tc.id)
				}
				if post.Username != "testuser" {
					t.Errorf("post.Username = %q, want %q", post.Username, "testuser")
				}
			}
			
			// Verify repository method was called if we expect it to be
			if !tc.expectError || tc.errorType != domain.ErrInvalidPostID {
				if !postRepo.getByIDCalled {
					t.Errorf("Expected GetByID to be called")
				}
			}
		})
	}
}

// TestCreate tests the Create method
func TestCreate(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		userID      string
		content     string
		setupRepos  func(*MockPostRepository, *MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name:    "valid post creation",
			userID:  "user_123",
			content: "Test post content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: false,
		},
		{
			name:       "empty user ID",
			userID:     "",
			content:    "Test post content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUserID,
		},
		{
			name:       "empty content",
			userID:     "user_123",
			content:    "",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidPostContent,
		},
		{
			name:    "non-existent user",
			userID:  "user_456",
			content: "Test post content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			postRepo := NewMockPostRepository()
			userRepo := NewMockUserRepository()
			tc.setupRepos(postRepo, userRepo)
			service := NewPostService(postRepo, userRepo)
			
			// Test
			post, err := service.Create(tc.userID, tc.content)
			
			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tc.errorType != nil && !errors.Is(err, tc.errorType) {
					t.Errorf("Expected error type %v, got %v", tc.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if post == nil {
					t.Errorf("Expected post, got nil")
				}
				if post.UserID != tc.userID {
					t.Errorf("post.UserID = %q, want %q", post.UserID, tc.userID)
				}
				if post.Content != tc.content {
					t.Errorf("post.Content = %q, want %q", post.Content, tc.content)
				}
				if post.ID == "" {
					t.Errorf("post.ID is empty")
				}
				if post.CreatedAt.IsZero() {
					t.Errorf("post.CreatedAt is zero")
				}
				if post.UpdatedAt.IsZero() {
					t.Errorf("post.UpdatedAt is zero")
				}
			}
			
			// Verify repository methods were called
			if !userRepo.getByIDCalled && !tc.expectError {
				t.Errorf("Expected GetByID to be called on user repository")
			}
			
			if !tc.expectError {
				if !postRepo.createCalled {
					t.Errorf("Expected Create to be called on post repository")
				}
			}
		})
	}
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		id          string
		userID      string
		content     string
		setupRepos  func(*MockPostRepository, *MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name:    "valid post update",
			id:      "post_123",
			userID:  "user_123",
			content: "Updated content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123",
					Content: "Original content",
				}
			},
			expectError: false,
		},
		{
			name:       "empty post ID",
			id:         "",
			userID:     "user_123",
			content:    "Updated content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidPostID,
		},
		{
			name:       "empty user ID",
			id:         "post_123",
			userID:     "",
			content:    "Updated content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUserID,
		},
		{
			name:       "empty content",
			id:         "post_123",
			userID:     "user_123",
			content:    "",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidPostContent,
		},
		{
			name:    "non-existent post",
			id:      "post_456",
			userID:  "user_123",
			content: "Updated content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123",
					Content: "Original content",
				}
			},
			expectError: true,
			errorType:   domain.ErrPostNotFound,
		},
		{
			name:    "post owned by different user",
			id:      "post_123",
			userID:  "user_456",
			content: "Updated content",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123", // Different user
					Content: "Original content",
				}
			},
			expectError: true,
			errorType:   domain.ErrPostNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			postRepo := NewMockPostRepository()
			userRepo := NewMockUserRepository()
			tc.setupRepos(postRepo, userRepo)
			service := NewPostService(postRepo, userRepo)
			
			// Record the time before the update
			beforeUpdate := time.Now()
			
			// Test
			post, err := service.Update(tc.id, tc.userID, tc.content)
			
			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tc.errorType != nil && !errors.Is(err, tc.errorType) {
					t.Errorf("Expected error type %v, got %v", tc.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if post == nil {
					t.Errorf("Expected post, got nil")
				}
				if post.ID != tc.id {
					t.Errorf("post.ID = %q, want %q", post.ID, tc.id)
				}
				if post.UserID != tc.userID {
					t.Errorf("post.UserID = %q, want %q", post.UserID, tc.userID)
				}
				if post.Content != tc.content {
					t.Errorf("post.Content = %q, want %q", post.Content, tc.content)
				}
				if post.UpdatedAt.Before(beforeUpdate) {
					t.Errorf("post.UpdatedAt was not updated")
				}
			}
			
			// Verify repository methods were called
			if !tc.expectError || tc.errorType != domain.ErrInvalidPostID && tc.errorType != domain.ErrInvalidUserID && tc.errorType != domain.ErrInvalidPostContent {
				if !postRepo.getByIDCalled {
					t.Errorf("Expected GetByID to be called on post repository")
				}
			}
			
			if !tc.expectError {
				if !postRepo.updateCalled {
					t.Errorf("Expected Update to be called on post repository")
				}
			}
		})
	}
}

// TestPostDelete tests the Delete method for posts
func TestPostDelete(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		id          string
		userID      string
		setupRepos  func(*MockPostRepository, *MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name:   "valid post deletion",
			id:     "post_123",
			userID: "user_123",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123",
					Content: "Test post",
				}
			},
			expectError: false,
		},
		{
			name:       "empty post ID",
			id:         "",
			userID:     "user_123",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidPostID,
		},
		{
			name:       "empty user ID",
			id:         "post_123",
			userID:     "",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUserID,
		},
		{
			name:   "non-existent post",
			id:     "post_456",
			userID: "user_123",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123",
					Content: "Test post",
				}
			},
			expectError: true,
			errorType:   domain.ErrPostNotFound,
		},
		{
			name:   "post owned by different user",
			id:     "post_123",
			userID: "user_456",
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				postRepo.posts["post_123"] = &domain.Post{
					ID:      "post_123",
					UserID:  "user_123", // Different user
					Content: "Test post",
				}
			},
			expectError: true,
			errorType:   domain.ErrPostNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			postRepo := NewMockPostRepository()
			userRepo := NewMockUserRepository()
			tc.setupRepos(postRepo, userRepo)
			service := NewPostService(postRepo, userRepo)
			
			// Test
			err := service.Delete(tc.id, tc.userID)
			
			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tc.errorType != nil && !errors.Is(err, tc.errorType) {
					t.Errorf("Expected error type %v, got %v", tc.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				
				// Verify post was deleted
				_, err := postRepo.GetByID(tc.id)
				if err == nil {
					t.Errorf("Expected post to be deleted")
				}
			}
			
			// Verify repository methods were called
			if !tc.expectError || tc.errorType != domain.ErrInvalidPostID && tc.errorType != domain.ErrInvalidUserID {
				if !postRepo.getByIDCalled {
					t.Errorf("Expected GetByID to be called on post repository")
				}
			}
			
			if !tc.expectError {
				if !postRepo.deleteCalled {
					t.Errorf("Expected Delete to be called on post repository")
				}
			}
		})
	}
}

// TestListByUser tests the ListByUser method
func TestListByUser(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		userID      string
		page        int
		limit       int
		setupRepos  func(*MockPostRepository, *MockUserRepository)
		expectError bool
		errorType   error
		expectedLen int
	}{
		{
			name:   "valid list with default pagination",
			userID: "user_123",
			page:   0,
			limit:  0,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
				for i := 1; i <= 15; i++ {
					postRepo.posts[fmt.Sprintf("post_%d", i)] = &domain.Post{
						ID:      fmt.Sprintf("post_%d", i),
						UserID:  "user_123",
						Content: fmt.Sprintf("Post %d", i),
					}
				}
			},
			expectError: false,
			expectedLen: 10, // Default limit is 10
		},
		{
			name:   "valid list with custom pagination",
			userID: "user_123",
			page:   2,
			limit:  5,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
				for i := 1; i <= 15; i++ {
					postRepo.posts[fmt.Sprintf("post_%d", i)] = &domain.Post{
						ID:      fmt.Sprintf("post_%d", i),
						UserID:  "user_123",
						Content: fmt.Sprintf("Post %d", i),
					}
				}
			},
			expectError: false,
			expectedLen: 5, // Page 2 with limit 5 should return 5 posts
		},
		{
			name:       "empty user ID",
			userID:     "",
			page:       1,
			limit:      10,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUserID,
		},
		{
			name:   "non-existent user",
			userID: "user_456",
			page:   1,
			limit:  10,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserNotFound,
		},
		{
			name:   "user with no posts",
			userID: "user_123",
			page:   1,
			limit:  10,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				userRepo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
				// No posts for this user
			},
			expectError: false,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			postRepo := NewMockPostRepository()
			userRepo := NewMockUserRepository()
			tc.setupRepos(postRepo, userRepo)
			service := NewPostService(postRepo, userRepo)
			
			// Test
			posts, _, err := service.ListByUser(tc.userID, tc.page, tc.limit)
			
			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if posts == nil {
					t.Errorf("Expected posts, got nil")
				}
				if len(posts) != tc.expectedLen {
					t.Errorf("len(posts) = %d, want %d", len(posts), tc.expectedLen)
				}
				
				// Verify all posts belong to the user
				for _, post := range posts {
					if post.UserID != tc.userID {
						t.Errorf("post.UserID = %q, want %q", post.UserID, tc.userID)
					}
				}
			}
			
			// Verify repository methods were called
			if !tc.expectError || tc.errorType != domain.ErrInvalidUserID {
				if !userRepo.getByIDCalled {
					t.Errorf("Expected GetByID to be called on user repository")
				}
			}
			
			if !tc.expectError {
				if !postRepo.listByUserCalled {
					t.Errorf("Expected ListByUser to be called on post repository")
				}
				if !postRepo.countByUserCalled {
					t.Errorf("Expected CountByUser to be called on post repository")
				}
			}
		})
	}
}

// TestPostList tests the List method for posts
func TestPostList(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		page        int
		limit       int
		setupRepos  func(*MockPostRepository, *MockUserRepository)
		expectError bool
		expectedLen int
	}{
		{
			name:  "valid list with default pagination",
			page:  0,
			limit: 0,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				for i := 1; i <= 15; i++ {
					postRepo.posts[fmt.Sprintf("post_%d", i)] = &domain.Post{
						ID:      fmt.Sprintf("post_%d", i),
						UserID:  fmt.Sprintf("user_%d", i%3+1), // Distribute posts among 3 users
						Content: fmt.Sprintf("Post %d", i),
					}
				}
			},
			expectError: false,
			expectedLen: 10, // Default limit is 10
		},
		{
			name:  "valid list with custom pagination",
			page:  2,
			limit: 5,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				for i := 1; i <= 15; i++ {
					postRepo.posts[fmt.Sprintf("post_%d", i)] = &domain.Post{
						ID:      fmt.Sprintf("post_%d", i),
						UserID:  fmt.Sprintf("user_%d", i%3+1), // Distribute posts among 3 users
						Content: fmt.Sprintf("Post %d", i),
					}
				}
			},
			expectError: false,
			expectedLen: 5, // Page 2 with limit 5 should return 5 posts
		},
		{
			name:  "empty repository",
			page:  1,
			limit: 10,
			setupRepos: func(postRepo *MockPostRepository, userRepo *MockUserRepository) {
				// No posts
			},
			expectError: false,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			postRepo := NewMockPostRepository()
			userRepo := NewMockUserRepository()
			tc.setupRepos(postRepo, userRepo)
			service := NewPostService(postRepo, userRepo)
			
			// Test
			posts, count, err := service.List(tc.page, tc.limit)
			
			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if posts == nil {
					t.Errorf("Expected posts, got nil")
				}
				if len(posts) != tc.expectedLen {
					t.Errorf("len(posts) = %d, want %d", len(posts), tc.expectedLen)
				}
				if count != len(postRepo.posts) {
					t.Errorf("count = %d, want %d", count, len(postRepo.posts))
				}
				
				// Verify all posts have a username
				for _, post := range posts {
					if post.Username == "" {
						t.Errorf("post.Username is empty")
					}
				}
			}
			
			// Verify repository methods were called
			if !postRepo.listCalled {
				t.Errorf("Expected List to be called on post repository")
			}
			if !postRepo.countCalled {
				t.Errorf("Expected Count to be called on post repository")
			}
		})
	}
}
