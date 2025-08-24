# Local Quick Start

1. **Configure database credentials**

   Create `.env` file with your settings:
   ```env
   # Required
   DB_HOST=db
   DB_PORT=5432
   DB_USER=monitor_user
   DB_PASSWORD=monitor_password
   DB_NAME=monitor_db
   
   # Optional
   DB_HOST_PORT=5432   # Only needed for local Docker (defaults to 5432)
   DB_SSL_MODE=require # Use 'require' for remote host, 'disable' for local
   ```

2. **Start the application**
   ```bash
   # Start PostgreSQL database (only if you want to use local DB)
   docker compose up -d db
   
   # Run database migrations
   docker compose run --rm monitor ./migrate

   # Add urls to the DB
   
   # Start the monitor
   docker compose up monitor
   ```

# Technical Decisions

## Checks

- For the check there is a timeout of 30 seconds.
- The regex is checked against the first 64KB of the page

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DB_HOST` | Yes | PostgreSQL host (e.g., remote host or `db` for local Docker) |
| `DB_PORT` | Yes | PostgreSQL port |
| `DB_USER` | Yes | PostgreSQL username |
| `DB_PASSWORD` | Yes | PostgreSQL password |
| `DB_NAME` | Yes | PostgreSQL database name |
| `DB_HOST_PORT` | No | Host port to expose PostgreSQL (Docker only, defaults to 5432) |
| `DB_SSL_MODE` | No | SSL mode (`disable`, `require`, `prefer`, etc.) - defaults to `require` |

## Graceful Shutdown

The application handles `SIGINT` and `SIGTERM` signals for graceful shutdown:
- Stops all monitoring goroutines
- Waits for in-flight checks to complete

## Migrations

`golang-migrate` is used for migrations. Migration files can be found in `internal/migrations`

## Database Schema

### monitored_urls table
- `id`: Serial primary key
- `url`: Website URL (unique)
- `check_interval_sec`: Check interval in seconds (5-300)
- `regex_pattern`: Optional regex pattern for page validation

### checks table
- `id`: Serial primary key
- `url`: Website URL
- `check_timestamp`: When the check was performed
- `response_time_ms`: HTTP response time in milliseconds
- `http_status`: HTTP status code
- `regex_match`: Regex pattern match indicator (if pattern provided)
- `error`: Error message if check failed

# Testing

**Run all tests:**
```bash
go test ./...
```