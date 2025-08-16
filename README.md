# Local Quick Start

1. **Configure database credentials**

   Create `.env` file with your settings:
   ```env
   # Required
   DB_HOST=your-aiven-host.com  # or 'db' for local Docker
   DB_PORT=5432
   DB_USER=monitor_user
   DB_PASSWORD=secret
   DB_NAME=monitor_db
   
   # Optional
   DB_HOST_PORT=5432          # Only needed for Docker (defaults to 5432)
   DB_SSL_MODE=require        # Use 'require' for Aiven, 'disable' for local
   ```

2. **Start the application**
   ```bash
   # Start PostgreSQL database (if you want to use local DB)
   docker compose up db -d
   
   # Run database migrations
   docker compose run --rm monitor ./migrate
   
   # Seed database with sample URLs (optional)
   docker compose run --rm monitor ./seed
   
   # Start the monitor
   docker compose up monitor
   ```

# Production Setup

# Technical Decisions

## URLs List

I would store URLs in the database (as implemented) because to iterate over data periodically, it should be stored persistently. We could also use some kind of file storage, e.g., JSON or YAML config, but updating big formatted files is error-prone. Also, with a database, we can introduce pagination when the list becomes too big, to avoid loading it all into memory.
However, I introduced the repository pattern here, so the storage mechanism can be easily replaced in the future.

## Checks

- For the check there is a timeout of 30 seconds.
- The regex is checked against the first 64KB of the page

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DB_HOST` | Yes | PostgreSQL host (e.g., Aiven host or `db` for Docker) |
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

We used `golang-migrate` for migrations. That is the same thing I would use on prod, however with a couple of changes:
- Migrations would be run in their own container. Something similar is actually implemented in docker-compose
- There would be down migrations

Migration files can be found in `internal/migrations`

## Local

Some features here are just for demo purposes, e.g. seeder (see the comment in cmd/seed/main.go)

## Docker

Since containerization won't affect the evaluation, I prepared a simple Dockerfile and docker-compose.yml only for demo purposes. They were almost entirely copied from other projects.

## Database Schema

### monitored_urls table
- `id`: Serial primary key
- `url`: Website URL (unique)
- `check_interval_sec`: Check interval in seconds (5-300)
- `regex_pattern`: Optional regex pattern for response validation

### checks table
- `id`: Serial primary key
- `url`: Website URL
- `check_timestamp`: When the check was performed
- `response_time_ms`: HTTP response time in milliseconds
- `http_status`: HTTP status code
- `regex_match`: Boolean indicating regex pattern match (if pattern provided)
- `error`: Error message if check failed

# Testing

**Run all tests:**
```bash
go test ./...
```