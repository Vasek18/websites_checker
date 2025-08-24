package url_repository

import "website-monitor/internal/models"

// UrlRepository defines the interface for url data sources
type UrlRepository interface {
	// GetMonitoredUrls returns all urls that should be monitored
	GetMonitoredUrls() ([]models.MonitoredUrl, error)
}
