# conntest

[![Actions Status][actions-image]][actions] [![Go Report Card][goreport-image]][goreport] [![Release][release-image]][releases] [![License][license-image]][license]

## Overview

Conntest is a command-line utility for validating connections to Snowplow-supported destinations.

## Running

To check your database connection, run:

```shell
conntest check --dsn your://database/uri --retry-times 0 --tags 'aTag=value;anotherTag=value'
```

### Example

For Snowflake (success):-

```shell
conntest check --dsn snowflake://username:password@snowplow.snowflakecomputing.com/database --tags 'aTag=value'

{"id":"be863f45-9d16-4271-94f1-c019df7d8d49","name":"fabric:warehouse-connection-check","version":1,"emittedBy":"conntest","timestamp":"2024-04-08T11:50:05.816247+01:00","data":{"host":"snowplow.snowflakecomputing.com","complete":true,"messages":[],"tags":{"aTag":"value"},"attempts":1}}
```

For Snowflake (failure):-

```shell
conntest check --dsn snowflake://lorem:ipsum@abcdefg-ab01234.snowflakecomputing.com/lorem --tags 'aTag=value'

{"id":"9ef56c9b-8f30-4ef3-8e6a-0b063e938ce4","name":"fabric:warehouse-connection-check","version":1,"emittedBy":"conntest","timestamp":"2024-04-08T11:45:38.086054+01:00","data":{"host":"abcdefg-ab01234.snowflakecomputing.com","complete":false,"messages":["260008 (08004): failed to connect to db. verify account name is correct. HTTP: 403, URL: https://abcdefg-ab01234.snowflakecomputing.com:443/session/v1/login-request?databaseName=lorem\u0026requestId=227cfd45-7335-44e1-6495-bfc57a0ca842\u0026request_guid=5b8b3eff-ee45-4c48-7498-e6b2ba707082"],"tags":{"aTag":"value"},"attempts":1}}
```

For BigQuery (success):-

```shell
conntest check --dsn bigquery://:@engineering-sandbox/testantonis_derived --retry-times 0 --tags 'aTag=value;anotherTag=value'

{"id":"8be8275c-1aa4-4907-8e28-c027dbe58ad7","name":"fabric:warehouse-connection-check","version":1,"emittedBy":"conntest","timestamp":"2024-04-08T11:36:50.593231+01:00","data":{"host":"engineering-sandbox","complete":true,"messages":[],"tags":{"aTag":"value","anotherTag":"value"},"attempts":0}}
```

For BigQuery (failure):-

```shell
conntest check --dsn bigquery://:@engineering-sandbox/testantonis_invalid --retry-times 0 --tags 'aTag=value;anotherTag=value'

{"id":"e50cea6f-adc1-4399-af88-9e6d439f5f16","name":"fabric:warehouse-connection-check","version":1,"emittedBy":"conntest","timestamp":"2024-04-08T11:38:22.039055+01:00","data":{"host":"engineering-sandbox","complete":false,"messages":["googleapi: Error 404: Not found: Dataset engineering-sandbox:testantonis_invalid, notFound"],"tags":{"aTag":"value","anotherTag":"value"},"attempts":0}}
```

## Development

### Prerequisites

- Go 1.24+ 
- Docker (for integration tests)
- Make

### Quick Start

Assuming git, Go, and Make are installed:

```bash
host> git clone https://github.com/snowplow-devops/conntest
host> cd conntest
host> make test
host> make build
```

### Available Make Targets

#### Building
```bash
# Build for local development
make build

# Build all cross-platform binaries
make all

# Build specific platforms
make cli-linux-amd64
make cli-linux-arm64
make cli-darwin-amd64
make cli-darwin-arm64
```

#### Testing
```bash
# Run unit tests with coverage
make test

# Run integration tests
make integration-test
```

#### Code Quality
```bash
# Format code
make format

# Lint code
make lint

# Tidy Go modules
make tidy

# Update dependencies
make update
```

#### Cleanup
```bash
# Remove build artifacts
make clean
```

All compiled assets are available under `build/compiled`.

**Note:** Always run `make format` before submitting any code.

**Note:** The `make test` command generates a code coverage file which can be found at `build/coverage/coverage.html`.

### Copyright and license

Conntest is copyright 2022-2025 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[actions-image]: https://github.com/snowplow-devops/conntest/workflows/ci/badge.svg
[actions]: https://github.com/snowplow-devops/conntest/actions

[release-image]: https://img.shields.io/github/v/release/snowplow-devops/conntest?style=flat&color=6ad7e5
[releases]: https://github.com/snowplow-devops/conntest/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow-devops/conntest
[goreport]: https://goreportcard.com/report/github.com/snowplow-devops/conntest
