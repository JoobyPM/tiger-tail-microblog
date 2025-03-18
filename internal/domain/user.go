package domain

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	// GetByID retrieves a user by ID
	GetByID(id string) (*User, error)
	
	// GetByUsername retrieves a user by username
	GetByUsername(username string) (*User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(email string) (*User, error)
	
	// Create creates a new user
	Create(user *User) error
	
	// Update updates an existing user
	Update(user *User) error
	
	// Delete deletes a user
	Delete(id string) error
	
	// List retrieves a list of users with pagination
	List(offset, limit int) ([]*User, error)
	
	// Count returns the total number of users
	Count() (int, error)
}

// UserService defines the interface for user business logic
type UserService interface {
	// GetByID retrieves a user by ID
	GetByID(id string) (*User, error)
	
	// GetByUsername retrieves a user by username
	GetByUsername(username string) (*User, error)
	
	// Register registers a new user
	Register(username, email, password string) (*User, error)
	
	// Authenticate authenticates a user
	Authenticate(usernameOrEmail, password string) (*User, error)
	
	// UpdateProfile updates a user's profile
	UpdateProfile(id, bio string) (*User, error)
	
	// ChangePassword changes a user's password
	ChangePassword(id, currentPassword, newPassword string) error
	
	// Delete deletes a user
	Delete(id string) error
	
	// List retrieves a list of users with pagination
	List(page, limit int) ([]*User, int, error)
}
