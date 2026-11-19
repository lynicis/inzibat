# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.1] - 2026-11-19

### Fixed
- JSON Schema now mirrors runtime validation: required `serverPort`/`routes`, at least one route handler, method enums, `/`-prefixed paths, circuit-breaker minimums, and header values as a string or string array. Example configs point editors at it.
- Chocolatey package metadata: removed broken `raw.githubusercontent.com` `iconUrl` and added `<copyright>` matching `LICENSE`.

## [0.5.0] - 2026-11-16

### Added
- JSON Schema (`examples/inzibat.schema.json`) describing the real configuration; configs can be validated against it.
- Graphify knowledge-graph index (`graphify-out/`) replacing GitNexus, queryable with `graphify query`/`explain`/`path`.
- Headless end-to-end tests for every interactive huh form (route, mock response, client request, header/body collectors); total test coverage is now 95.5%.

### Changed
- Rewrote the README: corrected the quick start to the real config schema (`serverPort`, `fakeResponse.statusCode`, object `body`), fixed dead install URLs, stated project status, and moved contributing content to CONTRIBUTING.md.
- Whole-repository refactor: one koanf-backed config reader, sequential route registration, `formbuilder` package with a single terminal-wiring helper, one `Do` method in `client/http`, a guarded-map circuit breaker store, and recorder/CLI helpers deduplicated.
- Upgraded the Go toolchain requirement to 1.27 and tidied `go.mod` (`tool` directives pin golangci-lint, gosec, mockgen, govulncheck).

### Fixed
- Proxy routes now honour `passWithRequestHeaders` and `passWithRequestBody`, forward multi-value headers intact, and no longer send a JSON `null` body.
- HTTP client retries use real exponential backoff with jitter and only retry idempotent methods; response status is read before release.
- Circuit breakers: the half-open in-flight gauge is released, and seeded breakers get the documented defaults.
- Server: listener failure is returned as an error instead of calling `log.Fatal`.
- Config: defaults and normalization run before validation, bodyless fake responses are valid, `healthCheckRoute` survives a write/read round trip, and a route with both handlers warns and drops the mock.
- CLI: `create` refuses to write a typeless route, circuit breaker form defaults no longer leak between forms, the root command starts the server, and `record` rejects unknown export formats with bounded HTTP calls.
- Recorder skips HTTP methods a mock route cannot serve.

## [0.4.0] - 2026-06-19

### Added
- Request Recorder feature to capture incoming HTTP requests and responses flowing through the mock/proxy server.
- New `record` CLI command tree with `list`, `export`, and `clear` subcommands.
- Flag `--record`/`-r` for the `start` command to activate request recording on startup.
- Admin HTTP endpoints under `/_inzibat/recorder` to retrieve or clear recorded session data.

### Fixed
- Replaced deprecated fasthttp `VisitAll` calls with Go 1.23+ iterators (`All()`) in recorder middleware.

## [0.3.11] - 2026-06-16

### Fixed
- CI Chocolatey push parameter flag (changed from `--key` to `-k`).

## [0.3.10] - 2026-06-16

### Added
- Passed build archives as GitHub Action artifacts to `windows-latest` for Chocolatey packages.

## [0.3.9] - 2026-06-16

### Changed
- Refined archive passing for Chocolatey packaging steps.

## [0.3.8] - 2026-06-16

### Changed
- Split release into two distinct stages with native PowerShell Chocolatey packager on `windows-latest`.

## [0.3.7] - 2026-06-16

### Changed
- Unified GoReleaser config and run release on `windows-latest` to fix Chocolatey checksum mismatch.

## [0.3.6] - 2026-06-14

### Fixed
- URL template configuration in Chocolatey release-disabled mode.

## [0.3.5] - 2026-06-14

### Changed
- General repository and maintenance updates.

## [0.3.4] - 2026-06-08

### Changed
- Release enhancements.

## [0.3.3] - 2026-06-08

### Changed
- Routine maintenance release.

## [0.3.2] - 2026-06-07

### Changed
- Updated Homebrew tap release configurations.

## [0.3.1] - 2026-06-07

### Changed
- CI updates to publish Homebrew formulas instead of casks.

## [0.3.0] - 2026-06-07

### Added
- Automated publishing to Homebrew, Scoop, and NFPM packages (deb, rpm) via GoReleaser.

## [0.2.0] - 2026-06-05

### Changed
- Updated linter configuration (`golangci-lint` settings).

## [0.1.7] - 2025-12-15

### Changed
- Increased server `ReadBufferSize` for improved performance.

## [0.1.6] - 2025-12-03

### Changed
- Updated mock generation source scripts.
- Cleaned up source code comments.

## [0.1.5] - 2025-12-03

### Added
- Added global configuration file support (`~/.inzibat.config.json`).
- Enhanced CLI commands with aliases, global configuration flags (`--global` / `-g`), and listings.

## [0.1.4] - 2025-11-18

### Fixed
- Added proper Go versioning constraints to `go.mod`.

## [0.1.3] - 2025-11-18

### Changed
- Refined package configurations to make the repository installable via standard `go install`.

## [0.1.2] - 2025-11-17

### Fixed
- Code style and minor linter errors.

## [0.1.1] - 2025-11-14

### Fixed
- Resolved issues in the CI security analysis/gosec stage.

## [0.1.0] - 2025-11-11

### Changed
- Switched default JSON parsing and serialization package to high-performance library.

[Unreleased]: https://github.com/Lynicis/inzibat/compare/v0.5.1...HEAD
[0.5.1]: https://github.com/Lynicis/inzibat/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/Lynicis/inzibat/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/Lynicis/inzibat/compare/v0.3.11...v0.4.0
[0.3.11]: https://github.com/Lynicis/inzibat/compare/v0.3.10...v0.3.11
[0.3.10]: https://github.com/Lynicis/inzibat/compare/v0.3.9...v0.3.10
[0.3.9]: https://github.com/Lynicis/inzibat/compare/v0.3.8...v0.3.9
[0.3.8]: https://github.com/Lynicis/inzibat/compare/v0.3.7...v0.3.8
[0.3.7]: https://github.com/Lynicis/inzibat/compare/v0.3.6...v0.3.7
[0.3.6]: https://github.com/Lynicis/inzibat/compare/v0.3.5...v0.3.6
[0.3.5]: https://github.com/Lynicis/inzibat/compare/v0.3.4...v0.3.5
[0.3.4]: https://github.com/Lynicis/inzibat/compare/v0.3.3...v0.3.4
[0.3.3]: https://github.com/Lynicis/inzibat/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/Lynicis/inzibat/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/Lynicis/inzibat/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/Lynicis/inzibat/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/Lynicis/inzibat/compare/v0.1.7...v0.2.0
[0.1.7]: https://github.com/Lynicis/inzibat/compare/v0.1.6...v0.1.7
[0.1.6]: https://github.com/Lynicis/inzibat/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/Lynicis/inzibat/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/Lynicis/inzibat/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/Lynicis/inzibat/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/Lynicis/inzibat/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/Lynicis/inzibat/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/Lynicis/inzibat/releases/tag/v0.1.0
