package main

import (
	"context"
	"testing"
)

// Mocks
type mockScheduler struct {
	stopped bool
}

func (m *mockScheduler) Stop() {
	m.stopped = true
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
		t.Errorf("Expected no error from graceful shutdown, got: %v", err)
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
