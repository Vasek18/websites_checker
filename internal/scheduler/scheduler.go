package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"website-monitor/internal/checker"
	"website-monitor/internal/db"
	"website-monitor/internal/models"
	"website-monitor/internal/repository"
)

// Scheduler manages the periodic checking of URLs
type Scheduler struct {
	repo    repository.UrlRepository
	db      *db.DB
	checker *checker.HTTPChecker
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// New creates a new scheduler
func New(repo repository.UrlRepository, database *db.DB) *Scheduler {
	return &Scheduler{
		repo:    repo,
		db:      database,
		checker: checker.NewHTTPChecker(),
	}
}

// Start begins monitoring all URLs from the repository
func (s *Scheduler) Start(ctx context.Context) error {
	urls, err := s.repo.GetMonitoredUrls()
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		log.Println("No URLs to monitor")
		return nil
	}

	// Create a context that we can cancel // todo it's already context with cancel
	ctx, s.cancel = context.WithCancel(ctx)

	log.Printf("Starting monitoring for %d URLs", len(urls))

	// Start a goroutine for each url
	for _, url := range urls {
		s.wg.Add(1)
		go s.startMonitorUrl(ctx, url)
	}

	return nil
}

// Stop gracefully stops all monitoring goroutines
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		log.Println("Stopping scheduler...")
		s.cancel()
		s.wg.Wait()
		log.Println("Scheduler stopped")
	}
}

// startMonitorUrl runs in a goroutine to monitor a single url
func (s *Scheduler) startMonitorUrl(ctx context.Context, url models.MonitoredUrl) {
	defer s.wg.Done()

	log.Printf("Starting monitoring for %s (interval: %d seconds)", url.Url, url.CheckIntervalSec)

	ticker := time.NewTicker(time.Duration(url.CheckIntervalSec) * time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		select {
		case <-ctx.Done():
			log.Printf("Stopping monitoring for %s", url.Url)
			return
		default:
		}

		s.performCheck(url)
	}
}

// performCheck executes a single check for a url and stores the result
func (s *Scheduler) performCheck(url models.MonitoredUrl) {
	log.Printf("Checking %s", url.Url)

	result := s.checker.Check(url)

	// Log the result
	if result.Error != "" {
		log.Printf("Check failed for %s: %s", url.Url, result.Error)
	} else {
		status := "unknown"
		if result.HTTPStatus != nil {
			status = fmt.Sprintf("%d", *result.HTTPStatus)
		}
		responseTime := "unknown"
		if result.ResponseTimeMs != nil {
			responseTime = fmt.Sprintf("%dms", *result.ResponseTimeMs)
		}

		regexStatus := "not defined"
		if result.RegexMatch != nil {
			if *result.RegexMatch {
				regexStatus = "match"
			} else {
				regexStatus = "no match"
			}
		}

		log.Printf("Check successful for %s: status=%s, time=%s, regex: %s",
			url.Url, status, responseTime, regexStatus)
	}

	// Store the result in the database
	if err := s.db.InsertCheckResult(result); err != nil {
		log.Printf("Failed to store check result for %s: %v", url.Url, err)
	}
}
