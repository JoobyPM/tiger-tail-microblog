package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig        `json:"server"`
	Database DatabaseCredentials `json:"database"`
	Cache    RedisCredentials    `json:"cache"`
	Auth     AuthCredentials     `json:"auth"`
	UseRealDB    bool            `json:"use_real_db"`
	UseRealCache bool            `json:"use_real_cache"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port    int    `json:"port"`
	Host    string `json:"host"`
	BaseURL string `json:"base_url"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    8080,
			Host:    "0.0.0.0",
			BaseURL: "http://localhost:8080",
		},
		Database: DatabaseCredentials{
			Host:     "localhost",
			Port:     "5432",
			User:     SensitiveString("postgres"),
			Password: SensitiveString("postgres"),
			Name:     "tigertail",
			SSLMode:  "disable",
		},
		Cache: RedisCredentials{
			Host:     "localhost",
			Port:     "6379",
			Password: SensitiveString(""),
			DB:       0,
		},
		Auth: AuthCredentials{
			Username: SensitiveString("admin"),
			Password: SensitiveString("password"),
		},
		UseRealDB:    false,
		UseRealCache: false,
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
	if port := os.Getenv("SERVER_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &config.Server.Port)
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if baseURL := os.Getenv("SERVER_BASE_URL"); baseURL != "" {
		config.Server.BaseURL = baseURL
	}

	// Database config
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = SensitiveString(user)
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = SensitiveString(password)
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		config.Database.SSLMode = sslMode
	}

	// Cache config
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Cache.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		config.Cache.Port = port
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Cache.Password = SensitiveString(password)
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		fmt.Sscanf(db, "%d", &config.Cache.DB)
	}
	
	// Auth config
	if username := os.Getenv("AUTH_USERNAME"); username != "" {
		config.Auth.Username = SensitiveString(username)
	}
	if password := os.Getenv("AUTH_PASSWORD"); password != "" {
		config.Auth.Password = SensitiveString(password)
	}
	
	// Use real services
	if useRealDB := os.Getenv("USE_REAL_DB"); useRealDB == "true" {
		config.UseRealDB = true
	}
	if useRealCache := os.Getenv("USE_REAL_REDIS"); useRealCache == "true" {
		config.UseRealCache = true
	}

	return config
}
