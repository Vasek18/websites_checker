# Website Monitor

A Go application that monitors multiple websites periodically, collecting metrics (response time, HTTP status, regex matches) and storing them in a PostgreSQL database.

## Features

- Monitors multiple websites concurrently using goroutines
- Configurable check intervals (5-300 seconds) per URL
- Optional regex pattern matching on response content
- PostgreSQL database storage for check results
- Graceful shutdown handling
- No external dependencies except PostgreSQL driver

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database

## Installation

1. Clone or download the project
2. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

### Database Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your PostgreSQL connection details:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=monitor_user
   DB_PASSWORD=your_password_here
   DB_NAME=monitor_db
   ```

### URLs Configuration

Create a `urls.yaml` file in the project root with the websites you want to monitor:

```yaml
urls:
  - url: "https://example.com"
    interval: 60              # Check every 60 seconds
    regex: "Example Domain"   # Optional: regex pattern to match in response
  - url: "https://github.com"
    interval: 120             # Check every 2 minutes
    regex: "GitHub"
  - url: "https://httpbin.org/status/200"
    interval: 30              # Check every 30 seconds
```

Alternatively, you can use JSON format (`urls.json`):

```json
{
  "urls": [
    {
      "url": "https://example.com",
      "interval": 60,
      "regex": "Example Domain"
    },
    {
      "url": "https://github.com",
      "interval": 120,
      "regex": "GitHub"
    }
  ]
}
```

## Database Setup

1. Create a PostgreSQL database and user
2. Run database migrations:
   ```bash
   go run ./cmd/migrate
   ```

This will create the required tables:
- `monitored_urls`: Stores URL configurations
- `checks`: Stores check results with timestamps, response times, status codes, and errors

## Running the Application

Start the website monitor:
```bash
go run ./cmd/monitor
```

The application will:
1. Load configuration from `.env` and `urls.yaml`
2. Connect to the PostgreSQL database
3. Start monitoring all configured URLs in separate goroutines
4. Log check results and store them in the database
5. Continue running until interrupted (Ctrl+C)

## Building

To build the application:
```bash
# Build monitor
go build -o monitor ./cmd/monitor

# Build migration tool
go build -o migrate ./cmd/migrate

# Run built binaries
./migrate
./monitor
```

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

## Configuration Validation

- Check intervals must be between 5 and 300 seconds
- URLs must be valid and non-empty
- Database connection parameters are required
- Regex patterns are validated before use

## Graceful Shutdown

The application handles `SIGINT` and `SIGTERM` signals for graceful shutdown:
- Stops all monitoring goroutines
- Waits for in-flight checks to complete
- Closes database connections

## Logging

The application logs:
- Startup and shutdown events
- Check results (success/failure)
- Database operations
- Configuration loading

## Architecture

- **Repository Pattern**: Abstracts URL data sources (currently file-based, easily extensible to database)
- **Goroutine per URL**: Each monitored URL runs in its own goroutine with independent timing
- **Raw SQL**: Uses `database/sql` with PostgreSQL driver, no ORM
- **Standard Library**: Minimal external dependencies

## Development

Project structure:
```
├── cmd/
│   ├── monitor/main.go     # Main application
│   └── migrate/main.go     # Database migration tool
├── internal/
│   ├── config/             # Configuration loading
│   ├── db/                 # Database operations
│   ├── models/             # Data structures
│   ├── repository/         # URL data sources
│   ├── checker/            # HTTP checking logic
│   └── scheduler/          # Goroutine-based scheduling
├── urls.yaml               # URL configuration
├── .env                    # Database configuration
└── README.md
```