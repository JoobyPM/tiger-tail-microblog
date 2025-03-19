package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/JoobyPM/tiger-tail-microblog/internal/config"
	"github.com/JoobyPM/tiger-tail-microblog/internal/domain"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresStub creates a new stub PostgreSQL connection for testing
func NewPostgresStub() *PostgresDB {
	log.Println("Creating PostgreSQL stub")
	return &PostgresDB{
		db: nil,
	}
}

// NewPostgresConnection creates a new PostgreSQL connection
func NewPostgresConnection(dsn string) (*PostgresDB, error) {
	log.Printf("Connecting to PostgreSQL with DSN: %s", config.SanitizeConnectionString(dsn))
	
	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	// Initialize database
	postgres := &PostgresDB{
		db: db,
	}
	
	if err := postgres.initializeDatabase(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
	}
	
	return postgres, nil
}

// initializeDatabase creates the necessary tables if they don't exist
func (p *PostgresDB) initializeDatabase() error {
	if p.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	
	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)
	`
	
	_, err := p.db.Exec(usersTable)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	
	// Create posts table
	postsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)
	`
	
	_, err = p.db.Exec(postsTable)
	if err != nil {
		return fmt.Errorf("error creating posts table: %w", err)
	}
	
	// Check if default user exists
	var count int
	err = p.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for default user: %w", err)
	}
	
	// Create default user if it doesn't exist
	if count == 0 {
		_, err = p.db.Exec(
			"INSERT INTO users (id, username, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
			"user_1",
			"admin",
			"password", // In a real app, this would be hashed
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("error creating default user: %w", err)
		}
		log.Println("Created default user: admin")
	}
	
	log.Println("Database initialized successfully")
	return nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	if p.db == nil {
		return nil
	}
	
	log.Println("Closing PostgreSQL connection")
	return p.db.Close()
}

// Ping checks if the database connection is alive
func (p *PostgresDB) Ping() error {
	if p.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	
	return p.db.Ping()
}

// Exec executes a query without returning any rows
func (p *PostgresDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	
	return p.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (p *PostgresDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	
	return p.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (p *PostgresDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if p.db == nil {
		log.Printf("Error: database connection not initialized")
		return nil
	}
	
	return p.db.QueryRow(query, args...)
}

// PostRepository implements the domain.PostRepository interface
type PostRepository struct {
	db *PostgresDB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *PostgresDB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// GetByID retrieves a post by ID
func (r *PostRepository) GetByID(id string) (*domain.Post, error) {
	if r.db.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	
	query := "SELECT id, user_id, content, created_at, updated_at FROM posts WHERE id = $1"
	row := r.db.QueryRow(query, id)
	
	var post domain.Post
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPostNotFound
		}
		return nil, fmt.Errorf("error scanning post row: %w", err)
	}
	
	return &post, nil
}

// Create creates a new post
func (r *PostRepository) Create(post *domain.Post) error {
	if r.db.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	
	query := "INSERT INTO posts (id, user_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, post.ID, post.UserID, post.Content, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}
	
	return nil
}

// Update updates an existing post
func (r *PostRepository) Update(post *domain.Post) error {
	if r.db.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	
	query := "UPDATE posts SET content = $1, updated_at = $2 WHERE id = $3"
	result, err := r.db.Exec(query, post.Content, post.UpdatedAt, post.ID)
	if err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrPostNotFound
	}
	
	return nil
}

// Delete deletes a post
func (r *PostRepository) Delete(id string) error {
	if r.db.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	
	query := "DELETE FROM posts WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return domain.ErrPostNotFound
	}
	
	return nil
}

// ListByUser retrieves posts by a specific user with pagination
func (r *PostRepository) ListByUser(userID string, offset, limit int) ([]*domain.Post, error) {
	if r.db.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	
	query := "SELECT id, user_id, content, created_at, updated_at FROM posts WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying posts by user: %w", err)
	}
	defer rows.Close()
	
	posts := make([]*domain.Post, 0)
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		posts = append(posts, &post)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}
	
	return posts, nil
}

// List retrieves a list of posts with pagination
func (r *PostRepository) List(offset, limit int) ([]*domain.PostWithUser, error) {
	if r.db.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	
	// First, try to get posts with user information
	query := `
		SELECT p.id, p.user_id, p.content, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		// If the query fails (e.g., no users table yet), fall back to just getting posts
		log.Printf("Error querying posts with users: %v, falling back to posts-only query", err)
		return r.listPostsOnly(offset, limit)
	}
	defer rows.Close()
	
	posts := make([]*domain.PostWithUser, 0)
	for rows.Next() {
		var post domain.PostWithUser
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Username)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		posts = append(posts, &post)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}
	
	// If no posts were found with the JOIN, try the fallback method
	if len(posts) == 0 {
		return r.listPostsOnly(offset, limit)
	}
	
	return posts, nil
}

// listPostsOnly retrieves posts without user information
func (r *PostRepository) listPostsOnly(offset, limit int) ([]*domain.PostWithUser, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		// If this also fails, there might be no posts table yet
		if err.Error() == `pq: relation "posts" does not exist` {
			// Return empty result instead of error for better UX
			return []*domain.PostWithUser{}, nil
		}
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()
	
	posts := make([]*domain.PostWithUser, 0)
	for rows.Next() {
		var post domain.PostWithUser
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		// Set a default username since we don't have user information
		post.Username = "unknown"
		posts = append(posts, &post)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}
	
	return posts, nil
}

// CountByUser returns the total number of posts by a specific user
func (r *PostRepository) CountByUser(userID string) (int, error) {
	if r.db.db == nil {
		return 0, fmt.Errorf("database connection not initialized")
	}
	
	query := "SELECT COUNT(*) FROM posts WHERE user_id = $1"
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting posts by user: %w", err)
	}
	
	return count, nil
}

// Count returns the total number of posts
func (r *PostRepository) Count() (int, error) {
	if r.db.db == nil {
		return 0, fmt.Errorf("database connection not initialized")
	}
	
	query := "SELECT COUNT(*) FROM posts"
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		// If the table doesn't exist yet, return 0 instead of an error
		if err.Error() == `pq: relation "posts" does not exist` {
			return 0, nil
		}
		return 0, fmt.Errorf("error counting posts: %w", err)
	}
	
	return count, nil
}

// FetchAllPosts retrieves all posts from the database
func (r *PostRepository) FetchAllPosts() ([]*domain.Post, error) {
	if r.db.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	
	query := "SELECT id, user_id, content, created_at, updated_at FROM posts ORDER BY created_at DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying all posts: %w", err)
	}
	defer rows.Close()
	
	posts := make([]*domain.Post, 0)
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning post row: %w", err)
		}
		posts = append(posts, &post)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}
	
	return posts, nil
}

// CreatePost creates a new post in the database
func (r *PostRepository) CreatePost(post *domain.Post) error {
	return r.Create(post) // Reuse the Create method
}
