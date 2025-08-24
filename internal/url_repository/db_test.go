package url_repository_test

import (
	"database/sql"
	"testing"

	"website-monitor/internal/db"
	"website-monitor/internal/models"
	"website-monitor/internal/url_repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetMonitoredUrls_HappyPath(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock db: %v", err)
	}
	defer sqlDB.Close()

	mock.ExpectQuery(`SELECT id, url, check_interval_sec, COALESCE\(regex_pattern, ''\) FROM monitored_urls`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "url", "check_interval_sec", "regex_pattern"}).
				AddRow(1, "https://stackoverflow.com", 10, "Example").
				AddRow(2, "https://google.com", 120, "Google").
				AddRow(3, "https://github.com", 30, ""),
		)

	repo := url_repository.New(db.New(sqlDB))

	urls, err := repo.GetMonitoredUrls()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(urls) != 3 {
		t.Fatalf("expected 3 URLs, got %d", len(urls))
	}

	expected := []models.MonitoredUrl{
		{ID: 1, Url: "https://stackoverflow.com", CheckIntervalSec: 10, RegexPattern: "Example"},
		{ID: 2, Url: "https://google.com", CheckIntervalSec: 120, RegexPattern: "Google"},
		{ID: 3, Url: "https://github.com", CheckIntervalSec: 30, RegexPattern: ""},
	}

	for i, exp := range expected {
		if urls[i] != exp {
			t.Errorf("expected %+v, got %+v", exp, urls[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestGetMonitoredUrls_EmptyTable(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	mock.ExpectQuery(`SELECT id, url, check_interval_sec, COALESCE\(regex_pattern, ''\) FROM monitored_urls`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url", "check_interval_sec", "regex_pattern"}))

	repo := url_repository.New(db.New(sqlDB))

	urls, err := repo.GetMonitoredUrls()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(urls) != 0 {
		t.Errorf("expected empty result, got %d rows", len(urls))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestGetMonitoredUrls_QueryError(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	defer sqlDB.Close()

	mock.ExpectQuery(`SELECT id, url, check_interval_sec, COALESCE\(regex_pattern, ''\) FROM monitored_urls`).
		WillReturnError(sql.ErrConnDone)

	repo := url_repository.New(db.New(sqlDB))

	urls, err := repo.GetMonitoredUrls()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if urls != nil {
		t.Errorf("expected nil result, got %+v", urls)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}
