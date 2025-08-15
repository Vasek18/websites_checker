package checker

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"website-monitor/internal/db"
	"website-monitor/internal/models"
)

type HTTPChecker struct {
	client *http.Client
	db     *db.DB
}

// New creates a new HTTP checker with a configured client
func New(database *db.DB) *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		db: database,
	}
}

// Check performs an HTTP check on the given url and returns the result
func (c *HTTPChecker) Check(url models.MonitoredUrl) models.CheckResult {
	result := models.CheckResult{
		URL:            url.Url,
		CheckTimestamp: time.Now(),
	}

	start := time.Now()
	resp, err := c.client.Get(url.Url)
	responseTime := int(time.Since(start).Milliseconds())
	result.ResponseTimeMs = &responseTime

	if err != nil {
		result.Error = err.Error()

		return result
	}
	defer func(Body io.ReadCloser) { // todo should be before the error return?
		_ = Body.Close()
	}(resp.Body)

	result.HTTPStatus = &resp.StatusCode

	// Check a regexp pattern if provided
	if url.RegexPattern != "" {
		regexMatch, err := c.checkRegexPattern(resp, url.RegexPattern)
		if err != nil {
			result.Error = fmt.Sprintf("regex check failed: %s", err.Error())
		} else {
			result.RegexMatch = &regexMatch // todo why pointer?
		}
	}

	return result
}

// checkRegexPattern checks if the response body matches the given regex pattern
func (c *HTTPChecker) checkRegexPattern(resp *http.Response, pattern string) (bool, error) {
	regex, err := regexp.Compile(pattern) // todo safe?
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Read response body (limit to 64KB for performance)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	return regex.Match(body), nil
}

// InsertCheckResult inserts a check result into the database
func (c *HTTPChecker) InsertCheckResult(result models.CheckResult) error {
	query := `
		INSERT INTO checks (url, check_timestamp, response_time_ms, http_status, regex_match, error)
		VALUES ($1, $2, $3, $4, $5, $6)`

	err := c.db.Exec(query,
		result.URL,
		result.CheckTimestamp,
		result.ResponseTimeMs,
		result.HTTPStatus,
		result.RegexMatch,
		result.Error)

	if err != nil {
		return fmt.Errorf("failed to insert check result: %w", err)
	}

	return nil
}
