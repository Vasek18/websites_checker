package main

import (
	"log"

	"website-monitor/internal/config"
	"website-monitor/internal/db"
	"website-monitor/internal/models"
)

func main() {
	log.Println("Seeding database with mock URLs...")

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

	// Define mock URLs to seed
	mockURLs := []models.MonitoredURL{
		{
			URL:              "https://aiven.io/",
			CheckIntervalSec: 5,
			RegexPattern:     "Example Domain",
		},
		{
			URL:              "https://www.google.com",
			CheckIntervalSec: 90,
			RegexPattern:     "Google",
		},
		{
			URL:              "https://github.com",
			CheckIntervalSec: 120,
			RegexPattern:     "GitHub",
		},
		{
			URL:              "https://stackoverflow.com",
			CheckIntervalSec: 180,
			RegexPattern:     "Stack Overflow",
		},
		{
			URL:              "https://httpbin.org/status/418",
			CheckIntervalSec: 300,
			RegexPattern:     "",
		},
	}

	// Insert each URL into the database
	for _, url := range mockURLs {
		if err := database.UpsertMonitoredURL(url); err != nil {
			log.Printf("Failed to insert URL %s: %v", url.URL, err)
		} else {
			log.Printf("Successfully seeded URL: %s (interval: %ds)", url.URL, url.CheckIntervalSec)
		}
	}

	log.Println("Database seeding completed successfully")
}
