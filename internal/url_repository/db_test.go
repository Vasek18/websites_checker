package url_repository

import (
	"errors"
	"testing"

	"website-monitor/internal/models"
)

// todo mb remove

// RowsInterface defines the interface for sql.Rows operations
type RowsInterface interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// DBInterface defines the interface that our DB type implements
type DBInterface interface {
	Query(query string, args ...interface{}) (RowsInterface, error)
	Exec(query string, args ...interface{}) error
	Close() error
}

// MockDB implements DBInterface for testing
type MockDB struct {
	queryFunc func(query string, args ...interface{}) (RowsInterface, error)
	execFunc  func(query string, args ...interface{}) error
	closeFunc func() error
}

func (m *MockDB) Query(query string, args ...interface{}) (RowsInterface, error) {
	if m.queryFunc != nil {
		return m.queryFunc(query, args...)
	}
	return nil, errors.New("mock query not implemented")
}

func (m *MockDB) Exec(query string, args ...interface{}) error {
	if m.execFunc != nil {
		return m.execFunc(query, args...)
	}
	return errors.New("mock exec not implemented")
}

func (m *MockDB) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

// MockRows simulates sql.Rows for testing
type MockRows struct {
	data    [][]interface{}
	current int
	closed  bool
	err     error
}

func NewMockRows(data [][]interface{}) *MockRows {
	return &MockRows{
		data:    data,
		current: -1,
	}
}

func (m *MockRows) Next() bool {
	if m.closed || m.current >= len(m.data)-1 {
		return false
	}
	m.current++
	return true
}

func (m *MockRows) Scan(dest ...interface{}) error {
	if m.closed {
		return errors.New("rows are closed")
	}
	if m.current < 0 || m.current >= len(m.data) {
		return errors.New("no current row")
	}

	row := m.data[m.current]
	if len(dest) != len(row) {
		return errors.New("destination count does not match column count")
	}

	for i, val := range row {
		switch d := dest[i].(type) {
		case *int:
			if v, ok := val.(int); ok {
				*d = v
			} else {
				return errors.New("type mismatch for int")
			}
		case *string:
			if v, ok := val.(string); ok {
				*d = v
			} else {
				return errors.New("type mismatch for string")
			}
		default:
			return errors.New("unsupported destination type")
		}
	}
	return nil
}

func (m *MockRows) Close() error {
	m.closed = true
	return nil
}

func (m *MockRows) Err() error {
	return m.err
}

// MockDbUrlRepository wraps DbUrlRepository to use our interface
type MockDbUrlRepository struct {
	db DBInterface
}

func NewMockRepository(db DBInterface) *MockDbUrlRepository {
	return &MockDbUrlRepository{
		db: db,
	}
}

// GetMonitoredUrls returns all URLs that should be monitored from the database
func (r *MockDbUrlRepository) GetMonitoredUrls() ([]models.MonitoredUrl, error) {
	query := `SELECT id, url, check_interval_sec, COALESCE(regex_pattern, '') FROM monitored_urls`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	urls := make([]models.MonitoredUrl, 0)
	for rows.Next() {
		var url models.MonitoredUrl
		if err := rows.Scan(&url.ID, &url.Url, &url.CheckIntervalSec, &url.RegexPattern); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func TestDbUrlRepository_GetMonitoredUrls_HappyPath(t *testing.T) {
	// Setup mock data
	mockData := [][]interface{}{
		{1, "https://example.com", 60, "Example"},
		{2, "https://google.com", 120, "Google"},
		{3, "https://test.com", 30, ""},
	}

	mockRows := NewMockRows(mockData)
	mockDB := &MockDB{
		queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
			// Verify correct query is called
			expectedQuery := `SELECT id, url, check_interval_sec, COALESCE(regex_pattern, '') FROM monitored_urls`
			if query != expectedQuery {
				t.Errorf("Expected query: %s, got: %s", expectedQuery, query)
			}
			return mockRows, nil
		},
	}

	repo := NewMockRepository(mockDB)
	urls, err := repo.GetMonitoredUrls()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(urls) != 3 {
		t.Fatalf("Expected 3 URLs, got %d", len(urls))
	}

	// Verify first URL
	expectedURL1 := models.MonitoredUrl{
		ID:               1,
		Url:              "https://example.com",
		CheckIntervalSec: 60,
		RegexPattern:     "Example",
	}
	if urls[0] != expectedURL1 {
		t.Errorf("Expected first URL %+v, got %+v", expectedURL1, urls[0])
	}

	// Verify second URL
	expectedURL2 := models.MonitoredUrl{
		ID:               2,
		Url:              "https://google.com",
		CheckIntervalSec: 120,
		RegexPattern:     "Google",
	}
	if urls[1] != expectedURL2 {
		t.Errorf("Expected second URL %+v, got %+v", expectedURL2, urls[1])
	}

	// Verify third URL (with empty regex)
	expectedURL3 := models.MonitoredUrl{
		ID:               3,
		Url:              "https://test.com",
		CheckIntervalSec: 30,
		RegexPattern:     "",
	}
	if urls[2] != expectedURL3 {
		t.Errorf("Expected third URL %+v, got %+v", expectedURL3, urls[2])
	}
}

func TestDbUrlRepository_GetMonitoredUrls_EmptyTable(t *testing.T) {
	// Setup mock with no data
	mockRows := NewMockRows([][]interface{}{})
	mockDB := &MockDB{
		queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
			return mockRows, nil
		},
	}

	repo := NewMockRepository(mockDB)
	urls, err := repo.GetMonitoredUrls()

	if err != nil {
		t.Fatalf("Expected no error with empty table, got: %v", err)
	}

	if urls == nil {
		t.Error("Expected non-nil slice, got nil")
	}

	if len(urls) != 0 {
		t.Errorf("Expected empty slice, got %d URLs", len(urls))
	}
}

func TestDbUrlRepository_GetMonitoredUrls_QueryError(t *testing.T) {
	expectedError := errors.New("database connection failed")
	mockDB := &MockDB{
		queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
			return nil, expectedError
		},
	}

	repo := NewMockRepository(mockDB)
	urls, err := repo.GetMonitoredUrls()

	if err == nil {
		t.Fatal("Expected error from database query")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if urls != nil {
		t.Error("Expected nil URLs on error")
	}
}

func TestDbUrlRepository_GetMonitoredUrls_ScanError(t *testing.T) {
	// Setup mock data with invalid types to cause scan error
	mockData := [][]interface{}{
		{"invalid_id", "https://example.com", 60, "Example"}, // string instead of int for ID
	}

	mockRows := NewMockRows(mockData)
	mockDB := &MockDB{
		queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
			return mockRows, nil
		},
	}

	repo := NewMockRepository(mockDB)
	urls, err := repo.GetMonitoredUrls()

	if err == nil {
		t.Fatal("Expected error from scan operation")
	}

	if urls != nil {
		t.Error("Expected nil URLs on error")
	}

	expectedErrorMessage := "type mismatch for int"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestDbUrlRepository_GetMonitoredUrls_RowsError(t *testing.T) {
	expectedRowsError := errors.New("rows iteration error")
	mockRows := NewMockRows([][]interface{}{
		{1, "https://example.com", 60, "Example"},
	})
	mockRows.err = expectedRowsError

	mockDB := &MockDB{
		queryFunc: func(query string, args ...interface{}) (RowsInterface, error) {
			return mockRows, nil
		},
	}

	repo := NewMockRepository(mockDB)
	urls, err := repo.GetMonitoredUrls()

	if err == nil {
		t.Fatal("Expected error from rows.Err()")
	}

	if err != expectedRowsError {
		t.Errorf("Expected error %v, got %v", expectedRowsError, err)
	}

	if urls != nil {
		t.Error("Expected nil URLs on error")
	}
}
