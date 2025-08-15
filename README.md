# Website Monitor

A Go application that monitors multiple websites periodically, collecting metrics (response time, HTTP status, regex matches) and storing them in a PostgreSQL database.

## Features

- Monitors multiple websites concurrently using goroutines
- Configurable check intervals (5-300 seconds) per URL
- Optional regex pattern matching on response content
- PostgreSQL database storage for check results
- Graceful shutdown handling
- Docker Compose setup for easy local development
- Versioned database migrations using golang-migrate

## Prerequisites

- Docker and Docker Compose (recommended)
- **OR** Go 1.21+ and PostgreSQL (for manual setup)

## Quick Start with Docker Compose

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd websites_checker
   ```

2. **Configure database credentials**
   
   Edit `.env` file with your database settings:
   ```env
   DB_HOST=db
   DB_PORT=5432
   DB_HOST_PORT=5432
   DB_USER=monitor_user
   DB_PASSWORD=secret
   DB_NAME=monitor_db
   ```

3. **Start the application**
   ```bash
   # Start PostgreSQL database
   docker compose up db -d
   
   # Run database migrations
   docker compose run --rm monitor ./migrate
   
   # Seed database with sample URLs
   docker compose run --rm monitor ./seed
   
   # Start the monitor
   docker compose up monitor
   ```

   Or run everything at once:
   ```bash
   docker compose up -d db
   docker compose run --rm monitor ./migrate
   docker compose run --rm monitor ./seed
   docker compose up monitor
   ```

4. **Stop the application**
   ```bash
   docker compose down
   ```

## Manual Setup (without Docker)

### Prerequisites
- Go 1.21 or higher
- PostgreSQL database

### Installation

1. **Clone and install dependencies**
   ```bash
   git clone <repository-url>
   cd websites_checker
   go mod download
   ```

2. **Set up PostgreSQL database**
   ```bash
   # Create database and user
   createdb monitor_db
   createuser monitor_user
   # Grant permissions as needed
   ```

3. **Set environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=monitor_user
   export DB_PASSWORD=your_password_here
   export DB_NAME=monitor_db
   ```

4. **Run database migrations**
   ```bash
   go run ./cmd/migrate
   ```

5. **Seed database with sample URLs**
   ```bash
   go run ./cmd/seed
   ```

6. **Start the monitor**
   ```bash
   go run ./cmd/monitor
   ```

   The application will:
   - Load configuration from environment variables
   - Connect to the PostgreSQL database
   - Load monitored URLs from the database
   - Start monitoring all configured URLs in separate goroutines
   - Log check results and store them in the database
   - Continue running until interrupted (Ctrl+C)

## Building

**For Docker deployment:**
```bash
# Build the Docker image
docker-compose build

# Or build manually
docker build -t website-monitor .
```

**For manual deployment:**
```bash
# Build all binaries
go build -o monitor ./cmd/monitor
go build -o migrate ./cmd/migrate
go build -o seed ./cmd/seed

# Run built binaries
./migrate
./seed
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

## Testing

The project includes comprehensive unit tests for all internal packages.

### Running Tests

**Run all tests:**
```bash
go test ./internal/...
```

**Run tests with verbose output:**
```bash
go test -v ./internal/...
```

**Run tests with coverage:**
```bash
go test -cover ./internal/...
```

**Run tests for specific package:**
```bash
go test ./internal/config
go test ./internal/checker
go test ./internal/models
```

**Run specific test function:**
```bash
go test -run TestHTTPChecker_Check_Success ./internal/checker
```

**Generate detailed coverage report:**
```bash
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

The last command generates an HTML coverage report you can open in your browser.

### Test Coverage

- **config**: 100% coverage - Environment variable loading and validation
- **checker**: 79.3% coverage - HTTP checking, regex matching, error handling
- **models**: No statements to cover - Pure data structures with JSON marshaling tests
- **db**: 25.9% coverage - Database connection and URL generation
- **scheduler**: 17.8% coverage - Goroutine management and graceful shutdown
- **url_repository**: 6.7% coverage - Repository pattern implementation

### Test Features

- **Mock implementations** for testing without real database connections
- **HTTP test servers** for safe HTTP functionality testing
- **Environment variable management** with proper cleanup
- **Error scenario testing** including network failures and invalid inputs
- **Interface compliance verification**
- **JSON serialization testing** for data models
- **Concurrent execution safety** testing

## Development

Project structure:
```
├── cmd/
│   ├── monitor/main.go         # Main application
│   ├── migrate/main.go         # Database migration tool
│   └── seed/main.go            # Database seeding tool
├── internal/
│   ├── config/                 # Configuration loading
│   ├── db/                     # Database operations
│   ├── migrations/             # SQL migration files
│   │   ├── 000001_create_monitored_urls.up.sql
│   │   ├── 000001_create_monitored_urls.down.sql
│   │   ├── 000002_create_checks.up.sql
│   │   └── 000002_create_checks.down.sql
│   ├── models/                 # Data structures
│   ├── repository/             # URL data sources (file & database)
│   │   ├── repository.go       # Repository interface
│   │   ├── file.go            # File-based repository (legacy)
│   │   └── db.go              # Database repository
│   ├── checker/                # HTTP checking logic
│   └── scheduler/              # Goroutine-based scheduling
├── .env                        # Environment variables for Docker
├── docker-compose.yml          # Docker Compose setup
├── Dockerfile                  # Docker image definition
└── README.md
```

## Environment Variables

The application uses the following environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | PostgreSQL host | - | Yes |
| `DB_PORT` | PostgreSQL port | - | Yes |
| `DB_HOST_PORT` | Host port to expose PostgreSQL | - | Yes (Docker only) |
| `DB_USER` | PostgreSQL username | - | Yes |
| `DB_PASSWORD` | PostgreSQL password | - | Yes |
| `DB_NAME` | PostgreSQL database name | - | Yes |

In Docker Compose, these are automatically set in the `docker-compose.yml` file.