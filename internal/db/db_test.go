package db

import (
	"os"
	"testing"
)

func TestGetUrl_Success(t *testing.T) {
	setTestEnvVars()
	defer clearTestEnvVars()

	url, err := GetUrl()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedURL := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=require"
	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}
}

func TestGetUrl_MissingConfig(t *testing.T) {
	clearTestEnvVars()

	_, err := GetUrl()
	if err == nil {
		t.Fatal("Expected error for missing configuration")
	}

	if err.Error() == "" {
		t.Error("Expected error message to be non-empty")
	}
}

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
