# conntest

[![Actions Status][actions-image]][actions] [![Go Report Card][goreport-image]][goreport] [![Release][release-image]][releases] [![License][license-image]][license]

## Overview

Conntest is a command-line utility for validating connections to Snowplow-supported destinations.

## Running

To check your database connection, run:
```shell
conntest check --driver snowflake --dsn your://database/uri --retry-times 0 --tags 'aTag=value;anotherTag=value'
```
### Example

```shell
$ conntest check --driver snowflake --dsn username:password@abcdefg-ab01234.snowflakecomputing.com/snowplow --tags=aTag=value

# Successful response
{
  "id": "6504878f-e1e5-4108-b902-c125d5a75ee4",
  "name": "fabric:warehouse-connection-check",
  "version": 1,
  "emittedBy": "conntest",
  "timestamp": "2023-07-10T17:00:34.793223+02:00",
  "data": {
    "complete": true,
    "messages": [],
    "tags": {
      "aTag": "value"
    },
    "attempts": 1
  }
}

# Failure
{
  "id": "7025ed2b-b0e9-434c-81ea-4a9633d681bd",
  "name": "fabric:warehouse-connection-check",
  "version": 1,
  "emittedBy": "conntest",
  "timestamp": "2023-07-10T17:01:08.855494+02:00",
  "data": {
    "complete": false,
    "messages": [
      "All attempts fail:\n#1: can't query snowflake: 390100 (08004): Incorrect username or password was specified."
    ],
    "tags": {
      "aTag": "value"
    },
    "attempts": 1
  }
}
```


## Development

This repo uses nix to provide [reproducible development environment](https://nixos.org/guides/ad-hoc-developer-environments.html). To make use of the provided setup:

1. Install `nix`:
```shell
sh <(curl -L https://nixos.org/nix/install)
```
2. Enable experimental flags
``` shell
mkdir -p ~/.config/nix && echo 'experimental-features = nix-command flakes' > ~/.config/nix/nix.conf
```
3. Enter development environment
```shell
nix develop
```

> **Note**
> If you want the convenience of getting the development environment upon `cd` into directory use [direnv](https://direnv.net)

4. Develop
```
# format
go fmt ./...
# build
go build
# test
go test -v ./... -test.short
# test with integration tests
go test -v ./...
```


### Copyright and license

Conntest is copyright 2022-2022 Snowplow Analytics Ltd.

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
