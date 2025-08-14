package repository

import (
	"website-monitor/internal/db"
	"website-monitor/internal/models"
)

// DBRepository implements URLRepository using database as the data source
type DBRepository struct { // todo DbRepository
	db *db.DB
}

// NewDBRepository creates a new database-backed repository
func NewDBRepository(database *db.DB) *DBRepository {
	return &DBRepository{
		db: database,
	}
}

// GetMonitoredURLs returns all URLs that should be monitored from the database
func (r *DBRepository) GetMonitoredURLs() ([]models.MonitoredUrl, error) {
	return r.db.GetMonitoredURLs()
}
