package url_repository

import (
	"database/sql"
	"fmt"
	"website-monitor/internal/db"
	"website-monitor/internal/models"
)

// DbUrlRepository implements UrlRepository using database as the data source
type DbUrlRepository struct {
	db *db.DB
}

func New(database *db.DB) *DbUrlRepository {
	return &DbUrlRepository{
		db: database,
	}
}

// GetMonitoredUrls returns all URLs that should be monitored from the database
func (r *DbUrlRepository) GetMonitoredUrls() ([]models.MonitoredUrl, error) {
	query := `SELECT id, url, check_interval_sec, COALESCE(regex_pattern, '') FROM monitored_urls`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query monitored URLs: %w", err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

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
