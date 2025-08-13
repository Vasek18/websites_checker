package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"website-monitor/internal/models"
	"gopkg.in/yaml.v3"
)

// Load loads configuration from .env file and URLs from yaml/json file
func Load() (*models.Config, error) {
	// Load environment variables from .env file
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

// loadDatabaseConfig loads database configuration from .env file
func loadDatabaseConfig() (*models.DatabaseConfig, error) {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file not found")
	}

	env := make(map[string]string)
	file, err := os.Open(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	// Validate required environment variables
	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, key := range required {
		if _, exists := env[key]; !exists {
			return nil, fmt.Errorf("required environment variable %s not found", key)
		}
	}

	return &models.DatabaseConfig{
		Host:     env["DB_HOST"],
		Port:     env["DB_PORT"],
		User:     env["DB_USER"],
		Password: env["DB_PASSWORD"],
		Name:     env["DB_NAME"],
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