package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"website-monitor/internal/checker"
	"website-monitor/internal/db"
	"website-monitor/internal/scheduler"
	"website-monitor/internal/url_repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	log.Println("Starting Website Monitor...")

	database, err := connectToDatabase()
	if err != nil {
		return err
	}
	defer func(database *db.DB) {
		err := database.Close()
		if err != nil {
			fmt.Printf("Error closing database connection: %v", err)
		}
	}(database)

	sched, cancel, err := setupScheduler(database)
	if err != nil {
		return err
	}
	defer cancel()

	// Set up signal handling and wait for shutdown
	return waitForShutdown(cancel, sched)
}

func connectToDatabase() (*db.DB, error) {
	database, err := db.Connect()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)

		return nil, err
	}

	log.Println("Connected to database successfully")

	return database, nil
}

func setupScheduler(database *db.DB) (*scheduler.Scheduler, context.CancelFunc, error) {
	repo := url_repository.New(database)
	chk := checker.New(database)
	sched := scheduler.New(repo, database, chk)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	if err := sched.Start(ctx); err != nil {
		cancel()
		log.Printf("Failed to start scheduler: %v", err)

		return nil, nil, err
	}

	log.Println("Scheduler started successfully")

	return sched, cancel, nil
}

func waitForShutdown(cancel context.CancelFunc, sched *scheduler.Scheduler) error {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Website Monitor is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal")

	// Graceful shutdown
	return performGracefulShutdown(cancel, sched)
}

func performGracefulShutdown(cancel context.CancelFunc, sched scheduler.Stoppable) error {
	// Cancel context to signal all goroutines to stop
	cancel()

	// Stop scheduler and wait for all goroutines to finish
	sched.Stop()

	log.Println("Website Monitor stopped")

	return nil
}
