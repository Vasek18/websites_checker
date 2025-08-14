package repository

import (
	"fmt"
	"website-monitor/internal/db"
	"website-monitor/internal/models"
)

// DbRepository implements UrlRepository using database as the data source
type DbRepository struct {
	db *db.DB
}

// NewDBRepository creates a new database-backed repository
func NewDBRepository(database *db.DB) *DbRepository {
	return &DbRepository{
		db: database,
	}
}

// GetMonitoredURLs returns all URLs that should be monitored from the database
func (r *DbRepository) GetMonitoredUrls() ([]models.MonitoredUrl, error) {
	query := `SELECT id, url, check_interval_sec, COALESCE(regex_pattern, '') FROM monitored_urls`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query monitored URLs: %w", err)
	}
	defer rows.Close()

	var urls []models.MonitoredUrl
	for rows.Next() {
		var url models.MonitoredUrl
		if err := rows.Scan(&url.ID, &url.Url, &url.CheckIntervalSec, &url.RegexPattern); err != nil {
			return nil, fmt.Errorf("failed to scan monitored url: %w", err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over monitored urls: %w", err)
	}

	return urls, nil
}
