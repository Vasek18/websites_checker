package main

import (
	"log"

	"website-monitor/internal/db"
	"website-monitor/internal/models"
)

func main() {
	log.Println("Seeding database with mock URLs...")

	// Connect to database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Define mock urls to seed
	mockURLs := []models.MonitoredUrl{
		{
			Url:              "https://aiven.io/",
			CheckIntervalSec: 5,
			RegexPattern:     "Example Domain",
		},
		{
			Url:              "https://www.google.com",
			CheckIntervalSec: 90,
			RegexPattern:     "Google",
		},
		{
			Url:              "https://github.com",
			CheckIntervalSec: 120,
			RegexPattern:     "GitHub",
		},
		{
			Url:              "https://stackoverflow.com",
			CheckIntervalSec: 180,
			RegexPattern:     "Stack Overflow",
		},
		{
			Url:              "https://httpbin.org/status/418",
			CheckIntervalSec: 300,
			RegexPattern:     "",
		},
	}

	// Insert each url into the database
	for _, url := range mockURLs {
		if err := database.UpsertMonitoredURL(url); err != nil {
			log.Printf("Failed to insert url %s: %v", url.Url, err)
		} else {
			log.Printf("Successfully seeded url: %s (interval: %ds)", url.Url, url.CheckIntervalSec)
		}
	}

	log.Println("Database seeding completed successfully")
}
