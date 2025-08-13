package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"website-monitor/internal/config"
	"website-monitor/internal/db"
	"website-monitor/internal/repository"
	"website-monitor/internal/scheduler"
)

func main() {
	log.Println("Starting Website Monitor...")

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

	// Create repository
	repo, err := repository.NewFileRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create and start scheduler
	sched := scheduler.New(repo, database)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := sched.Start(ctx); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Website Monitor is running. Press Ctrl+C to stop.")
	
	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal")

	// Graceful shutdown
	cancel()
	sched.Stop()
	
	log.Println("Website Monitor stopped")
}