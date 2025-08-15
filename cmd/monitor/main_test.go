package main

import (
	"context"
	"os"
	"testing"
	"time"

	"website-monitor/internal/models"
)

// Mock repository for testing
type mockRepository struct {
	urls []models.MonitoredUrl
	err  error
}

func (m *mockRepository) GetMonitoredUrls() ([]models.MonitoredUrl, error) {
	return m.urls, m.err
}

// Mock scheduler interface for testing graceful shutdown
type testableScheduler interface {
	Stop()
}

type mockScheduler struct {
	stopped bool
}

func (m *mockScheduler) Stop() {
	m.stopped = true
}

func TestConnectToDatabase_MissingEnv(t *testing.T) {
	// Clear environment variables to simulate missing config
	clearTestEnvVars()

	_, err := connectToDatabase()
	if err == nil {
		t.Error("Expected error when environment variables are missing")
	}
}

func TestConnectToDatabase_ValidEnv(t *testing.T) {
	// Set valid environment variables
	setTestEnvVars()
	defer clearTestEnvVars()

	// This will likely fail without a real database, but should not panic
	_, err := connectToDatabase()
	
	// We expect an error because there's no real database running
	// The important thing is that it doesn't panic and returns an error
	if err == nil {
		t.Log("Unexpected success - you might have a test database running")
	}
}

func TestSetupScheduler_NilDatabase(t *testing.T) {
	// Test with nil database - this will panic because scheduler.New will try to use the database
	// We expect this to panic, so we'll catch it
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected setupScheduler to panic with nil database")
		}
	}()

	_, _, _, _ = setupScheduler(nil)
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

func TestRun_MissingEnvironment(t *testing.T) {
	// Clear environment to test error handling
	clearTestEnvVars()

	err := run()
	if err == nil {
		t.Error("Expected error when environment is not set")
	}
}

func TestRun_ComponentCreation(t *testing.T) {
	// Test that run() properly handles component creation errors
	setTestEnvVars()
	defer clearTestEnvVars()

	// This will fail due to no database, but should handle the error gracefully
	err := run()
	if err == nil {
		t.Log("Unexpected success - you might have a test database running")
	}

	// The important thing is that it returns an error rather than panicking
}

func TestSignalHandlingSetup(t *testing.T) {
	// Test that signal handling can be set up without errors
	// We can't easily test actual signal reception in a unit test
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Signal handling setup panicked: %v", r)
		}
	}()

	// Create a mock scheduler for testing
	mock := &mockScheduler{}

	// Test the signal handling setup (but don't wait for actual signals)
	sigChan := make(chan os.Signal, 1)
	if sigChan == nil {
		t.Error("Failed to create signal channel")
	}

	// Test that we can call performGracefulShutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := performGracefulShutdown(cancel, mock)
	if err != nil {
		t.Errorf("Graceful shutdown failed: %v", err)
	}
}

func TestContextManagement(t *testing.T) {
	// Test context creation and cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Initially, context should not be cancelled
	select {
	case <-ctx.Done():
		t.Error("Context should not be cancelled initially")
	default:
		// Expected - context is not cancelled
	}

	// Cancel the context
	cancel()

	// Now context should be cancelled
	select {
	case <-ctx.Done():
		// Expected - context is cancelled
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled after calling cancel()")
	}
}

func TestMain_Integration(t *testing.T) {
	// Integration test that verifies main components work together
	// This is more of a smoke test to ensure no panics occur

	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test - set INTEGRATION_TEST=true to run")
	}

	// Set up environment
	setTestEnvVars()
	defer clearTestEnvVars()

	// Test that main doesn't panic when run with proper environment
	// We can't test the full main() because it would wait for signals
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Main components caused panic: %v", r)
		}
	}()

	// Test individual components
	t.Run("DatabaseConnection", func(t *testing.T) {
		_, err := connectToDatabase()
		// We expect this to fail without a real database
		if err == nil {
			t.Log("Database connection succeeded - you may have a test database")
		}
	})
}

func TestApplicationFlow(t *testing.T) {
	// Test the overall application flow without actually running it
	tests := []struct {
		name        string
		setupEnv    bool
		expectError bool
	}{
		{
			name:        "Missing environment",
			setupEnv:    false,
			expectError: true,
		},
		{
			name:        "Valid environment",
			setupEnv:    true,
			expectError: true, // Still expect error due to no database
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv {
				setTestEnvVars()
				defer clearTestEnvVars()
			} else {
				clearTestEnvVars()
			}

			err := run()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// Helper functions for test setup
func setTestEnvVars() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
}

func clearTestEnvVars() {
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
}

func TestEnvironmentHelpers(t *testing.T) {
	// Test our helper functions work correctly
	clearTestEnvVars()

	if os.Getenv("DB_HOST") != "" {
		t.Error("Environment should be clear after clearTestEnvVars()")
	}

	setTestEnvVars()

	if os.Getenv("DB_HOST") != "localhost" {
		t.Error("Environment should be set after setTestEnvVars()")
	}

	clearTestEnvVars()

	if os.Getenv("DB_HOST") != "" {
		t.Error("Environment should be clear after second clearTestEnvVars()")
	}
}