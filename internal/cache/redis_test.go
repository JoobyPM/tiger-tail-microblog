package cache

import (
	"testing"
	"time"
)

// Using the MockRedisClient from post_cache_test.go

func TestNewRedisClient(t *testing.T) {
	// Skip this test in CI environments since we don't have a real Redis server
	t.Skip("Skipping test that requires a real Redis server")
	
	// Test cases
	testCases := []struct {
		name     string
		addr     string
		password string
		db       int
	}{
		{
			name:     "default connection",
			addr:     "localhost:6379",
			password: "",
			db:       0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			client, err := NewRedisClient(tc.addr, tc.password, tc.db)

			// Assert
			if err != nil {
				t.Fatalf("NewRedisClient(%q, %q, %d) returned error: %v", tc.addr, tc.password, tc.db, err)
			}
			if client == nil {
				t.Fatalf("NewRedisClient(%q, %q, %d) returned nil client", tc.addr, tc.password, tc.db)
			}
			
			// Clean up
			client.Close()
		})
	}
}

func TestMockRedisClientMethods(t *testing.T) {
	// Setup
	client := NewMockRedisClient()

	// Test Set and Get
	t.Run("Set and Get", func(t *testing.T) {
		// Set a value
		err := client.Set("test_key", []byte("test_value"), time.Minute)
		if err != nil {
			t.Errorf("Set returned error: %v", err)
		}
		
		// Get the value
		val, err := client.Get("test_key")
		if err != nil {
			t.Errorf("Get returned error: %v", err)
		}
		if string(val) != "test_value" {
			t.Errorf("Get returned %q, want %q", string(val), "test_value")
		}
		
		// Get a non-existent key
		_, err = client.Get("non_existent_key")
		if err != ErrCacheMiss {
			t.Errorf("Expected ErrCacheMiss, got: %v", err)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		// Set a value
		err := client.Set("test_key", []byte("test_value"), time.Minute)
		if err != nil {
			t.Errorf("Set returned error: %v", err)
		}
		
		// Delete the key
		err = client.Delete("test_key")
		if err != nil {
			t.Errorf("Delete returned error: %v", err)
		}
		
		// Verify the key is gone
		_, err = client.Get("test_key")
		if err != ErrCacheMiss {
			t.Errorf("Expected ErrCacheMiss, got: %v", err)
		}
	})

	// Test Exists
	t.Run("Exists", func(t *testing.T) {
		// Set a value
		err := client.Set("test_key", []byte("test_value"), time.Minute)
		if err != nil {
			t.Errorf("Set returned error: %v", err)
		}
		
		// Check if the key exists
		exists, err := client.Exists("test_key")
		if err != nil {
			t.Errorf("Exists returned error: %v", err)
		}
		if !exists {
			t.Errorf("Expected exists to be true, got false")
		}
		
		// Check if a non-existent key exists
		exists, err = client.Exists("non_existent_key")
		if err != nil {
			t.Errorf("Exists returned error: %v", err)
		}
		if exists {
			t.Errorf("Expected exists to be false, got true")
		}
	})

	// Test FlushDB
	t.Run("FlushDB", func(t *testing.T) {
		// Set some values
		err := client.Set("test_key1", []byte("test_value1"), time.Minute)
		if err != nil {
			t.Errorf("Set returned error: %v", err)
		}
		err = client.Set("test_key2", []byte("test_value2"), time.Minute)
		if err != nil {
			t.Errorf("Set returned error: %v", err)
		}
		
		// Flush the database
		err = client.FlushDB()
		if err != nil {
			t.Errorf("FlushDB returned error: %v", err)
		}
		
		// Verify the keys are gone
		_, err = client.Get("test_key1")
		if err != ErrCacheMiss {
			t.Errorf("Expected ErrCacheMiss, got: %v", err)
		}
		_, err = client.Get("test_key2")
		if err != ErrCacheMiss {
			t.Errorf("Expected ErrCacheMiss, got: %v", err)
		}
	})
}
