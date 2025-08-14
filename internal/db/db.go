package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"website-monitor/internal/config"
	"website-monitor/internal/models"
)

// DB wraps the database connection and provides methods for database operations
type DB struct {
	conn *sql.DB
}

// Connect creates a new database connection
func Connect() (*DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to database %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// DB returns the underlying *sql.DB connection
func (db *DB) DB() *sql.DB {
	return db.conn
}

// InsertCheckResult inserts a check result into the database
func (db *DB) InsertCheckResult(result models.CheckResult) error {
	query := `
		INSERT INTO checks (url, check_timestamp, response_time_ms, http_status, regex_match, error)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.conn.Exec(query,
		result.URL,
		result.CheckTimestamp,
		result.ResponseTimeMs,
		result.HTTPStatus,
		result.RegexMatch,
		result.Error)

	if err != nil {
		return fmt.Errorf("failed to insert check result: %w", err)
	}

	return nil
}

// UpsertMonitoredURL inserts or updates a monitored URL
func (db *DB) UpsertMonitoredURL(url models.MonitoredURL) error {
	query := `
		INSERT INTO monitored_urls (url, check_interval_sec, regex_pattern)
		VALUES ($1, $2, $3)
		ON CONFLICT (url) DO UPDATE SET
			check_interval_sec = EXCLUDED.check_interval_sec,
			regex_pattern = EXCLUDED.regex_pattern`

	_, err := db.conn.Exec(query, url.URL, url.CheckIntervalSec, url.RegexPattern)
	if err != nil {
		return fmt.Errorf("failed to upsert monitored URL: %w", err)
	}

	return nil
}

// GetMonitoredURLs retrieves all monitored URLs from the database
func (db *DB) GetMonitoredURLs() ([]models.MonitoredURL, error) {
	query := `SELECT id, url, check_interval_sec, COALESCE(regex_pattern, '') FROM monitored_urls`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query monitored URLs: %w", err)
	}
	defer rows.Close()

	var urls []models.MonitoredURL
	for rows.Next() {
		var url models.MonitoredURL
		if err := rows.Scan(&url.ID, &url.URL, &url.CheckIntervalSec, &url.RegexPattern); err != nil {
			return nil, fmt.Errorf("failed to scan monitored URL: %w", err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over monitored URLs: %w", err)
	}

	return urls, nil
}

// GetURL returns a database connection URL for migration tools
func GetURL() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name), nil
}
