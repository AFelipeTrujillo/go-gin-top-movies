package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// DefaultPort is the default HTTP server port.
	DefaultPort = 8080
	// DefaultDBPath is the default SQLite database file path.
	DefaultDBPath = "internal/Infrastructure/Persistence/movies.sqlite"
	// DefaultPageSize is the default number of items per page.
	DefaultPageSize = 10
	// DefaultMaxPageSize is the maximum allowed page size to prevent abuse.
	DefaultMaxPageSize = 100
)

// Config holds the application configuration.
// Values are loaded from environment variables with sensible defaults.
type Config struct {
	// Port is the HTTP server port.
	Port int
	// DBPath is the filesystem path to the SQLite database.
	DBPath string
	// PageSize is the default number of items per page for paginated endpoints.
	PageSize int
	// MaxPageSize is the maximum allowed page size.
	MaxPageSize int
}

// Load reads configuration from environment variables.
// If a variable is not set or invalid, the default value is used.
func Load() *Config {
	return &Config{
		Port:        getEnvInt("PORT", DefaultPort),
		DBPath:      getEnvString("DB_PATH", DefaultDBPath),
		PageSize:    getEnvInt("PAGE_SIZE", DefaultPageSize),
		MaxPageSize: getEnvInt("MAX_PAGE_SIZE", DefaultMaxPageSize),
	}
}

// Address returns the server address string (e.g., ":8080").
func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

// getEnvString returns the value of an environment variable or a default value.
func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns the integer value of an environment variable or a default value.
func getEnvInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists || valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
