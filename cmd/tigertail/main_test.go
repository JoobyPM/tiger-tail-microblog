package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test cases
	testCases := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable is set",
			key:          "TEST_ENV_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "environment variable is not set",
			key:          "TEST_ENV_VAR_NONEXISTENT",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "environment variable is empty",
			key:          "TEST_ENV_VAR_EMPTY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			if tc.envValue != "" {
				os.Setenv(tc.key, tc.envValue)
				defer os.Unsetenv(tc.key)
			}

			// Test
			result := getEnv(tc.key, tc.defaultValue)

			// Assert
			if result != tc.expected {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tc.key, tc.defaultValue, result, tc.expected)
			}
		})
	}
}
