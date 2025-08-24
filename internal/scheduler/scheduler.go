package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"website-monitor/internal/checker"
	"website-monitor/internal/db"
	"website-monitor/internal/models"
	"website-monitor/internal/url_repository"
)

// Stoppable defines the interface for components that can be gracefully stopped
type Stoppable interface {
	Stop()
}

type Scheduler struct {
	repo    url_repository.UrlRepository
	db      *db.DB
	checker checker.IChecker
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func New(repo url_repository.UrlRepository, database *db.DB, chk checker.IChecker) *Scheduler {
	return &Scheduler{
		repo:    repo,
		db:      database,
		checker: chk,
	}
}

// Start begins monitoring of all URLs from the repository
func (s *Scheduler) Start(ctx context.Context) error {
	urls, err := s.repo.GetMonitoredUrls()
	if err != nil {
		return err
	}

	if len(urls) == 0 {
		log.Println("No URLs to monitor")

		return nil
	}

	// Wrap context with cancel to ensure Stop() can immediately signal all goroutines and wait for them to exit
	ctx, s.cancel = context.WithCancel(ctx)

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

	if err := s.checker.InsertCheckResult(result); err != nil {
		log.Printf("Failed to store check result for %s: %v", url.Url, err)
	}
}
