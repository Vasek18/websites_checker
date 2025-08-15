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

func TestScheduler_New(t *testing.T) {
	repo := &mockRepository{}
	
	scheduler := New(repo, nil)
	
	if scheduler == nil {
		t.Fatal("Expected non-nil scheduler")
	}
	
	if scheduler.repo != repo {
		t.Error("Expected repository to be set")
	}
	
	if scheduler.checker == nil {
		t.Error("Expected checker to be initialized")
	}
}

func TestScheduler_Start_NoUrls(t *testing.T) {
	repo := &mockRepository{
		urls: []models.MonitoredUrl{},
	}
	
	scheduler := New(repo, nil)
	ctx := context.Background()
	
	err := scheduler.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
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

func TestScheduler_Stop_BeforeStart(t *testing.T) {
	repo := &mockRepository{}
	
	scheduler := New(repo, nil)
	
	scheduler.Stop()
}