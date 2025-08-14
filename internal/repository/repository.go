package repository

import "website-monitor/internal/models"

// URLRepository defines the interface for url data sources
type URLRepository interface {
	// GetMonitoredURLs returns all URLs that should be monitored
	GetMonitoredURLs() ([]models.MonitoredUrl, error)
}
