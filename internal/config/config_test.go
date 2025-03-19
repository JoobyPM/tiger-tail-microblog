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
	if config.Database.Port != "5432" {
		t.Errorf("Default database port = %s, want %s", config.Database.Port, "5432")
	}
	if string(config.Database.User) != "postgres" {
		t.Errorf("Default database user = %s, want %s", config.Database.User, "postgres")
	}
	if string(config.Database.Password) != "postgres" {
		t.Errorf("Default database password = %s, want %s", config.Database.Password, "postgres")
	}
	if config.Database.Name != "tigertail" {
		t.Errorf("Default database name = %s, want %s", config.Database.Name, "tigertail")
	}
	if config.Database.SSLMode != "disable" {
		t.Errorf("Default database SSL mode = %s, want %s", config.Database.SSLMode, "disable")
	}

	// Verify default cache config
	if config.Cache.Host != "localhost" {
		t.Errorf("Default cache host = %s, want %s", config.Cache.Host, "localhost")
	}
	if config.Cache.Port != "6379" {
		t.Errorf("Default cache port = %s, want %s", config.Cache.Port, "6379")
	}
	if string(config.Cache.Password) != "" {
		t.Errorf("Default cache password = %s, want %s", config.Cache.Password, "")
	}
	if config.Cache.DB != 0 {
		t.Errorf("Default cache DB = %d, want %d", config.Cache.DB, 0)
	}
	
	// Verify default auth config
	if string(config.Auth.Username) != "admin" {
		t.Errorf("Default auth username = %s, want %s", config.Auth.Username, "admin")
	}
	if string(config.Auth.Password) != "password" {
		t.Errorf("Default auth password = %s, want %s", config.Auth.Password, "password")
	}
	
	// Verify default use real services
	if config.UseRealDB != false {
		t.Errorf("Default use real DB = %t, want %t", config.UseRealDB, false)
	}
	if config.UseRealCache != false {
		t.Errorf("Default use real cache = %t, want %t", config.UseRealCache, false)
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
					Database: DatabaseCredentials{
						Host:     "db.example.com",
						Port:     "5433",
						User:     SensitiveString("testuser"),
						Password: SensitiveString("testpass"),
						Name:     "testdb",
						SSLMode:  "require",
					},
					Cache: RedisCredentials{
						Host:     "cache.example.com",
						Port:     "6380",
						Password: SensitiveString("cachepass"),
						DB:       1,
					},
					Auth: AuthCredentials{
						Username: SensitiveString("testadmin"),
						Password: SensitiveString("testadminpass"),
					},
					UseRealDB:    true,
					UseRealCache: true,
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
				if config.Database.Port != "5433" {
					t.Errorf("Database port = %s, want %s", config.Database.Port, "5433")
				}
				if string(config.Database.User) != "testuser" {
					t.Errorf("Database user = %s, want %s", config.Database.User, "testuser")
				}
				if string(config.Database.Password) != "testpass" {
					t.Errorf("Database password = %s, want %s", config.Database.Password, "testpass")
				}
				if config.Database.Name != "testdb" {
					t.Errorf("Database name = %s, want %s", config.Database.Name, "testdb")
				}
				if config.Database.SSLMode != "require" {
					t.Errorf("Database SSL mode = %s, want %s", config.Database.SSLMode, "require")
				}
				if config.Cache.Host != "cache.example.com" {
					t.Errorf("Cache host = %s, want %s", config.Cache.Host, "cache.example.com")
				}
				if config.Cache.Port != "6380" {
					t.Errorf("Cache port = %s, want %s", config.Cache.Port, "6380")
				}
				if string(config.Cache.Password) != "cachepass" {
					t.Errorf("Cache password = %s, want %s", config.Cache.Password, "cachepass")
				}
				if config.Cache.DB != 1 {
					t.Errorf("Cache DB = %d, want %d", config.Cache.DB, 1)
				}
				if string(config.Auth.Username) != "testadmin" {
					t.Errorf("Auth username = %s, want %s", config.Auth.Username, "testadmin")
				}
				if string(config.Auth.Password) != "testadminpass" {
					t.Errorf("Auth password = %s, want %s", config.Auth.Password, "testadminpass")
				}
				if config.UseRealDB != true {
					t.Errorf("Use real DB = %t, want %t", config.UseRealDB, true)
				}
				if config.UseRealCache != true {
					t.Errorf("Use real cache = %t, want %t", config.UseRealCache, true)
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
		"SERVER_PORT", "SERVER_HOST", "SERVER_BASE_URL",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB",
		"AUTH_USERNAME", "AUTH_PASSWORD", "USE_REAL_DB", "USE_REAL_REDIS",
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
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_BASE_URL", "http://example.com")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("REDIS_HOST", "cache.example.com")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_PASSWORD", "cachepass")
	os.Setenv("REDIS_DB", "1")
	os.Setenv("AUTH_USERNAME", "testadmin")
	os.Setenv("AUTH_PASSWORD", "testadminpass")
	os.Setenv("USE_REAL_DB", "true")
	os.Setenv("USE_REAL_REDIS", "true")

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
	if config.Database.Port != "5433" {
		t.Errorf("Database port = %s, want %s", config.Database.Port, "5433")
	}
	if string(config.Database.User) != "testuser" {
		t.Errorf("Database user = %s, want %s", config.Database.User, "testuser")
	}
	if string(config.Database.Password) != "testpass" {
		t.Errorf("Database password = %s, want %s", config.Database.Password, "testpass")
	}
	if config.Database.Name != "testdb" {
		t.Errorf("Database name = %s, want %s", config.Database.Name, "testdb")
	}
	if config.Database.SSLMode != "require" {
		t.Errorf("Database SSL mode = %s, want %s", config.Database.SSLMode, "require")
	}
	if config.Cache.Host != "cache.example.com" {
		t.Errorf("Cache host = %s, want %s", config.Cache.Host, "cache.example.com")
	}
	if config.Cache.Port != "6380" {
		t.Errorf("Cache port = %s, want %s", config.Cache.Port, "6380")
	}
	if string(config.Cache.Password) != "cachepass" {
		t.Errorf("Cache password = %s, want %s", config.Cache.Password, "cachepass")
	}
	if config.Cache.DB != 1 {
		t.Errorf("Cache DB = %d, want %d", config.Cache.DB, 1)
	}
	if string(config.Auth.Username) != "testadmin" {
		t.Errorf("Auth username = %s, want %s", config.Auth.Username, "testadmin")
	}
	if string(config.Auth.Password) != "testadminpass" {
		t.Errorf("Auth password = %s, want %s", config.Auth.Password, "testadminpass")
	}
	if config.UseRealDB != true {
		t.Errorf("Use real DB = %t, want %t", config.UseRealDB, true)
	}
	if config.UseRealCache != true {
		t.Errorf("Use real cache = %t, want %t", config.UseRealCache, true)
	}

	// Test with invalid port values
	os.Setenv("SERVER_PORT", "invalid")
	os.Setenv("REDIS_DB", "invalid")

	// Test
	config = LoadConfigFromEnv()

	// Assert - should use default values for invalid ports
	if config.Server.Port != 8080 {
		t.Errorf("Server port = %d, want %d", config.Server.Port, 8080)
	}
	if config.Cache.DB != 0 {
		t.Errorf("Cache DB = %d, want %d", config.Cache.DB, 0)
	}
}
