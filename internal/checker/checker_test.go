package checker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"website-monitor/internal/models"
)

func TestHTTPChecker_Check_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response body"))
	}))
	defer server.Close()

	checker := &HTTPChecker{
		client: &http.Client{Timeout: 5 * time.Second},
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

	if result.HTTPStatus == nil || *result.HTTPStatus != 200 { // todo why pointer?
		t.Errorf("Expected HTTP status 200, got %v", result.HTTPStatus)
	}

	if result.ResponseTimeMs == nil || *result.ResponseTimeMs < 0 { // todo why pointer?
		t.Error("Expected response time to be set and positive")
	}

	if !result.CheckTimestamp.After(time.Now().Add(-time.Minute)) {
		t.Error("Expected check timestamp to be recent")
	}
}

func TestHTTPChecker_Check_WithRegexMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Example Domain - This is a test page"))
	}))
	defer server.Close()

	checker := &HTTPChecker{
		client: &http.Client{Timeout: 5 * time.Second},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
		RegexPattern:     "Example Domain",
	}

	result := checker.Check(url)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}

	if result.RegexMatch == nil || !*result.RegexMatch { // todo why pointer?
		t.Error("Expected regex to match")
	}
}

func TestHTTPChecker_Check_WithRegexNoMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Different content here"))
	}))
	defer server.Close()

	checker := &HTTPChecker{
		client: &http.Client{Timeout: 5 * time.Second},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
		RegexPattern:     "Example Domain",
	}

	result := checker.Check(url)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}

	if result.RegexMatch == nil || *result.RegexMatch { // todo why pointer?
		t.Error("Expected regex not to match")
	}
}

func TestHTTPChecker_Check_InvalidRegex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test content"))
	}))
	defer server.Close()

	checker := &HTTPChecker{
		client: &http.Client{Timeout: 5 * time.Second},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
		RegexPattern:     "[invalid",
	}

	result := checker.Check(url)

	if result.Error == "" {
		t.Error("Expected error for invalid regex pattern")
	}

	if !strings.Contains(result.Error, "regex check failed") {
		t.Errorf("Expected regex error, got: %s", result.Error)
	}
}

func TestHTTPChecker_Check_HTTPError(t *testing.T) {
	checker := &HTTPChecker{
		client: &http.Client{Timeout: 1 * time.Millisecond},
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

	if result.HTTPStatus != nil {
		t.Error("Expected no HTTP status for failed request")
	}

	if result.ResponseTimeMs == nil {
		t.Error("Expected response time to be set even for failed requests")
	}
}

func TestHTTPChecker_Check_404Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer server.Close()

	checker := &HTTPChecker{
		client: &http.Client{Timeout: 5 * time.Second},
	}

	url := models.MonitoredUrl{
		ID:               1,
		Url:              server.URL,
		CheckIntervalSec: 60,
	}

	result := checker.Check(url)

	if result.Error != "" {
		t.Errorf("Expected no error for 404 status, got: %s", result.Error)
	}

	if result.HTTPStatus == nil || *result.HTTPStatus != 404 {
		t.Errorf("Expected HTTP status 404, got %v", result.HTTPStatus)
	}
}
