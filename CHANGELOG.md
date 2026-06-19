# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/Lynicis/inzibat/compare/v0.4.0...HEAD
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
