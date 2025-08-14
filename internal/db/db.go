package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"website-monitor/internal/config"
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

// Exec executes a query with parameters and returns an error if it fails
func (db *DB) Exec(query string, args ...interface{}) error {
	_, err := db.conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// Query executes a query and returns rows
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return rows, nil
}

// GetUrl returns a database connection url for migration tools
func GetUrl() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name), nil
}
