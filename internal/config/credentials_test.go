package config

import (
	"testing"
)

func TestSensitiveString(t *testing.T) {
	testCases := []struct {
		name     string
		input    SensitiveString
		expected string
	}{
		{
			name:     "empty string",
			input:    SensitiveString(""),
			expected: "",
		},
		{
			name:     "short string",
			input:    SensitiveString("abc"),
			expected: "***",
		},
		{
			name:     "longer string",
			input:    SensitiveString("password123"),
			expected: "p***3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.String()
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestSensitiveStringValue(t *testing.T) {
	original := "secret-password"
	sensitive := SensitiveString(original)
	
	if sensitive.Value() != original {
		t.Errorf("Expected Value() to return the original string %q, got %q", original, sensitive.Value())
	}
}

func TestSanitizeConnectionString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "postgres connection string",
			input:    "postgres://user:password@localhost:5432/dbname?sslmode=disable",
			expected: "postgres://[redacted]@localhost:5432/dbname?sslmode=disable",
		},
		{
			name:     "redis connection string",
			input:    "localhost:6379",
			expected: "localhost:6379", // Redis address doesn't contain credentials
		},
		{
			name:     "other connection string",
			input:    "mongodb://user:pass@localhost:27017/dbname",
			expected: "[redacted-connection-string]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SanitizeConnectionString(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestDatabaseCredentials(t *testing.T) {
	creds := DatabaseCredentials{
		Host:     "localhost",
		Port:     "5432",
		User:     SensitiveString("user"),
		Password: SensitiveString("password"),
		Name:     "dbname",
		SSLMode:  "disable",
	}

	// Test GetDSN
	expectedDSN := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	if dsn := creds.GetDSN(); dsn != expectedDSN {
		t.Errorf("Expected DSN %q, got %q", expectedDSN, dsn)
	}

	// Test GetSanitizedDSN
	expectedSanitizedDSN := "postgres://[redacted]@localhost:5432/dbname?sslmode=disable"
	if sanitizedDSN := creds.GetSanitizedDSN(); sanitizedDSN != expectedSanitizedDSN {
		t.Errorf("Expected sanitized DSN %q, got %q", expectedSanitizedDSN, sanitizedDSN)
	}
}

func TestRedisCredentials(t *testing.T) {
	creds := RedisCredentials{
		Host:     "localhost",
		Port:     "6379",
		Password: SensitiveString("password"),
		DB:       0,
	}

	// Test GetAddr
	expectedAddr := "localhost:6379"
	if addr := creds.GetAddr(); addr != expectedAddr {
		t.Errorf("Expected address %q, got %q", expectedAddr, addr)
	}
}
