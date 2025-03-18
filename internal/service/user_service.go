package service

import (
	"errors"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// UserService implements the domain.UserService interface
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidUserID
	}

	return s.userRepo.GetByID(id)
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(username string) (*domain.User, error) {
	if username == "" {
		return nil, domain.ErrInvalidUsername
	}

	return s.userRepo.GetByUsername(username)
}

// Register registers a new user
func (s *UserService) Register(username, email, password string) (*domain.User, error) {
	// Validate input
	if username == "" {
		return nil, domain.ErrInvalidUsername
	}
	if email == "" {
		return nil, domain.ErrInvalidEmail
	}
	if password == "" {
		return nil, domain.ErrInvalidPassword
	}

	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(username)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Create user
	// In a real application, we would hash the password here
	now := time.Now()
	user := &domain.User{
		ID:        generateID(), // This would be a real ID generation function
		Username:  username,
		Email:     email,
		Password:  password, // This would be hashed in a real application
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save user
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Authenticate authenticates a user
func (s *UserService) Authenticate(usernameOrEmail, password string) (*domain.User, error) {
	if usernameOrEmail == "" {
		return nil, errors.New("username or email is required")
	}
	if password == "" {
		return nil, domain.ErrInvalidPassword
	}

	// Try to find user by username
	user, err := s.userRepo.GetByUsername(usernameOrEmail)
	if err != nil {
		// If not found by username, try by email
		user, err = s.userRepo.GetByEmail(usernameOrEmail)
		if err != nil {
			return nil, domain.ErrUserNotFound
		}
	}

	// Check password
	// In a real application, we would compare hashed passwords
	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// UpdateProfile updates a user's profile
func (s *UserService) UpdateProfile(id, bio string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidUserID
	}

	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update user
	user.Bio = bio
	user.UpdatedAt = time.Now()

	// Save user
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(id, currentPassword, newPassword string) error {
	if id == "" {
		return domain.ErrInvalidUserID
	}
	if currentPassword == "" || newPassword == "" {
		return domain.ErrInvalidPassword
	}

	// Get user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check current password
	// In a real application, we would compare hashed passwords
	if user.Password != currentPassword {
		return errors.New("invalid current password")
	}

	// Update password
	// In a real application, we would hash the new password
	user.Password = newPassword
	user.UpdatedAt = time.Now()

	// Save user
	return s.userRepo.Update(user)
}

// Delete deletes a user
func (s *UserService) Delete(id string) error {
	if id == "" {
		return domain.ErrInvalidUserID
	}

	return s.userRepo.Delete(id)
}

// List retrieves a list of users with pagination
func (s *UserService) List(page, limit int) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get users
	users, err := s.userRepo.List(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// generateID generates a unique ID
// In a real application, this would use a proper ID generation method
func generateID() string {
	return "user_" + time.Now().Format("20060102150405")
}
