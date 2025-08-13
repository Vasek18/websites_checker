# Project Specification: Website Availability Monitor

## Overview
We are building a Go application that monitors the availability of multiple websites, collects metrics, and stores them into a PostgreSQL database.

## Requirements

### General
- Language: **Go**
- Main purpose: Periodically check multiple URLs for availability, response time, HTTP status, and optional regex match in response body.
- Scale target: Should handle **thousands of sites** efficiently.
- Code quality: Production-ready, tested, maintainable.
- Tests: Unit tests and basic integration tests.
- Concurrency: **No external scheduling libraries** — implement scheduling with `time.Ticker` and goroutines.
- DB access: **No ORM** — use `database/sql` with `lib/pq` and raw SQL.
- Structure: Use repository pattern for data sources.

---

## Architecture

### Directory Structure
cmd/monitor/main.go # Application entry point
cmd/migrate/main.go # Runs DB migrations
internal/config/config.go # Loads and validates config from .env and yaml/json
internal/db/db.go # PostgreSQL connection
internal/db/migrations.go # Embedded SQL migrations
internal/repository/repository.go # Interface for URL source
internal/repository/file.go # Reads monitored URLs from config file
internal/models/models.go # Data structures
internal/checker/checker.go # HTTP check logic
internal/scheduler/scheduler.go # Goroutine-based scheduling

---

## Configuration

### Sources
- **DB credentials** and other sensitive data: from `.env` file.
- **Monitored URLs**: from `urls.yaml` or `urls.json` in project root.

### Example `.env`
DB_HOST=localhost
DB_PORT=5432
DB_USER=monitor_user
DB_PASSWORD=secret
DB_NAME=monitor_db

### Example `urls.yaml`
```yaml
urls:
  - url: "https://example.com"
    interval: 60
    regex: "Example Domain"
  - url: "https://aiven.io"
    interval: 120
```

## Database
### Engine
PostgreSQL (via Docker Compose)

## Schema
```sql
CREATE TABLE monitored_urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    check_interval_sec INT NOT NULL CHECK (check_interval_sec BETWEEN 5 AND 300),
    regex_pattern TEXT
);

CREATE TABLE checks (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    check_timestamp TIMESTAMPTZ NOT NULL,
    response_time_ms INT,
    http_status INT,
    regex_match BOOLEAN,
    error TEXT
);
```

## Migrations
Migrations should be run with a separate command: "go run ./cmd/migrate"
SQL migrations are embedded in the code.

## Scheduling
- One goroutine per monitored URL.
- Uses time.Ticker for interval-based execution.
- Graceful shutdown via context.Context.

## Repository Pattern
- Initial implementation reads monitored URLs from urls.yaml / urls.json.
- Later, can be swapped to DB-backed repository without changing scheduler or checker logic.

## Testing
- Unit tests for:
  - Checker (use httptest.Server to simulate responses)
  - Repository (load from file)
- Integration tests for:
    - DB insertions (Postgres via Docker Compose or testcontainers)
- Scheduler tests using mocked tickers to avoid real-time waiting.

## Logging
- Use Go standard library log for simplicity.
- No separate logging helper unless needed.
