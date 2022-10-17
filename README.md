# conntest

[![Actions Status][actions-image]][actions] [![Go Report Card][goreport-image]][goreport] [![Release][release-image]][releases] [![License][license-image]][license]

## Overview

Conntest is a simple utility that checks connections to Snowplow-supported destinations written in Go.

## Quick start

Install `nix` using `brew install nix`.

To check the format of the go code :-

```nix
nix develop --extra-experimental-features nix-command --extra-experimental-features flakes -c go fmt ./...
```

To build the go code :-

```nix
nix develop --extra-experimental-features nix-command --extra-experimental-features flakes -c go build
```

To run short tests (no Docker required) :-

```nix
nix develop --extra-experimental-features nix-command --extra-experimental-features flakes -c go test -v ./... -test.short
```

To run longer tests (running Docker required) :-

```nix
nix develop --extra-experimental-features nix-command --extra-experimental-features flakes -c go test -v ./...
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