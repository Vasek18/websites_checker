package main

import (
	"log"

	"website-monitor/internal/config"
	"website-monitor/internal/db"
)

func main() {
	log.Println("Running database migrations...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	database, err := db.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create database tables
	if err := database.CreateTables(); err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}

	log.Println("Database migrations completed successfully")
}