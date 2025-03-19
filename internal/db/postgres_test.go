package db

import (
	"testing"
)

func TestNewPostgresConnection(t *testing.T) {
	// Skip this test in CI environments since we don't have a real PostgreSQL server
	t.Skip("Skipping test that requires a real PostgreSQL server")
	
	// Test cases
	testCases := []struct {
		name string
		dsn  string
	}{
		{
			name: "default connection string",
			dsn:  "postgres://postgres:postgres@localhost:5432/tigertail?sslmode=disable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			db, err := NewPostgresConnection(tc.dsn)

			// Assert
			if err != nil {
				t.Fatalf("NewPostgresConnection(%q) returned error: %v", tc.dsn, err)
			}
			if db == nil {
				t.Fatalf("NewPostgresConnection(%q) returned nil db", tc.dsn)
			}
			
			// Clean up
			db.Close()
		})
	}
}

func TestPostgresDBMethods(t *testing.T) {
	// Skip this test in CI environments since we don't have a real PostgreSQL server
	t.Skip("Skipping test that requires a real PostgreSQL server")
}

// TestWithEmptyDB tests the behavior with an empty database
func TestWithEmptyDB(t *testing.T) {
	// Create an empty PostgresDB
	db := &PostgresDB{
		db: nil,
	}
	
	// Test methods with nil db
	t.Run("Methods with nil db", func(t *testing.T) {
		// Ping
		err := db.Ping()
		if err == nil || err.Error() != "database connection not initialized" {
			t.Errorf("Expected 'database connection not initialized' error, got: %v", err)
		}
		
		// Exec
		_, err = db.Exec("SELECT 1")
		if err == nil || err.Error() != "database connection not initialized" {
			t.Errorf("Expected 'database connection not initialized' error, got: %v", err)
		}
		
		// Query
		_, err = db.Query("SELECT 1")
		if err == nil || err.Error() != "database connection not initialized" {
			t.Errorf("Expected 'database connection not initialized' error, got: %v", err)
		}
		
		// QueryRow
		row := db.QueryRow("SELECT 1")
		if row != nil {
			t.Errorf("Expected nil, got: %v", row)
		}
		
		// Close
		err = db.Close()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}
