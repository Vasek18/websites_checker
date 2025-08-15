package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Connect to database
	database, err := connectToDatabase()
	if err != nil {
		return err
	}
	defer database.Close()

	// Create and start scheduler
	sched, ctx, cancel, err := setupScheduler(database)
	if err != nil {
		return err
	}
	defer cancel()

	// Set up signal handling and wait for shutdown
	return waitForShutdown(ctx, cancel, sched)
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

func setupScheduler(database *db.DB) (*scheduler.Scheduler, context.Context, context.CancelFunc, error) {
	// Create repository
	repo := url_repository.New(database)

	// Create scheduler
	sched := scheduler.New(repo, database)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Start scheduler
	if err := sched.Start(ctx); err != nil {
		cancel()
		log.Printf("Failed to start scheduler: %v", err)
		return nil, nil, nil, err
	}

	log.Println("Scheduler started successfully")
	return sched, ctx, cancel, nil
}

func waitForShutdown(ctx context.Context, cancel context.CancelFunc, sched *scheduler.Scheduler) error {
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

type Stoppable interface {
	Stop()
}

func performGracefulShutdown(cancel context.CancelFunc, sched Stoppable) error {
	log.Println("Starting graceful shutdown...")

	// Cancel context to signal all goroutines to stop
	cancel()

	// Stop scheduler and wait for all goroutines to finish
	sched.Stop()

	log.Println("Website Monitor stopped")
	return nil
}
