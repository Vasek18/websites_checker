package checker

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"website-monitor/internal/models"
)

// HTTPChecker performs HTTP checks on URLs
type HTTPChecker struct {
	client *http.Client
}

// NewHTTPChecker creates a new HTTP checker with a configured client
func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Check performs an HTTP check on the given URL and returns the result
func (c *HTTPChecker) Check(url models.MonitoredURL) models.CheckResult {
	result := models.CheckResult{
		URL:            url.URL,
		CheckTimestamp: time.Now(),
	}

	start := time.Now()
	resp, err := c.client.Get(url.URL)
	responseTime := int(time.Since(start).Milliseconds())
	result.ResponseTimeMs = &responseTime

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.HTTPStatus = &resp.StatusCode

	// Check regex pattern if provided
	if url.RegexPattern != "" {
		regexMatch, err := c.checkRegexPattern(resp, url.RegexPattern)
		if err != nil {
			result.Error = fmt.Sprintf("regex check failed: %s", err.Error())
		} else {
			result.RegexMatch = &regexMatch
		}
	}

	return result
}

// checkRegexPattern checks if the response body matches the given regex pattern
func (c *HTTPChecker) checkRegexPattern(resp *http.Response, pattern string) (bool, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Read response body
	buf := make([]byte, 4096) // Read first 4KB for regex matching
	n, err := resp.Body.Read(buf)
	if err != nil && n == 0 {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	return regex.Match(buf[:n]), nil
}
