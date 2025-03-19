package config

import (
	"fmt"
	"strings"
)

// SensitiveString is a type that wraps a string containing sensitive information
// It prevents accidental logging of the sensitive value by implementing custom String() method
type SensitiveString string

// String returns a redacted version of the sensitive string for logging
func (s SensitiveString) String() string {
	if s == "" {
		return ""
	}
	
	// For very short strings, just return "***"
	if len(s) < 4 {
		return "***"
	}
	
	// For longer strings, show first and last character with stars in between
	return fmt.Sprintf("%c***%c", s[0], s[len(s)-1])
}

// Value returns the actual sensitive string value
// This should only be used when the actual value is needed (e.g., for authentication)
func (s SensitiveString) Value() string {
	return string(s)
}

// SanitizeConnectionString redacts sensitive information from a connection string
func SanitizeConnectionString(dsn string) string {
	// Handle postgres connection strings
	if strings.HasPrefix(dsn, "postgres://") {
		// Extract components
		parts := strings.Split(dsn, "@")
		if len(parts) != 2 {
			return "postgres://[redacted]"
		}
		
		credParts := strings.Split(parts[0], "://")
		if len(credParts) != 2 {
			return "postgres://[redacted]@" + parts[1]
		}
		
		// Redact username and password
		return "postgres://[redacted]@" + parts[1]
	}
	
	// Handle redis connection strings (host:port)
	if strings.Contains(dsn, ":") && !strings.Contains(dsn, "://") {
		return dsn // Redis address doesn't contain credentials
	}
	
	// For other connection strings, do a generic redaction
	return "[redacted-connection-string]"
}

// DatabaseCredentials holds database connection credentials
type DatabaseCredentials struct {
	Host     string
	Port     string
	User     SensitiveString
	Password SensitiveString
	Name     string
	SSLMode  string
}

// GetDSN returns the database connection string
func (c *DatabaseCredentials) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User.Value(), c.Password.Value(), c.Host, c.Port, c.Name, c.SSLMode)
}

// GetSanitizedDSN returns a sanitized version of the DSN for logging
func (c *DatabaseCredentials) GetSanitizedDSN() string {
	return fmt.Sprintf("postgres://[redacted]@%s:%s/%s?sslmode=%s",
		c.Host, c.Port, c.Name, c.SSLMode)
}

// RedisCredentials holds Redis connection credentials
type RedisCredentials struct {
	Host     string
	Port     string
	Password SensitiveString
	DB       int
}

// GetAddr returns the Redis address (host:port)
func (c *RedisCredentials) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// AuthCredentials holds authentication credentials
type AuthCredentials struct {
	Username SensitiveString
	Password SensitiveString
}
