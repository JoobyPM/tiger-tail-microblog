package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Verify default server config
	if config.Server.Port != 8080 {
		t.Errorf("Default server port = %d, want %d", config.Server.Port, 8080)
	}
	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Default server host = %s, want %s", config.Server.Host, "0.0.0.0")
	}
	if config.Server.BaseURL != "http://localhost:8080" {
		t.Errorf("Default server base URL = %s, want %s", config.Server.BaseURL, "http://localhost:8080")
	}

	// Verify default database config
	if config.Database.Host != "localhost" {
		t.Errorf("Default database host = %s, want %s", config.Database.Host, "localhost")
	}
	if config.Database.Port != 5432 {
		t.Errorf("Default database port = %d, want %d", config.Database.Port, 5432)
	}
	if config.Database.User != "postgres" {
		t.Errorf("Default database user = %s, want %s", config.Database.User, "postgres")
	}
	if config.Database.Password != "postgres" {
		t.Errorf("Default database password = %s, want %s", config.Database.Password, "postgres")
	}
	if config.Database.Name != "tigertail" {
		t.Errorf("Default database name = %s, want %s", config.Database.Name, "tigertail")
	}
	if config.Database.SSLMode != "disable" {
		t.Errorf("Default database SSL mode = %s, want %s", config.Database.SSLMode, "disable")
	}

	// Verify default cache config
	if config.Cache.Enabled != false {
		t.Errorf("Default cache enabled = %t, want %t", config.Cache.Enabled, false)
	}
	if config.Cache.Host != "localhost" {
		t.Errorf("Default cache host = %s, want %s", config.Cache.Host, "localhost")
	}
	if config.Cache.Port != 6379 {
		t.Errorf("Default cache port = %d, want %d", config.Cache.Port, 6379)
	}
	if config.Cache.Password != "" {
		t.Errorf("Default cache password = %s, want %s", config.Cache.Password, "")
	}
	if config.Cache.DB != 0 {
		t.Errorf("Default cache DB = %d, want %d", config.Cache.DB, 0)
	}
}

func TestLoadConfig(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		configPath  string
		setupConfig func() string
		expectError bool
		validate    func(*Config, *testing.T)
	}{
		{
			name:       "empty path returns default config",
			configPath: "",
			setupConfig: func() string {
				return ""
			},
			expectError: false,
			validate: func(config *Config, t *testing.T) {
				// Should be default config
				if config.Server.Port != 8080 {
					t.Errorf("Server port = %d, want %d", config.Server.Port, 8080)
				}
			},
		},
		{
			name:       "valid config file",
			configPath: "test_config.json",
			setupConfig: func() string {
				// Create a test config file
				testConfig := Config{
					Server: ServerConfig{
						Port:    9090,
						Host:    "127.0.0.1",
						BaseURL: "http://example.com",
					},
					Database: DatabaseConfig{
						Host:     "db.example.com",
						Port:     5433,
						User:     "testuser",
						Password: "testpass",
						Name:     "testdb",
						SSLMode:  "require",
					},
					Cache: CacheConfig{
						Enabled:  true,
						Host:     "cache.example.com",
						Port:     6380,
						Password: "cachepass",
						DB:       1,
					},
				}

				data, _ := json.Marshal(testConfig)
				tmpFile := filepath.Join(os.TempDir(), "test_config.json")
				_ = os.WriteFile(tmpFile, data, 0644)
				return tmpFile
			},
			expectError: false,
			validate: func(config *Config, t *testing.T) {
				if config.Server.Port != 9090 {
					t.Errorf("Server port = %d, want %d", config.Server.Port, 9090)
				}
				if config.Server.Host != "127.0.0.1" {
					t.Errorf("Server host = %s, want %s", config.Server.Host, "127.0.0.1")
				}
				if config.Server.BaseURL != "http://example.com" {
					t.Errorf("Server base URL = %s, want %s", config.Server.BaseURL, "http://example.com")
				}
				if config.Database.Host != "db.example.com" {
					t.Errorf("Database host = %s, want %s", config.Database.Host, "db.example.com")
				}
				if config.Database.Port != 5433 {
					t.Errorf("Database port = %d, want %d", config.Database.Port, 5433)
				}
				if config.Database.User != "testuser" {
					t.Errorf("Database user = %s, want %s", config.Database.User, "testuser")
				}
				if config.Database.Password != "testpass" {
					t.Errorf("Database password = %s, want %s", config.Database.Password, "testpass")
				}
				if config.Database.Name != "testdb" {
					t.Errorf("Database name = %s, want %s", config.Database.Name, "testdb")
				}
				if config.Database.SSLMode != "require" {
					t.Errorf("Database SSL mode = %s, want %s", config.Database.SSLMode, "require")
				}
				if config.Cache.Enabled != true {
					t.Errorf("Cache enabled = %t, want %t", config.Cache.Enabled, true)
				}
				if config.Cache.Host != "cache.example.com" {
					t.Errorf("Cache host = %s, want %s", config.Cache.Host, "cache.example.com")
				}
				if config.Cache.Port != 6380 {
					t.Errorf("Cache port = %d, want %d", config.Cache.Port, 6380)
				}
				if config.Cache.Password != "cachepass" {
					t.Errorf("Cache password = %s, want %s", config.Cache.Password, "cachepass")
				}
				if config.Cache.DB != 1 {
					t.Errorf("Cache DB = %d, want %d", config.Cache.DB, 1)
				}
			},
		},
		{
			name:       "non-existent config file",
			configPath: "non_existent_config.json",
			setupConfig: func() string {
				return "non_existent_config.json"
			},
			expectError: true,
			validate:    func(config *Config, t *testing.T) {},
		},
		{
			name:       "invalid JSON in config file",
			configPath: "invalid_config.json",
			setupConfig: func() string {
				// Create an invalid JSON file
				tmpFile := filepath.Join(os.TempDir(), "invalid_config.json")
				_ = os.WriteFile(tmpFile, []byte("invalid json"), 0644)
				return tmpFile
			},
			expectError: true,
			validate:    func(config *Config, t *testing.T) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			configPath := tc.setupConfig()
			if configPath != "" && configPath != tc.configPath {
				defer os.Remove(configPath)
				tc.configPath = configPath
			}

			// Test
			config, err := LoadConfig(tc.configPath)

			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if config == nil {
					t.Errorf("Expected config, got nil")
				} else {
					tc.validate(config, t)
				}
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Save original environment variables
	origEnv := make(map[string]string)
	envVars := []string{
		"TT_SERVER_PORT", "TT_SERVER_HOST", "TT_SERVER_BASE_URL",
		"TT_DB_HOST", "TT_DB_PORT", "TT_DB_USER", "TT_DB_PASSWORD", "TT_DB_NAME", "TT_DB_SSL_MODE",
		"TT_CACHE_ENABLED", "TT_CACHE_HOST", "TT_CACHE_PORT", "TT_CACHE_PASSWORD", "TT_CACHE_DB",
	}
	for _, env := range envVars {
		origEnv[env] = os.Getenv(env)
	}

	// Restore environment variables after test
	defer func() {
		for env, val := range origEnv {
			if val == "" {
				os.Unsetenv(env)
			} else {
				os.Setenv(env, val)
			}
		}
	}()

	// Clear all environment variables
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Test with custom environment variables
	os.Setenv("TT_SERVER_PORT", "9090")
	os.Setenv("TT_SERVER_HOST", "127.0.0.1")
	os.Setenv("TT_SERVER_BASE_URL", "http://example.com")
	os.Setenv("TT_DB_HOST", "db.example.com")
	os.Setenv("TT_DB_PORT", "5433")
	os.Setenv("TT_DB_USER", "testuser")
	os.Setenv("TT_DB_PASSWORD", "testpass")
	os.Setenv("TT_DB_NAME", "testdb")
	os.Setenv("TT_DB_SSL_MODE", "require")
	os.Setenv("TT_CACHE_ENABLED", "true")
	os.Setenv("TT_CACHE_HOST", "cache.example.com")
	os.Setenv("TT_CACHE_PORT", "6380")
	os.Setenv("TT_CACHE_PASSWORD", "cachepass")
	os.Setenv("TT_CACHE_DB", "1")

	// Test
	config := LoadConfigFromEnv()

	// Assert
	if config.Server.Port != 9090 {
		t.Errorf("Server port = %d, want %d", config.Server.Port, 9090)
	}
	if config.Server.Host != "127.0.0.1" {
		t.Errorf("Server host = %s, want %s", config.Server.Host, "127.0.0.1")
	}
	if config.Server.BaseURL != "http://example.com" {
		t.Errorf("Server base URL = %s, want %s", config.Server.BaseURL, "http://example.com")
	}
	if config.Database.Host != "db.example.com" {
		t.Errorf("Database host = %s, want %s", config.Database.Host, "db.example.com")
	}
	if config.Database.Port != 5433 {
		t.Errorf("Database port = %d, want %d", config.Database.Port, 5433)
	}
	if config.Database.User != "testuser" {
		t.Errorf("Database user = %s, want %s", config.Database.User, "testuser")
	}
	if config.Database.Password != "testpass" {
		t.Errorf("Database password = %s, want %s", config.Database.Password, "testpass")
	}
	if config.Database.Name != "testdb" {
		t.Errorf("Database name = %s, want %s", config.Database.Name, "testdb")
	}
	if config.Database.SSLMode != "require" {
		t.Errorf("Database SSL mode = %s, want %s", config.Database.SSLMode, "require")
	}
	if config.Cache.Enabled != true {
		t.Errorf("Cache enabled = %t, want %t", config.Cache.Enabled, true)
	}
	if config.Cache.Host != "cache.example.com" {
		t.Errorf("Cache host = %s, want %s", config.Cache.Host, "cache.example.com")
	}
	if config.Cache.Port != 6380 {
		t.Errorf("Cache port = %d, want %d", config.Cache.Port, 6380)
	}
	if config.Cache.Password != "cachepass" {
		t.Errorf("Cache password = %s, want %s", config.Cache.Password, "cachepass")
	}
	if config.Cache.DB != 1 {
		t.Errorf("Cache DB = %d, want %d", config.Cache.DB, 1)
	}

	// Test with invalid port values
	os.Setenv("TT_SERVER_PORT", "invalid")
	os.Setenv("TT_DB_PORT", "invalid")
	os.Setenv("TT_CACHE_PORT", "invalid")
	os.Setenv("TT_CACHE_DB", "invalid")

	// Test
	config = LoadConfigFromEnv()

	// Assert - should use default values for invalid ports
	if config.Server.Port != 8080 {
		t.Errorf("Server port = %d, want %d", config.Server.Port, 8080)
	}
	if config.Database.Port != 5432 {
		t.Errorf("Database port = %d, want %d", config.Database.Port, 5432)
	}
	if config.Cache.Port != 6379 {
		t.Errorf("Cache port = %d, want %d", config.Cache.Port, 6379)
	}
	if config.Cache.DB != 0 {
		t.Errorf("Cache DB = %d, want %d", config.Cache.DB, 0)
	}
}
