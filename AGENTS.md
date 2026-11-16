# Inzibat — Agent Guide

## Entrypoints

- `main.go`: calls `cmd.Execute()` only. The root Cobra command starts the server; `inzibat start` is an alias for the same runner.
- Server can be invoked standalone (bare `inzibat` or `go run .`) or via CLI (`inzibat start`).

## Build & Test

| Command | What |
|---------|------|
| `go build -o inzibat .` | Build binary |
| `go test ./... -v` | All tests (testify) |
| `go test ./handler -v` | Single package |
| `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out` | Coverage + HTML |
| `make lint` | golangci-lint v2 (strict: cyclop ≤11, funlen ≤100, lll) |
| `make generate-mock` | Regenerate `go.uber.org/mock` mocks (3 sources) |

CI order: lint → build → test → SonarCloud coverage → gosec security.

## Architecture

```
config/ → router/ (sequential route registration) → handler/
  handler/endpoint.go  — mock routes (fakeResponse)
  handler/client.go    — proxy routes (requestTo) with circuit breaker
server/                — Fiber bootstrap, graceful shutdown
```

- Routes with both `fakeResponse` and `requestTo` serve `requestTo`; the loader warns and the mock response is dropped.
- Circuit breaker: proxy routes only, tripped on network errors + 5xx (not 4xx). Global + per-route override.

## Config

- Formats: JSON, TOML, YAML (via koanf parsers). Default extension `.json`.
- **Precedence**: `--config` flag > `INZIBAT_CONFIG_FILE` env > `inzibat.json` in CWD > `~/.inzibat.config.json` (with `--global`/`-g`)
- Default server port: `8080`
- `healthCheckRoute: true` auto-adds `GET /health → 200`

## CLI

- `inzibat start` (aliases: `s`, `server`, `start-server`)
- `inzibat create` (alias: `c`) — interactive form builder for routes
- `inzibat list` (aliases: `ls`, `l`)

## Gotchas

- `.gitignore` excludes `*.json` root-wide; only `examples/*.json` is tracked. Root `inzibat.json` is local-only.
- `*.html` and `coverage.out` are gitignored (coverage artifacts).
- `go.mod` requires `go 1.27` — ensure toolchain matches. `tool` directives pin golangci-lint, gosec, mockgen, govulncheck.
- Linter disables **all default linters**; only those listed in `.golangci.yml` run.
- `goimports` is the formatter (enforced via golangci config).

## Mocks

Run `make generate-mock` after changing interfaces in:
- `config/reader.go`
- `cmd/formbuilder/form_collector.go`

Mock files (`*_mock.go`) are excluded from SonarCloud coverage and linting.

## Graphify — Code Intelligence

This project is indexed by Graphify; the knowledge graph lives in `graphify-out/graph.json` (god nodes, communities, EXTRACTED/INFERRED edge audit trail).

- Query the graph instead of grepping when exploring unfamiliar code: `graphify query "<question>"`, `graphify explain "<node>"`, `graphify path "A" "B"`.
- Index stale? Run `graphify --update` from the project root (AST re-extraction is free and fast for Go).
- Full rebuild: `graphify .` — see `graphify-out/GRAPH_REPORT.md` for the audit report.
- `GRAPH_REPORT.md` and the `.graphify_*` sidecars are committed; `graph.json` itself is local-only (root-wide `*.json` ignore).
