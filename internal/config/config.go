package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load loads configuration from environment variables and URLs from yaml/json file
func Load() (*models.Config, error) {
	// Load database configuration from environment variables
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Load URLs from urls.yaml or urls.json
	urls, err := loadURLsConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load URLs config: %w", err)
	}

	return &models.Config{
		Database: *dbConfig,
		URLs:     urls,
	}, nil
}

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig() (*models.DatabaseConfig, error) {
	// Get required environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

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
	}, nil
}

// URLsConfig represents the structure of the URLs configuration file
type URLsConfig struct {
	URLs []models.MonitoredURL `yaml:"urls" json:"urls"`
}

// loadURLsConfig loads URLs configuration from urls.yaml or urls.json
func loadURLsConfig() ([]models.MonitoredURL, error) {
	// Try urls.yaml first, then urls.json
	configFiles := []string{"urls.yaml", "urls.yml", "urls.json"}

	for _, filename := range configFiles {
		if _, err := os.Stat(filename); err == nil {
			return loadURLsFromFile(filename)
		}
	}

	return nil, fmt.Errorf("no URLs configuration file found (tried: %s)", strings.Join(configFiles, ", "))
}

// loadURLsFromFile loads URLs from a specific file based on its extension
func loadURLsFromFile(filename string) ([]models.MonitoredURL, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	var config URLsConfig
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML file %s: %w", filename, err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON file %s: %w", filename, err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	// Validate URLs configuration
	for i, url := range config.URLs {
		if url.URL == "" {
			return nil, fmt.Errorf("URL at index %d is empty", i)
		}
		if url.CheckIntervalSec < 5 || url.CheckIntervalSec > 300 {
			return nil, fmt.Errorf("check interval for URL %s must be between 5 and 300 seconds", url.URL)
		}
	}

	return config.URLs, nil
}
