package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"website-monitor/internal/models"
)

type mockRepository struct {
	urls []models.MonitoredUrl
	err  error
}

func (m *mockRepository) GetMonitoredUrls() ([]models.MonitoredUrl, error) {
	return m.urls, m.err
}

type mockChecker struct {
	checkResult     models.CheckResult
	insertError     error
	checkCalls      []models.MonitoredUrl
	insertCalls     []models.CheckResult
	checkCallCount  int
	insertCallCount int
}

func (m *mockChecker) Check(url models.MonitoredUrl) models.CheckResult {
	m.checkCalls = append(m.checkCalls, url)
	m.checkCallCount++
	return m.checkResult
}

func (m *mockChecker) InsertCheckResult(result models.CheckResult) error {
	m.insertCalls = append(m.insertCalls, result)
	m.insertCallCount++
	return m.insertError
}

func TestScheduler_Start_RepositoryError(t *testing.T) {
	repo := &mockRepository{
		err: errors.New("repository error"),
	}
	checker := &mockChecker{}

	scheduler := New(repo, nil, checker)
	ctx := context.Background()

	err := scheduler.Start(ctx)
	if err == nil {
		t.Fatal("Expected error from repository")
	}

	if err.Error() != "repository error" {
		t.Errorf("Expected 'repository error', got: %s", err.Error())
	}
}

func TestScheduler_Start_HappyPath(t *testing.T) {
	urls := []models.MonitoredUrl{
		{
			ID:               1,
			Url:              "https://example.com",
			CheckIntervalSec: 10,
			RegexPattern:     "",
		},
		{
			ID:               2,
			Url:              "https://google.com",
			CheckIntervalSec: 30,
			RegexPattern:     "Google",
		},
	}

	repo := &mockRepository{
		urls: urls,
		err:  nil,
	}
	checker := &mockChecker{}

	scheduler := New(repo, nil, checker)
	ctx := context.Background()

	err := scheduler.Start(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Clean up
	scheduler.Stop()
}

func TestScheduler_Start_EmptyUrls(t *testing.T) {
	repo := &mockRepository{
		urls: []models.MonitoredUrl{},
		err:  nil,
	}
	checker := &mockChecker{}

	scheduler := New(repo, nil, checker)
	ctx := context.Background()

	err := scheduler.Start(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestScheduler_PerformCheck_HappyPath(t *testing.T) {
	url := models.MonitoredUrl{
		ID:               1,
		Url:              "https://example.com",
		CheckIntervalSec: 30,
		RegexPattern:     "Example",
	}

	status := 200
	responseTime := 150
	regexMatch := true

	expectedResult := models.CheckResult{
		URL:            url.Url,
		CheckTimestamp: time.Now(),
		ResponseTimeMs: &responseTime,
		HttpStatus:     &status,
		RegexMatch:     &regexMatch,
		Error:          "",
	}

	checker := &mockChecker{
		checkResult: expectedResult,
		insertError: nil,
	}

	repo := &mockRepository{}
	scheduler := New(repo, nil, checker)

	scheduler.performCheck(url)

	if checker.checkCallCount != 1 {
		t.Errorf("Expected Check to be called once, got %d calls", checker.checkCallCount)
	}

	if checker.insertCallCount != 1 {
		t.Errorf("Expected InsertCheckResult to be called once, got %d calls", checker.insertCallCount)
	}

	if len(checker.checkCalls) != 1 {
		t.Fatalf("Expected 1 check call, got %d", len(checker.checkCalls))
	}

	if checker.checkCalls[0].Url != url.Url {
		t.Errorf("Expected check call with URL %s, got %s", url.Url, checker.checkCalls[0].Url)
	}

	if len(checker.insertCalls) != 1 {
		t.Fatalf("Expected 1 insert call, got %d", len(checker.insertCalls))
	}

	insertedResult := checker.insertCalls[0]
	if insertedResult.URL != expectedResult.URL {
		t.Errorf("Expected inserted result URL %s, got %s", expectedResult.URL, insertedResult.URL)
	}
}

func TestScheduler_PerformCheck_CheckerInsertError(t *testing.T) {
	url := models.MonitoredUrl{
		ID:               1,
		Url:              "https://example.com",
		CheckIntervalSec: 30,
		RegexPattern:     "",
	}

	status := 200
	responseTime := 100

	checkResult := models.CheckResult{
		URL:            url.Url,
		CheckTimestamp: time.Now(),
		ResponseTimeMs: &responseTime,
		HttpStatus:     &status,
		Error:          "",
	}

	checker := &mockChecker{
		checkResult: checkResult,
		insertError: errors.New("database error"),
	}

	repo := &mockRepository{}
	scheduler := New(repo, nil, checker)

	scheduler.performCheck(url)

	if checker.checkCallCount != 1 {
		t.Errorf("Expected Check to be called once, got %d calls", checker.checkCallCount)
	}

	if checker.insertCallCount != 1 {
		t.Errorf("Expected InsertCheckResult to be called once even with error, got %d calls", checker.insertCallCount)
	}
}

func TestScheduler_PerformCheck_CheckerMultipleInteractions(t *testing.T) {
	urls := []models.MonitoredUrl{
		{
			ID:               1,
			Url:              "https://example.com",
			CheckIntervalSec: 30,
			RegexPattern:     "",
		},
		{
			ID:               2,
			Url:              "https://google.com",
			CheckIntervalSec: 60,
			RegexPattern:     "Google",
		},
	}

	status := 200
	responseTime := 250

	checkResult := models.CheckResult{
		URL:            "test",
		CheckTimestamp: time.Now(),
		ResponseTimeMs: &responseTime,
		HttpStatus:     &status,
		Error:          "",
	}

	checker := &mockChecker{
		checkResult: checkResult,
		insertError: nil,
	}

	repo := &mockRepository{}
	scheduler := New(repo, nil, checker)

	for _, url := range urls {
		scheduler.performCheck(url)
	}

	if checker.checkCallCount != 2 {
		t.Errorf("Expected Check to be called twice, got %d calls", checker.checkCallCount)
	}

	if checker.insertCallCount != 2 {
		t.Errorf("Expected InsertCheckResult to be called twice, got %d calls", checker.insertCallCount)
	}

	if len(checker.checkCalls) != 2 {
		t.Fatalf("Expected 2 check calls, got %d", len(checker.checkCalls))
	}

	for i, url := range urls {
		if checker.checkCalls[i].Url != url.Url {
			t.Errorf("Expected check call %d with URL %s, got %s", i, url.Url, checker.checkCalls[i].Url)
		}
		if checker.checkCalls[i].RegexPattern != url.RegexPattern {
			t.Errorf("Expected check call %d with regex %s, got %s", i, url.RegexPattern, checker.checkCalls[i].RegexPattern)
		}
	}
}
