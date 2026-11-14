# Contributing

Issues and pull requests are welcome.

## Development setup

- Go 1.26+
- `make lint` — golangci-lint (CI runs the same)
- `go test ./...` — test suite
- `make generate-mock` — regenerate mocks after changing a mocked interface

## Pull requests

1. Fork and create a feature branch
2. Add tests for new behaviour; keep commits focused
3. Ensure `make lint` and `go test ./...` pass
4. Open the PR with a clear description

## Reporting issues

Include what happened, steps to reproduce, expected vs actual behaviour, and your OS and Inzibat version.

Security vulnerabilities: follow [SECURITY.md](SECURITY.md), not a public issue.
