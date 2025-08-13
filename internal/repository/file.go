package repository

import (
	"website-monitor/internal/config"
	"website-monitor/internal/models"
)

// FileRepository implements URLRepository by reading from configuration files
type FileRepository struct {
	urls []models.MonitoredURL
}

// NewFileRepository creates a new file-based repository
func NewFileRepository() (*FileRepository, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	
	return &FileRepository{
		urls: cfg.URLs,
	}, nil
}

// GetMonitoredURLs returns all URLs loaded from the configuration file
func (r *FileRepository) GetMonitoredURLs() ([]models.MonitoredURL, error) {
	return r.urls, nil
}