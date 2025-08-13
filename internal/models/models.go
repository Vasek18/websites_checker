package models

import "time"

// MonitoredURL represents a URL to be monitored
type MonitoredURL struct {
	ID                int    `json:"id"`
	URL               string `json:"url"`
	CheckIntervalSec  int    `json:"check_interval_sec"`
	RegexPattern      string `json:"regex_pattern,omitempty"`
}

// CheckResult represents the result of a website check
type CheckResult struct {
	ID               int       `json:"id"`
	URL              string    `json:"url"`
	CheckTimestamp   time.Time `json:"check_timestamp"`
	ResponseTimeMs   *int      `json:"response_time_ms,omitempty"`
	HTTPStatus       *int      `json:"http_status,omitempty"`
	RegexMatch       *bool     `json:"regex_match,omitempty"`
	Error            string    `json:"error,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig `json:"database"`
	URLs     []MonitoredURL `json:"urls"`
}

// DatabaseConfig holds database connection parameters
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}