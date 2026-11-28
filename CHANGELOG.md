# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `CHANGELOG.md`, with entries back-filled from the git history, GitHub releases and pull requests.

## [0.7.2] - 2026-11-06

### Changed

- Moved from `github.com/snowflakedb/gosnowflake` v1.19.1 to `gosnowflake/v2` v2.4.0. Only the import path changed; the `snowflake://` DSN scheme works as before.
- Bumped `golang.org/x/crypto` from v0.56.0 to v0.58.0 and other transitive `golang.org/x/*` modules.
- Bumped `github.com/gabriel-vasile/mimetype` from v1.4.15 to v1.4.16, because v1.4.15 was retracted.

### Security

- Bumped Go from 1.27.0 to 1.27.3 to fix GO-2026-6599 and GO-2026-6600 (`html/template`); GO-2026-6603, GO-2026-6605, GO-2026-6609, GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6613 and GO-2026-6617 (`net/http`); GO-2026-6604 (`os`, Windows only); GO-2026-6607 (`crypto/tls`); GO-2026-6608 (`net/textproto`); and GO-2026-6809 (`crypto/x509`).
- Bumped `golang.org/x/net` from v0.58.0 to v0.62.0 to fix the HTTP/2 advisories above, GO-2026-6860 (CVE-2026-78661) and GO-2026-6866 (CVE-2026-97045).
- The gosnowflake v2 move fixes GO-2026-6665 (CVE-2026-86597, GHSA-qqj6-54q6-cxv6): the v1 driver writes sensitive information to logs and has no fixed release.
- Bumped `github.com/klauspost/compress` from v1.20.0 to v1.20.3 to fix GO-2026-6865 (CVE-2026-116591).
- Bumped `github.com/andybalholm/brotli` from v1.2.3 to v1.2.8 to fix GO-2026-6888 (CVE-2026-111497).

## [0.7.1] - 2026-09-25

### Changed

- Pinned all GitHub Actions to full commit SHAs instead of mutable version tags.
- Bumped `release-drafter` from v6 to v7.7.0 to move off the Node 20 runtime.

### Fixed

- The Postgres integration test now waits until the database is fully ready, fixing a flaky `TestCheck` (`pq: the database system is starting up`).
- Updated the release-drafter config to the v7 schema.

### Security

- Bumped `testcontainers-go` (test-only) from v0.39.0 to v0.44.0, which replaces `github.com/docker/docker` with the `github.com/moby/moby` modules. This fixes GHSA-hfg8-hc9c-6c3h, GHSA-rg2x-37c3-w2rh and GHSA-x86f-5xw2-fm2r (high), GHSA-vp62-88p7-qqf5 and GHSA-pxq6-2prw-chj9 (medium), and GHSA-8wmf-6v46-5gfg (low).

## [0.7.0] - 2026-09-07

### Changed

- Raised the Go version from 1.25 to 1.27 in `go.mod`, the Dockerfile and CI, as required by newer `golang.org/x/*` modules.

### Removed

- The Snyk workflow.

### Security

- Bumped dependencies to fix 17 of the 18 CVEs Docker Scout reported in the image (8 high, 6 medium, 1 low, 3 unspecified), in `x/crypto`, `grpc`, `go-git`, `apache/thrift`, `x/net`, `x/text`, `otel`, `aws-sdk-go-v2` and `klauspost/compress`.

## [0.6.1] - 2026-06-24

### Added

- `Dockerfile` that builds a Docker image on the distroless `static-debian12:nonroot` base.

### Changed

- Raised the Go version from 1.24 to 1.25.
- Bumped GitHub Actions to move off the deprecated Node 20 runtime.

### Security

- Bumped dependencies to fix critical and high CVEs.

## [0.6.0] - 2025-10-24

### Added

- Git repository connection checks over SSH, using `git+ssh://` DSNs with a `keyfile` parameter (#19).

### Security

- Bumped `golang.org/x/crypto` to v0.43.0 and `testcontainers-go` to v0.39.0.

## [0.5.0] - 2025-09-02

### Added

- Checks for multiple DSNs in a single run (#20).

## [0.4.0] - 2025-07-29

### Added

- Databricks OAuth machine-to-machine authentication (#17).

### Changed

- Raised the Go version to 1.24 (#14).
- Moved to the standard release tooling flow (#15).

### Fixed

- Databricks connection test (#16).
- Flaky test when setting the tags map (#18).

### Security

- Bumped dependencies to fix known vulnerabilities (#14).

## [0.3.0] - 2024-04-18

### Added

- BigQuery connection test (#11).

## [0.2.5] - 2023-04-05

### Added

- More logging for connection and query attempts (#8).

## [0.2.4] - 2022-11-07

### Added

- Copyright headers (#6).
- Snyk scanning (#4).

### Fixed

- Snyk GitHub Action (#7).

## [0.2.3] - 2022-10-07

### Fixed

- Register the Databricks DSN scheme (#3).

## [0.2.2] - 2022-10-05

### Added

- Databricks support (#2).

## [0.2.1] - 2022-09-09

### Added

- PR labeler and release drafter config (#1).

## [0.2.0] - 2022-09-07

### Changed

- Changed the exposed CLI flags.

## [0.1.1] - 2022-09-07

### Changed

- Changed the exposed CLI flags.

## [0.1.0] - 2022-09-07

### Added

- Initial release.

[unreleased]: https://github.com/snowplow-devops/conntest/compare/v0.7.2...HEAD
[0.7.2]: https://github.com/snowplow-devops/conntest/compare/v0.7.1...v0.7.2
[0.7.1]: https://github.com/snowplow-devops/conntest/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/snowplow-devops/conntest/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/snowplow-devops/conntest/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/snowplow-devops/conntest/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/snowplow-devops/conntest/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/snowplow-devops/conntest/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/snowplow-devops/conntest/compare/v0.2.5...v0.3.0
[0.2.5]: https://github.com/snowplow-devops/conntest/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/snowplow-devops/conntest/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/snowplow-devops/conntest/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/snowplow-devops/conntest/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/snowplow-devops/conntest/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/snowplow-devops/conntest/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/snowplow-devops/conntest/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/snowplow-devops/conntest/releases/tag/v0.1.0
