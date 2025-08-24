package config

import (
	"os"
	"testing"
	"website-monitor/internal/models"
)

func TestLoad_Success(t *testing.T) {
	setTestEnvVars()
	defer clearTestEnvVars()

	config, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be non-nil")
	}

	expected := &models.Config{
		Database: models.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
			SSLMode:  "require",
		},
	}

	if config.Database != expected.Database {
		t.Errorf("Expected database config %+v, got %+v", expected.Database, config.Database)
	}
}

func TestLoad_MissingEnvVars(t *testing.T) {
	clearTestEnvVars()

	config, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing environment variables")
	}

	if config != nil {
		t.Fatal("Expected config to be nil when error occurs")
	}
}

func TestLoadDatabaseConfig_MissingHost(t *testing.T) {
	setTestEnvVars()
	os.Unsetenv("DB_HOST")
	defer clearTestEnvVars()

	_, err := loadDatabaseConfig()
	if err == nil {
		t.Fatal("Expected error for missing DB_HOST")
	}

	expectedError := "required environment variable DB_HOST not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadDatabaseConfig_MissingPort(t *testing.T) {
	setTestEnvVars()
	os.Unsetenv("DB_PORT")
	defer clearTestEnvVars()

	_, err := loadDatabaseConfig()
	if err == nil {
		t.Fatal("Expected error for missing DB_PORT")
	}

	expectedError := "required environment variable DB_PORT not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadDatabaseConfig_MissingUser(t *testing.T) {
	setTestEnvVars()
	os.Unsetenv("DB_USER")
	defer clearTestEnvVars()

	_, err := loadDatabaseConfig()
	if err == nil {
		t.Fatal("Expected error for missing DB_USER")
	}

	expectedError := "required environment variable DB_USER not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadDatabaseConfig_MissingPassword(t *testing.T) {
	setTestEnvVars()
	os.Unsetenv("DB_PASSWORD")
	defer clearTestEnvVars()

	_, err := loadDatabaseConfig()
	if err == nil {
		t.Fatal("Expected error for missing DB_PASSWORD")
	}

	expectedError := "required environment variable DB_PASSWORD not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadDatabaseConfig_MissingName(t *testing.T) {
	setTestEnvVars()
	os.Unsetenv("DB_NAME")
	defer clearTestEnvVars()

	_, err := loadDatabaseConfig()
	if err == nil {
		t.Fatal("Expected error for missing DB_NAME")
	}

	expectedError := "required environment variable DB_NAME not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
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
