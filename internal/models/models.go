package models

import (
	"time"
)

// MonitoredUrl represents a url to be monitored
type MonitoredUrl struct {
	ID               int    `json:"id"`
	Url              string `json:"url"`
	CheckIntervalSec int    `json:"check_interval_sec"`
	RegexPattern     string `json:"regex_pattern,omitempty"`
}

// CheckResult represents the result of a website check
type CheckResult struct {
	ID             int       `json:"id"`
	URL            string    `json:"url"`
	CheckTimestamp time.Time `json:"check_timestamp"`
	ResponseTimeMs *int      `json:"response_time_ms,omitempty"` // todo why pointer?
	HTTPStatus     *int      `json:"http_status,omitempty"`      // todo why pointer?
	RegexMatch     *bool     `json:"regex_match,omitempty"`      // todo why pointer?
	Error          string    `json:"error,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig `json:"database"`
}

// DatabaseConfig holds database connection parameters
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
