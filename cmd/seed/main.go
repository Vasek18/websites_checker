package main

import (
	"log"

	"website-monitor/internal/db"
	"website-monitor/internal/models"
)

// todo remove the file in case of release
func main() {
	log.Println("Seeding database with mock URLs...")

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(database *db.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(database)

	mockURLs := []models.MonitoredUrl{
		{
			Url:              "https://aiven.io/",
			CheckIntervalSec: 5,
			RegexPattern:     "",
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
		{
			Url:              "https://9gag.com/",
			CheckIntervalSec: 60,
			RegexPattern:     "",
		},
		{
			Url:              "https://jsonplaceholder.typicode.com/posts/1",
			CheckIntervalSec: 75,
			RegexPattern:     `"userId":\s*1`,
		},
	}

	upsertQuery := `
		INSERT INTO monitored_urls (url, check_interval_sec, regex_pattern)
		VALUES ($1, $2, $3)
		ON CONFLICT (url) DO UPDATE SET
			check_interval_sec = EXCLUDED.check_interval_sec,
			regex_pattern = EXCLUDED.regex_pattern`

	for _, url := range mockURLs {
		if err := database.Exec(upsertQuery, url.Url, url.CheckIntervalSec, url.RegexPattern); err != nil {
			log.Printf("Failed to insert url %s: %v", url.Url, err)
		} else {
			log.Printf("Successfully seeded url: %s (interval: %ds)", url.Url, url.CheckIntervalSec)
		}
	}

	log.Println("Database seeding completed")
}
