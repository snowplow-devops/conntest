# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Conntest is a command-line utility for validating connections to Snowplow-supported destinations including Snowflake, BigQuery, PostgreSQL, and Databricks. The tool outputs structured JSON events containing connection test results.

**Version Requirements:**
- Go 1.24+ (updated from 1.18)
- Uses latest dependency versions as of 2025

## Development Environment

This project uses standard Go tooling with Make for build automation:

**Prerequisites:**
- Go 1.24+
- Docker (for integration tests)
- Make

## Common Commands

### Building
```shell
# Build for local development
make build

# Build all cross-platform binaries
make all

# Build specific platform
make cli-linux-amd64
make cli-linux-arm64
make cli-darwin-amd64
make cli-darwin-arm64
```

### Testing
```shell
# Run unit tests with coverage
make test

# Run integration tests
make integration-test
```

### Code Quality
```shell
# Format code
make format

# Lint code
make lint

# Tidy modules
make tidy

# Update dependencies
make update
```

### Cleanup
```shell
# Remove build artifacts
make clean
```

### Running the Application
```shell
# Basic connection test
./conntest check --dsn "your://database/uri" --retry-times 0 --tags 'tag=value'

# Examples:
./conntest check --dsn "snowflake://user:pass@host.snowflakecomputing.com/db" --tags 'env=test'
./conntest check --dsn "bigquery://:@project-id/dataset" --retry-times 0 --tags 'env=prod'
```

## Architecture

### Core Structure
- `main.go`: Entry point that calls cmd.Execute()
- `cmd/`: Cobra CLI command definitions
  - `root.go`: Root command setup
  - `check.go`: Main "check" command implementation with DSN parsing and tag handling
- `pkg/`: Core business logic
  - `conntest.go`: Connection testing logic with database-specific handling
  - `types.go`: Event and Result type definitions for JSON output
  - `databricks.go`: Databricks-specific connection handling

### Key Patterns
- Uses Cobra for CLI structure with flags for DSN, retry-times, and tags
- Database connections handled via `github.com/xo/dburl` for unified DSN parsing (v0.23.8+)
- Different connection strategies for BigQuery (uses GORM) vs other databases (uses database/sql)
- Retry logic implemented with `github.com/avast/retry-go/v4` (v4.6.1+)
- Structured JSON output using Event/Result types with UUIDs and timestamps
- All database drivers imported as blank imports for registration
- Logging to stderr with timestamps, JSON results to stdout
- Custom Databricks scheme registration with updated API compatibility

### Supported Databases
- Snowflake: Uses gosnowflake driver v1.15.0+ with information_schema queries
- BigQuery: Uses GORM with bigquery driver, connection-only testing
- PostgreSQL: Uses lib/pq driver v1.10.9+ with information_schema queries  
- Databricks: Uses databricks-sql-go driver v1.8.0+ with simple SELECT queries

### Testing
- Unit tests: `*_test.go` files with `-test.short` flag
- Integration tests: `*_integration_test.go` files requiring full test run
- Test files include database connection testing and result formatting validation