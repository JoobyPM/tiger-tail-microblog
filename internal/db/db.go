package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Variable for sql.Open to allow mocking in tests
var sqlOpen = sql.Open

// Config represents the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Connection represents a database connection
type Connection struct {
	db *sql.DB
}

// New creates a new database connection
func New(config Config) (*Connection, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode,
	)

	// Open connection
	db, err := sqlOpen("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &Connection{db: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}

// DB returns the underlying database connection
func (c *Connection) DB() *sql.DB {
	return c.db
}

// Ping pings the database to verify the connection
func (c *Connection) Ping() error {
	return c.db.Ping()
}

// Begin starts a new transaction
func (c *Connection) Begin() (*sql.Tx, error) {
	return c.db.Begin()
}

// Exec executes a query without returning any rows
func (c *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (c *Connection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (c *Connection) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.db.QueryRow(query, args...)
}
