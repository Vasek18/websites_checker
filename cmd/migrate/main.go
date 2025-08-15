package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"website-monitor/internal/db"
)

func main() {
	log.Println("Running database migrations...")

	dbURL, err := db.GetUrl()
	if err != nil {
		log.Fatalf("Failed to get database url: %v", err)
	}

	m, err := migrate.New("file://internal/migrations", dbURL)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			log.Fatalf("Could not close migrate database: %v", err)
		}
	}(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not run migrations: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Printf("Could not get migration version: %v", err)
	} else {
		log.Printf("Current migration version: %d, dirty: %v", version, dirty)
	}

	log.Println("Database migrations completed successfully")
}
