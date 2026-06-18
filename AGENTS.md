# Inzibat — Agent Guide

## Entrypoints

- `main.go`: starts Fiber server directly when called with 0–1 args; dispatches to Cobra (`cmd/`) for 2+ args.
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
config/ → router/ (worker pool, concurrency channels) → handler/
  handler/endpoint.go  — mock routes (fakeResponse)
  handler/client.go    — proxy routes (requestTo) with circuit breaker
server/                — Fiber bootstrap, graceful shutdown
```

- Routes with both `fakeResponse` and `requestTo` register **both** handlers.
- Circuit breaker: proxy routes only, tripped on network errors + 5xx (not 4xx). Global + per-route override.

## Config

- Formats: JSON, TOML, YAML (via koanf parsers). Default extension `.json`.
- **Precedence**: `--config` flag > `INZIBAT_CONFIG_FILE` env > `inzibat.json` in CWD > `~/.inzibat.config.json` (with `--global`/`-g`)
- Default server port: `8080`, concurrency: `runtime.GOMAXPROCS(3)` (configurable)
- `healthCheckRoute: true` auto-adds `GET /health → 200`

## CLI

- `inzibat start` (aliases: `s`, `server`, `start-server`)
- `inzibat create` (alias: `c`) — interactive form builder for routes
- `inzibat list` (aliases: `ls`, `l`)

## Gotchas

- `.gitignore` excludes `*.json` root-wide; only `examples/*.json` is tracked. Root `inzibat.json` is local-only.
- `*.html` and `coverage.out` are gitignored (coverage artifacts).
- `go.mod` requires `go 1.26` — ensure toolchain matches.
- Linter disables **all default linters**; only those listed in `.golangci.yml` run.
- `goimports` is the formatter (enforced via golangci config).

## Mocks

Run `make generate-mock` after changing interfaces in:
- `config/reader.go`
- `cmd/form_builder/form_collector.go`
- `cmd/create.go`

Mock files (`*_mock.go`) are excluded from SonarCloud coverage and linting.

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **inzibat** (667 symbols, 2879 relationships, 55 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> Index stale? Run `node .gitnexus/run.cjs analyze` from the project root — it auto-selects an available runner. No `.gitnexus/run.cjs` yet? `npx gitnexus analyze` (npm 11 crash → `npm i -g gitnexus`; #1939).

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows. For regression review, compare against the default branch: `detect_changes({scope: "compare", base_ref: "master"})`.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `context({name: "symbolName"})`.

## Never Do

- NEVER edit a function, class, or method without first running `impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `rename` which understands the call graph.
- NEVER commit changes without running `detect_changes()` to check affected scope.

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/inzibat/context` | Codebase overview, check index freshness |
| `gitnexus://repo/inzibat/clusters` | All functional areas |
| `gitnexus://repo/inzibat/processes` | All execution flows |
| `gitnexus://repo/inzibat/process/{name}` | Step-by-step execution trace |

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
