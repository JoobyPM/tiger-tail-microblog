package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// MockUserRepository is a mock implementation of domain.UserRepository
type MockUserRepository struct {
	users map[string]*domain.User
	// Track method calls for verification
	getByIDCalled      bool
	getByUsernameCalled bool
	getByEmailCalled   bool
	createCalled       bool
	updateCalled       bool
	deleteCalled       bool
	listCalled         bool
	countCalled        bool
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

// GetByID retrieves a user by ID
func (m *MockUserRepository) GetByID(id string) (*domain.User, error) {
	m.getByIDCalled = true
	user, ok := m.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

// GetByUsername retrieves a user by username
func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	m.getByUsernameCalled = true
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// GetByEmail retrieves a user by email
func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	m.getByEmailCalled = true
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// Create creates a new user
func (m *MockUserRepository) Create(user *domain.User) error {
	m.createCalled = true
	if _, ok := m.users[user.ID]; ok {
		return domain.ErrUserAlreadyExists
	}
	m.users[user.ID] = user
	return nil
}

// Update updates an existing user
func (m *MockUserRepository) Update(user *domain.User) error {
	m.updateCalled = true
	if _, ok := m.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}

// Delete deletes a user
func (m *MockUserRepository) Delete(id string) error {
	m.deleteCalled = true
	if _, ok := m.users[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

// List retrieves a list of users with pagination
func (m *MockUserRepository) List(offset, limit int) ([]*domain.User, error) {
	m.listCalled = true
	users := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	
	// Apply pagination
	if offset >= len(users) {
		return []*domain.User{}, nil
	}
	
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}
	
	return users[offset:end], nil
}

// Count returns the total number of users
func (m *MockUserRepository) Count() (int, error) {
	m.countCalled = true
	return len(m.users), nil
}

// TestNewUserService tests the NewUserService function
func TestNewUserService(t *testing.T) {
	// Setup
	repo := NewMockUserRepository()
	
	// Test
	service := NewUserService(repo)
	
	// Assert
	if service == nil {
		t.Fatal("NewUserService returned nil")
	}
	if service.userRepo != repo {
		t.Errorf("userRepo = %v, want %v", service.userRepo, repo)
	}
}

// TestUserGetByID tests the GetByID method for users
func TestUserGetByID(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		id          string
		setupRepo   func(*MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name: "valid user ID",
			id:   "user_123",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: false,
		},
		{
			name:        "empty user ID",
			id:          "",
			setupRepo:   func(repo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUserID,
		},
		{
			name: "non-existent user ID",
			id:   "user_456",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
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
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Test
			user, err := service.GetByID(tc.id)
			
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
				if user == nil {
					t.Errorf("Expected user, got nil")
				}
				if user.ID != tc.id {
					t.Errorf("user.ID = %q, want %q", user.ID, tc.id)
				}
			}
			
			// Verify repository method was called if we expect it to be
			if !tc.expectError || tc.errorType != domain.ErrInvalidUserID {
				if !repo.getByIDCalled {
					t.Errorf("Expected GetByID to be called")
				}
			}
		})
	}
}

// TestGetByUsername tests the GetByUsername method
func TestGetByUsername(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		username    string
		setupRepo   func(*MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name:     "valid username",
			username: "testuser",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: false,
		},
		{
			name:        "empty username",
			username:    "",
			setupRepo:   func(repo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUsername,
		},
		{
			name:     "non-existent username",
			username: "nonexistent",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
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
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Test
			user, err := service.GetByUsername(tc.username)
			
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
				if user == nil {
					t.Errorf("Expected user, got nil")
				}
				if user.Username != tc.username {
					t.Errorf("user.Username = %q, want %q", user.Username, tc.username)
				}
			}
			
			// Verify repository method was called if we expect it to be
			if !tc.expectError || tc.errorType != domain.ErrInvalidUsername {
				if !repo.getByUsernameCalled {
					t.Errorf("Expected GetByUsername to be called")
				}
			}
		})
	}
}

// TestRegister tests the Register method
func TestRegister(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		username    string
		email       string
		password    string
		setupRepo   func(*MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name:      "valid registration",
			username:  "newuser",
			email:     "newuser@example.com",
			password:  "password123",
			setupRepo: func(repo *MockUserRepository) {},
			expectError: false,
		},
		{
			name:        "empty username",
			username:    "",
			email:       "newuser@example.com",
			password:    "password123",
			setupRepo:   func(repo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUsername,
		},
		{
			name:        "empty email",
			username:    "newuser",
			email:       "",
			password:    "password123",
			setupRepo:   func(repo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidEmail,
		},
		{
			name:        "empty password",
			username:    "newuser",
			email:       "newuser@example.com",
			password:    "",
			setupRepo:   func(repo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidPassword,
		},
		{
			name:      "username already exists",
			username:  "existinguser",
			email:     "newuser@example.com",
			password:  "password123",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "existinguser",
					Email:    "existing@example.com",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserAlreadyExists,
		},
		{
			name:      "email already exists",
			username:  "newuser",
			email:     "existing@example.com",
			password:  "password123",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "existinguser",
					Email:    "existing@example.com",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Test
			user, err := service.Register(tc.username, tc.email, tc.password)
			
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
				if user == nil {
					t.Errorf("Expected user, got nil")
				}
				if user.Username != tc.username {
					t.Errorf("user.Username = %q, want %q", user.Username, tc.username)
				}
				if user.Email != tc.email {
					t.Errorf("user.Email = %q, want %q", user.Email, tc.email)
				}
				if user.Password != tc.password {
					t.Errorf("user.Password = %q, want %q", user.Password, tc.password)
				}
				if user.ID == "" {
					t.Errorf("user.ID is empty")
				}
				if user.CreatedAt.IsZero() {
					t.Errorf("user.CreatedAt is zero")
				}
				if user.UpdatedAt.IsZero() {
					t.Errorf("user.UpdatedAt is zero")
				}
			}
			
			// Verify repository methods were called
			if !tc.expectError || tc.errorType != domain.ErrInvalidUsername && tc.errorType != domain.ErrInvalidEmail && tc.errorType != domain.ErrInvalidPassword {
				if !repo.getByUsernameCalled {
					t.Errorf("Expected GetByUsername to be called")
				}
			}
			
			if !tc.expectError || tc.errorType != domain.ErrInvalidUsername && tc.errorType != domain.ErrInvalidEmail && tc.errorType != domain.ErrInvalidPassword && tc.errorType != domain.ErrUserAlreadyExists {
				if !repo.getByEmailCalled {
					t.Errorf("Expected GetByEmail to be called")
				}
			}
			
			if !tc.expectError {
				if !repo.createCalled {
					t.Errorf("Expected Create to be called")
				}
			}
		})
	}
}

// TestAuthenticate tests the Authenticate method
func TestAuthenticate(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		usernameOrEmail string
		password       string
		setupRepo      func(*MockUserRepository)
		expectError    bool
		errorType      error
	}{
		{
			name:           "valid authentication by username",
			usernameOrEmail: "testuser",
			password:       "password123",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}
			},
			expectError: false,
		},
		{
			name:           "valid authentication by email",
			usernameOrEmail: "test@example.com",
			password:       "password123",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}
			},
			expectError: false,
		},
		{
			name:           "empty username or email",
			usernameOrEmail: "",
			password:       "password123",
			setupRepo:      func(repo *MockUserRepository) {},
			expectError:    true,
		},
		{
			name:           "empty password",
			usernameOrEmail: "testuser",
			password:       "",
			setupRepo:      func(repo *MockUserRepository) {},
			expectError:    true,
			errorType:      domain.ErrInvalidPassword,
		},
		{
			name:           "non-existent user",
			usernameOrEmail: "nonexistent",
			password:       "password123",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserNotFound,
		},
		{
			name:           "incorrect password",
			usernameOrEmail: "testuser",
			password:       "wrongpassword",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Test
			user, err := service.Authenticate(tc.usernameOrEmail, tc.password)
			
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
				if user == nil {
					t.Errorf("Expected user, got nil")
				}
				if user.Username != "testuser" {
					t.Errorf("user.Username = %q, want %q", user.Username, "testuser")
				}
			}
		})
	}
}

// TestUpdateProfile tests the UpdateProfile method
func TestUpdateProfile(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		id          string
		bio         string
		setupRepo   func(*MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name: "valid update",
			id:   "user_123",
			bio:  "New bio",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Bio:      "Old bio",
				}
			},
			expectError: false,
		},
		{
			name:        "empty user ID",
			id:          "",
			bio:         "New bio",
			setupRepo:   func(repo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUserID,
		},
		{
			name: "non-existent user ID",
			id:   "user_456",
			bio:  "New bio",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Bio:      "Old bio",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Record the time before the update
			beforeUpdate := time.Now()
			
			// Test
			user, err := service.UpdateProfile(tc.id, tc.bio)
			
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
				if user == nil {
					t.Errorf("Expected user, got nil")
				}
				if user.Bio != tc.bio {
					t.Errorf("user.Bio = %q, want %q", user.Bio, tc.bio)
				}
				if user.UpdatedAt.Before(beforeUpdate) {
					t.Errorf("user.UpdatedAt was not updated")
				}
			}
			
			// Verify repository methods were called
			if !tc.expectError || tc.errorType != domain.ErrInvalidUserID {
				if !repo.getByIDCalled {
					t.Errorf("Expected GetByID to be called")
				}
			}
			
			if !tc.expectError {
				if !repo.updateCalled {
					t.Errorf("Expected Update to be called")
				}
			}
		})
	}
}

// TestChangePassword tests the ChangePassword method
func TestChangePassword(t *testing.T) {
	// Test cases
	testCases := []struct {
		name            string
		id              string
		currentPassword string
		newPassword     string
		setupRepo       func(*MockUserRepository)
		expectError     bool
		errorType       error
	}{
		{
			name:            "valid password change",
			id:              "user_123",
			currentPassword: "oldpassword",
			newPassword:     "newpassword",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Password: "oldpassword",
				}
			},
			expectError: false,
		},
		{
			name:            "empty user ID",
			id:              "",
			currentPassword: "oldpassword",
			newPassword:     "newpassword",
			setupRepo:       func(repo *MockUserRepository) {},
			expectError:     true,
			errorType:       domain.ErrInvalidUserID,
		},
		{
			name:            "empty current password",
			id:              "user_123",
			currentPassword: "",
			newPassword:     "newpassword",
			setupRepo:       func(repo *MockUserRepository) {},
			expectError:     true,
			errorType:       domain.ErrInvalidPassword,
		},
		{
			name:            "empty new password",
			id:              "user_123",
			currentPassword: "oldpassword",
			newPassword:     "",
			setupRepo:       func(repo *MockUserRepository) {},
			expectError:     true,
			errorType:       domain.ErrInvalidPassword,
		},
		{
			name:            "non-existent user ID",
			id:              "user_456",
			currentPassword: "oldpassword",
			newPassword:     "newpassword",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Password: "oldpassword",
				}
			},
			expectError: true,
			errorType:   domain.ErrUserNotFound,
		},
		{
			name:            "incorrect current password",
			id:              "user_123",
			currentPassword: "wrongpassword",
			newPassword:     "newpassword",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
					Password: "oldpassword",
				}
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Record the time before the update
			beforeUpdate := time.Now()
			
			// Test
			err := service.ChangePassword(tc.id, tc.currentPassword, tc.newPassword)
			
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
				
				// Verify password was updated
				user, err := repo.GetByID(tc.id)
				if err != nil {
					t.Errorf("Unexpected error getting user: %v", err)
				}
				if user.Password != tc.newPassword {
					t.Errorf("user.Password = %q, want %q", user.Password, tc.newPassword)
				}
				if user.UpdatedAt.Before(beforeUpdate) {
					t.Errorf("user.UpdatedAt was not updated")
				}
			}
			
			// Verify repository methods were called
			if !tc.expectError || tc.errorType != domain.ErrInvalidUserID && tc.errorType != domain.ErrInvalidPassword {
				if !repo.getByIDCalled {
					t.Errorf("Expected GetByID to be called")
				}
			}
			
			if !tc.expectError {
				if !repo.updateCalled {
					t.Errorf("Expected Update to be called")
				}
			}
		})
	}
}

// TestUserDelete tests the Delete method for users
func TestUserDelete(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		id          string
		setupRepo   func(*MockUserRepository)
		expectError bool
		errorType   error
	}{
		{
			name: "valid delete",
			id:   "user_123",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
					ID:       "user_123",
					Username: "testuser",
				}
			},
			expectError: false,
		},
		{
			name:        "empty user ID",
			id:          "",
			setupRepo:   func(repo *MockUserRepository) {},
			expectError: true,
			errorType:   domain.ErrInvalidUserID,
		},
		{
			name: "non-existent user ID",
			id:   "user_456",
			setupRepo: func(repo *MockUserRepository) {
				repo.users["user_123"] = &domain.User{
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
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Test
			err := service.Delete(tc.id)
			
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
				
				// Verify user was deleted
				_, err := repo.GetByID(tc.id)
				if err == nil {
					t.Errorf("Expected user to be deleted")
				}
			}
			
			// Verify repository method was called
			if !tc.expectError || tc.errorType != domain.ErrInvalidUserID {
				if !repo.deleteCalled {
					t.Errorf("Expected Delete to be called")
				}
			}
		})
	}
}

// TestUserList tests the List method for users
func TestUserList(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		page        int
		limit       int
		setupRepo   func(*MockUserRepository)
		expectError bool
		expectedLen int
	}{
		{
			name:  "valid list with default pagination",
			page:  0,
			limit: 0,
			setupRepo: func(repo *MockUserRepository) {
				for i := 1; i <= 15; i++ {
					repo.users[fmt.Sprintf("user_%d", i)] = &domain.User{
						ID:       fmt.Sprintf("user_%d", i),
						Username: fmt.Sprintf("user%d", i),
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
			setupRepo: func(repo *MockUserRepository) {
				for i := 1; i <= 15; i++ {
					repo.users[fmt.Sprintf("user_%d", i)] = &domain.User{
						ID:       fmt.Sprintf("user_%d", i),
						Username: fmt.Sprintf("user%d", i),
					}
				}
			},
			expectError: false,
			expectedLen: 5, // Page 2 with limit 5 should return 5 users
		},
		{
			name:  "empty repository",
			page:  1,
			limit: 10,
			setupRepo: func(repo *MockUserRepository) {
				// No users
			},
			expectError: false,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			repo := NewMockUserRepository()
			tc.setupRepo(repo)
			service := NewUserService(repo)
			
			// Test
			users, count, err := service.List(tc.page, tc.limit)
			
			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if users == nil {
					t.Errorf("Expected users, got nil")
				}
				if len(users) != tc.expectedLen {
					t.Errorf("len(users) = %d, want %d", len(users), tc.expectedLen)
				}
				if count != len(repo.users) {
					t.Errorf("count = %d, want %d", count, len(repo.users))
				}
			}
			
			// Verify repository methods were called
			if !repo.listCalled {
				t.Errorf("Expected List to be called")
			}
			if !repo.countCalled {
				t.Errorf("Expected Count to be called")
			}
		})
	}
}
