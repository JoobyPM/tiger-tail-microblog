package db

import (
	"database/sql"
	"testing"
)

func TestNewPostgresConnection(t *testing.T) {
	// Test cases
	testCases := []struct {
		name string
		dsn  string
	}{
		{
			name: "default connection string",
			dsn:  "postgres://postgres:postgres@localhost:5432/tigertail?sslmode=disable",
		},
		{
			name: "custom connection string",
			dsn:  "postgres://user:password@db.example.com:5432/mydb?sslmode=require",
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
			// Since this is a stub, db.db should be nil
			if db.db != nil {
				t.Errorf("Expected db.db to be nil, got: %v", db.db)
			}
		})
	}
}

func TestPostgresDBMethods(t *testing.T) {
	// Setup
	db, err := NewPostgresConnection("postgres://postgres:postgres@localhost:5432/tigertail?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to create PostgreSQL connection: %v", err)
	}

	// Test Close
	t.Run("Close", func(t *testing.T) {
		err := db.Close()
		// Since this is a stub and db.db is nil, we expect no error
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	// Test Ping
	t.Run("Ping", func(t *testing.T) {
		err := db.Ping()
		// Since this is a stub and db.db is nil, we expect an error
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		expectedErr := "database connection not initialized"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message %q, got %q", expectedErr, err.Error())
		}
	})

	// Test Exec
	t.Run("Exec", func(t *testing.T) {
		_, err := db.Exec("SELECT 1")
		// Since this is a stub and db.db is nil, we expect an error
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		expectedErr := "database connection not initialized"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message %q, got %q", expectedErr, err.Error())
		}
	})

	// Test Query
	t.Run("Query", func(t *testing.T) {
		_, err := db.Query("SELECT 1")
		// Since this is a stub and db.db is nil, we expect an error
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		expectedErr := "database connection not initialized"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message %q, got %q", expectedErr, err.Error())
		}
	})

	// Test QueryRow
	t.Run("QueryRow", func(t *testing.T) {
		row := db.QueryRow("SELECT 1")
		// Since this is a stub and db.db is nil, we expect nil
		if row != nil {
			t.Errorf("Expected nil, got: %v", row)
		}
	})
}

// TestWithInitializedDB tests the behavior when db.db is not nil
func TestWithInitializedDB(t *testing.T) {
	// This test is more for demonstration purposes
	// In a real test, we would mock the sql.DB interface
	
	// Create a PostgresDB with a non-nil db field
	db := &PostgresDB{
		db: &sql.DB{}, // This will panic if any methods are called on it
	}

	// We're just testing that the methods check for nil before using db.db
	// So we don't actually call any methods on the sql.DB
	
	// Test Close
	t.Run("Close with initialized db", func(t *testing.T) {
		// We expect this not to panic
		_ = db.Close()
	})
}
