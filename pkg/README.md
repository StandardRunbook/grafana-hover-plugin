# Grafana Plugin API

Backend API service for the Grafana Hover Tracker Panel plugin. Analyzes log data using KL divergence to identify anomalous patterns.

## Features

- **KL Divergence Analysis**: Identifies anomalous log patterns by comparing current and baseline time windows
- **ClickHouse Integration**: Efficient log storage and querying
- **Auto-table Creation**: Automatically creates required database tables from schema
- **Multi-platform**: Builds for Linux, macOS, and Windows (amd64 and arm64)
- **REST API**: Simple HTTP API for Grafana integration
- **Comprehensive Tests**: Unit tests for all core functionality

## Architecture

```
cmd/main.go                     - HTTP server entry point
internal/
  ├── api/                      - HTTP handlers and request validation
  ├── analyzer/                 - Log analysis and KL divergence
  ├── clickhouse/              - Database client
  └── config/                   - Configuration loading
schema/                         - ClickHouse schema (git submodule)
```

## Quick Start

### Prerequisites

- Go 1.23+
- ClickHouse (for production)
- Docker (for testing with ClickHouse)

### Development

```bash
# Install dependencies
go mod download

# Run development server
mage dev

# Run tests
mage test

# Run tests with coverage
mage testCoverage
```

### Building

```bash
# Build for current platform
mage build

# Build for all platforms
mage buildAll

# Install to ../grafana-hover-tracker-panel/bin/
mage install

# Clean build artifacts
mage clean
```

### Configuration

Create `config.toml`:

```toml
[server]
host = "127.0.0.1"
port = 8080

[clickhouse]
url = "http://localhost:8123"
user = ""
password = ""
database = "default"
```

### ClickHouse Setup

For testing, use the included docker-compose:

```bash
docker-compose -f docker-compose.test.yml up -d
```

The service will automatically create required tables from the schema on startup.

## API

### POST /query_logs

Analyzes logs for anomalies in a given time window.

**Request:**
```json
{
  "org": "my-org",
  "dashboard": "my-dashboard",
  "panel_title": "my-panel",
  "metric_name": "A-series",
  "start_time": "2025-10-22T04:00:00Z",
  "end_time": "2025-10-22T05:00:00Z"
}
```

**Response:**
```json
{
  "log_groups": [
    {
      "representative_logs": ["ERROR: Out of memory", "ERROR: OOM killer invoked"],
      "relative_change": 1.5
    }
  ]
}
```

## Testing

Tests cover:
- KL divergence calculations
- Relative change calculations
- Log analyzer logic (window calculation, sorting, grouping)
- API request validation
- Error handling and response formatting

```bash
# Run all tests
mage test

# Run with coverage report
mage testCoverage

# Run specific package tests
go test -v ./internal/analyzer
```

## Deployment

### Building Binaries

Binaries are built for all platforms:
- `linux/amd64`
- `linux/arm64`
- `darwin/amd64` (Intel Mac)
- `darwin/arm64` (Apple Silicon)
- `windows/amd64`

```bash
mage buildAll
```

All binaries are statically compiled with `CGO_ENABLED=0` for easy deployment.

### Running in Production

1. Set up ClickHouse
2. Create `config.toml` with production settings
3. Run the binary:

```bash
./grafana-plugin-api
```

The service will:
- Load configuration
- Connect to ClickHouse
- Verify/create required tables
- Start HTTP server

## Development Commands

```bash
mage build          # Build for current platform
mage buildAll       # Build for all platforms
mage install        # Build and install to bin/
mage clean          # Remove build artifacts
mage test           # Run tests
mage testCoverage   # Run tests with coverage
mage tidy           # Run go mod tidy
mage dev            # Run development server
```
