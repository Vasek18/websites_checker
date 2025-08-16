package config

import (
	"fmt"
	"os"

	"website-monitor/internal/models"
)

// Load loads configuration from environment variables
func Load() (*models.Config, error) {
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	return &models.Config{
		Database: *dbConfig,
	}, nil
}

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig() (*models.DatabaseConfig, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")

	// Set default SSL mode if not specified
	if sslMode == "" {
		sslMode = "require"
	}

	// Validate required environment variables
	required := map[string]string{
		"DB_HOST":     host,
		"DB_PORT":     port,
		"DB_USER":     user,
		"DB_PASSWORD": password,
		"DB_NAME":     name,
	}

	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("required environment variable %s not set", key)
		}
	}

	return &models.DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		SSLMode:  sslMode,
	}, nil
}
