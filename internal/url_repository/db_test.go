package url_repository

import (
	"testing"
)

// todo GetMonitoredUrls test

func TestDbUrlRepository_New(t *testing.T) {
	repo := New(nil)

	if repo == nil {
		t.Fatal("Expected non-nil repository")
	}

	var _ UrlRepository = repo
}

func TestDbUrlRepository_InterfaceCompliance(t *testing.T) {
	repo := New(nil)

	var _ UrlRepository = repo
}
