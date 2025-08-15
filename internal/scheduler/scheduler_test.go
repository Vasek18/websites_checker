package scheduler

import (
	"context"
	"errors"
	"testing"

	"website-monitor/internal/models"
)

type mockRepository struct {
	urls []models.MonitoredUrl
	err  error
}

func (m *mockRepository) GetMonitoredUrls() ([]models.MonitoredUrl, error) {
	return m.urls, m.err
}

func TestScheduler_Start_RepositoryError(t *testing.T) {
	repo := &mockRepository{
		err: errors.New("repository error"),
	}

	scheduler := New(repo, nil)
	ctx := context.Background()

	err := scheduler.Start(ctx)
	if err == nil {
		t.Fatal("Expected error from repository")
	}

	if err.Error() != "repository error" {
		t.Errorf("Expected 'repository error', got: %s", err.Error())
	}
}

// todo happy path
