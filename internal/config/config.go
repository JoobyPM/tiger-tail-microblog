package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Cache    CacheConfig    `json:"cache"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port    int    `json:"port"`
	Host    string `json:"host"`
	BaseURL string `json:"base_url"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

// CacheConfig represents the cache configuration
type CacheConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    8080,
			Host:    "0.0.0.0",
			BaseURL: "http://localhost:8080",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "tigertail",
			SSLMode:  "disable",
		},
		Cache: CacheConfig{
			Enabled:  false,
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
	}
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	// Use default config as a base
	config := DefaultConfig()

	// If path is empty, return default config
	if path == "" {
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse config file
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

// LoadConfigFromEnv loads the configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// Server config
	if port := os.Getenv("TT_SERVER_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Server.Port)
	}
	if host := os.Getenv("TT_SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if baseURL := os.Getenv("TT_SERVER_BASE_URL"); baseURL != "" {
		config.Server.BaseURL = baseURL
	}

	// Database config
	if host := os.Getenv("TT_DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("TT_DB_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Database.Port)
	}
	if user := os.Getenv("TT_DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("TT_DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if name := os.Getenv("TT_DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if sslMode := os.Getenv("TT_DB_SSL_MODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Cache config
	if enabled := os.Getenv("TT_CACHE_ENABLED"); enabled == "true" {
		config.Cache.Enabled = true
	}
	if host := os.Getenv("TT_CACHE_HOST"); host != "" {
		config.Cache.Host = host
	}
	if port := os.Getenv("TT_CACHE_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Cache.Port)
	}
	if password := os.Getenv("TT_CACHE_PASSWORD"); password != "" {
		config.Cache.Password = password
	}
	if db := os.Getenv("TT_CACHE_DB"); db != "" {
		fmt.Sscanf(db, "%d", &config.Cache.DB)
	}

	return config
}
