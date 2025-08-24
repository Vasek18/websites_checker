package checker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"website-monitor/internal/models"
)

func TestChecker_Check_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response body"))
	}))
	defer server.Close()

	checker := &Checker{
		client: &http.Client{},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
	}

	result := checker.Check(url)

	if result.URL != url.Url {
		t.Errorf("Expected URL %s, got %s", url.Url, result.URL)
	}

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}

	if result.HttpStatus == nil || *result.HttpStatus != 200 {
		t.Errorf("Expected HTTP status 200, got %v", result.HttpStatus)
	}

	if result.ResponseTimeMs == nil || *result.ResponseTimeMs < 0 {
		t.Error("Expected response time to be set and positive")
	}

	if !result.CheckTimestamp.After(time.Now().Add(-time.Minute)) {
		t.Error("Expected check timestamp to be recent")
	}
}

func TestChecker_Check_WithRegexMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hi, my name is Test Server!"))
	}))
	defer server.Close()

	checker := &Checker{
		client: &http.Client{},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
		RegexPattern:     "Hi, my name is",
	}

	result := checker.Check(url)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}

	if result.RegexMatch == nil || !*result.RegexMatch {
		t.Error("Expected regex to match")
	}
}

func TestChecker_Check_WithRegexNoMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("In another castle"))
	}))
	defer server.Close()

	checker := &Checker{
		client: &http.Client{},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
		RegexPattern:     "Princess",
	}

	result := checker.Check(url)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}

	if result.RegexMatch == nil || *result.RegexMatch {
		t.Error("Expected regex not to match")
	}
}

func TestChecker_Check_InvalidRegex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Nothing to see here"))
	}))
	defer server.Close()

	checker := &Checker{
		client: &http.Client{},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
		RegexPattern:     "[invalid",
	}

	result := checker.Check(url)

	if result.Error == "" {
		t.Error("Expected error")
	}

	if !strings.Contains(result.Error, "regex check failed") {
		t.Errorf("Expected regex error, got: %s", result.Error)
	}
}

func TestChecker_Check_HttpError(t *testing.T) {
	checker := &Checker{
		client: &http.Client{},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              "http://invalid.nonexistent.domain.example",
		CheckIntervalSec: 60,
	}

	result := checker.Check(url)

	if result.Error == "" {
		t.Error("Expected error for invalid URL")
	}

	if result.HttpStatus != nil {
		t.Error("Expected no HTTP status for failed request")
	}

	if result.ResponseTimeMs == nil {
		t.Error("Expected response time to be set even for failed requests")
	}
}

func TestChecker_Check_404Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer server.Close()

	checker := &Checker{
		client: &http.Client{},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
	}

	result := checker.Check(url)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}

	if result.HttpStatus == nil || *result.HttpStatus != 404 {
		t.Errorf("Expected HTTP status 404, got %v", result.HttpStatus)
	}
}
