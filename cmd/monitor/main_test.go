package main

import (
	"context"
	"testing"
	"website-monitor/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
)

type mockScheduler struct {
	stopped bool
}

func (m *mockScheduler) Stop() {
	m.stopped = true
}

func TestSetupScheduler_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock db: %v", err)
	}
	defer sqlDB.Close()

	mock.ExpectQuery(`SELECT id, url, check_interval_sec, COALESCE\(regex_pattern, ''\) FROM monitored_urls`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url", "check_interval_sec", "regex_pattern"}))

	sched, cancel, err := setupScheduler(db.New(sqlDB))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if cancel == nil {
		t.Error("Expected cancel function to be created, got nil")
	}
	if sched == nil {
		t.Error("Expected scheduler to be created, got nil")
	}
}

func TestSetupScheduler_NilDatabase(t *testing.T) {
	// We expect this to panic, so we'll catch it
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected setupScheduler to panic with nil database")
		}
	}()

	_, _, _ = setupScheduler(nil)
}

func TestPerformGracefulShutdown(t *testing.T) {
	// Test graceful shutdown with mock scheduler
	mock := &mockScheduler{}
	ctx, cancel := context.WithCancel(context.Background())

	err := performGracefulShutdown(cancel, mock)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !mock.stopped {
		t.Error("Expected scheduler to be stopped")
	}

	// Check that context was cancelled
	select {
	case <-ctx.Done():
		// Context was cancelled as expected
	default:
		t.Error("Expected context to be cancelled")
	}
}
