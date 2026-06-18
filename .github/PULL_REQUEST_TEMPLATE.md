## Description

Please include a summary of the change, the rationale behind it, and any background context.

Fixes # (issue number or link)

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update (formatting, grammar, or README/wiki improvements)
- [ ] Refactoring (no functional changes, code formatting, optimization)

## Verification & Testing

Explain how you tested the changes (e.g., manual validation, unit/integration tests). Propose commands to verify the change.

```bash
# Example test command
go test -v ./...
```

## Checklist

- [ ] My code follows the Go style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented on my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] I have run the linters locally via `make lint` and fixed any warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] I have regenerated mocks using `make generate-mock` (if I modified interface definitions in `config/reader.go`, `cmd/form_builder/form_collector.go`, or `cmd/create.go`)
- [ ] My changes generate no new warnings or build errors
