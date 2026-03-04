# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Inzibat is a lightweight HTTP mock server written in Go that provides config-driven mocking and proxying capabilities. Built on Fiber (fasthttp), it allows users to define mock API responses or proxy requests through declarative JSON/TOML/YAML configuration files.

## Development Commands

### Testing
```bash
# Run all tests
go test ./... -v

# Run tests with coverage report
make coverage

# Run specific package tests
go test ./handler -v
go test ./config -v
```

### Building
```bash
# Build binary
go build -o inzibat .

# Install from source
go install github.com/lynicis/inzibat@latest
```

### Code Quality
```bash
# Run linter
make lint

# Generate mock files (after modifying interfaces)
make generate-mock
```

### Running the Server
```bash
# Start with default config (inzibat.json in current directory)
inzibat start
# or: inzibat s

# Use custom config file
inzibat start --config /path/to/config.yml
# or: inzibat start -c /path/to/config.yml

# Use global config (~/.inzibat.config.json)
inzibat start --global
# or: inzibat start -g
```

### CLI Commands
```bash
# Create route interactively (uses current directory inzibat.json)
inzibat create  # or: inzibat c

# Create route in custom config file
inzibat create --config /path/to/config.json
# or: inzibat create -c /path/to/config.json

# Create route in global config (~/.inzibat.config.json)
inzibat create --global
# or: inzibat create -g

# List all configured routes (from current directory inzibat.json)
inzibat list    # or: inzibat ls, inzibat l

# List routes from custom config file
inzibat list --config /path/to/config.json
# or: inzibat list -c /path/to/config.json

# List routes from global config (~/.inzibat.config.json)
inzibat list --global
# or: inzibat list -g
```

## Architecture

### Core Flow
1. **Configuration Loading** (`config/`): Reads config from JSON/TOML/YAML using koanf, validates using go-playground/validator
2. **Server Initialization** (`server/`): Creates Fiber app, initializes handlers
3. **Route Registration** (`router/`): Concurrent route setup using worker pool pattern (controlled by `concurrency` config)
4. **Request Handling** (`handler/`): Two handler types process requests:
   - **EndpointHandler**: Returns mock responses (status, headers, body)
   - **ClientHandler**: Proxies requests to external services using reflection to call HTTP methods

### Configuration Precedence
The server loads config in this order:
1. File specified by `--config` / `-c` flag
2. File from `INZIBAT_CONFIG_FILE` environment variable
3. `inzibat.json` in current working directory
4. Global config (`~/.inzibat.config.json`) if `--global` / `-g` flag is used

### Route Types
- **Mock Routes**: Return predefined responses via `FakeResponse` (body/bodyString, headers, statusCode)
- **Proxy Routes**: Forward requests via `RequestTo` (host, path, method, headers, body)

Routes can have both `FakeResponse` and `RequestTo`, but at least one is required.

### Concurrency Model
Routes are registered concurrently using a worker pool:
- Workers: Controlled by `concurrency` config field (defaults to `runtime.GOMAXPROCS(3)`)
- Pattern: Channel-based work distribution in `router/router.go:CreateRoutes()`
- Each worker processes routes from shared channel using `processRoute()`

### HTTP Client (Proxy) Implementation
The `ClientHandler` uses reflection to dynamically call HTTP methods:
1. Converts method name (GET/POST/etc.) to title case
2. Uses `reflect.ValueOf().MethodByName()` to call the appropriate `Client` method
3. Passes URL, headers, and body (body excluded for GET requests)

## Key Files

### Entry Point
- `main.go`: Entry point, starts server or executes CLI command

### Configuration
- `config/model.go`: Core data structures (`Cfg`, `Route`, `RequestTo`, `FakeResponse`)
- `config/config.go`: Config loading logic with precedence rules
- `config/reader.go`: Strategy pattern for reading different config formats
- `config/file_loader.go`: Koanf-based file parsers (JSON, TOML, YAML)

### Server & Routing
- `server/server.go`: Server lifecycle (start, graceful shutdown)
- `router/router.go`: Concurrent route registration using worker pool

### Handlers
- `handler/endpoint.go`: Returns mock responses
- `handler/client.go`: Proxies requests using reflection-based HTTP method dispatch

### CLI
- `cmd/root.go`: Cobra root command setup
- `cmd/startServer.go`: Server start command with config flag handling
- `cmd/create.go`: Interactive route creation using charmbracelet/huh forms
- `cmd/list.go`: Displays routes in formatted table

## Testing Patterns

All tests follow AAA (Arrange, Act, Assert) pattern with explicit comments:

```go
func TestFunction(t *testing.T) {
    // Arrange
    setup := prepareTestData()

    // Act
    result := functionUnderTest(setup)

    // Assert
    assert.Equal(t, expected, result)
}
```

### Testing Guidelines
- Use `github.com/stretchr/testify/assert` for assertions
- Use `go.uber.org/mock` for mocks (generated via `make generate-mock`)
- Test both happy paths and error scenarios
- For handlers: Use `net/http/httptest` or Fiber test utilities
- Constructors should panic when critical dependencies are nil (test this)
- Handlers should return errors, not panic (Fiber handles error responses)

## Code Conventions

### Error Handling
- Always return errors explicitly (no silent failures)
- Use custom errors from `config/error.go` when appropriate
- Log errors with context: `zap.L().Error("message", zap.Error(err), zap.String("key", value))`

### Logging
- Use `zap.L()` globally (initialized in `log/log.go`)
- Include structured fields for context
- Log levels: Fatal (unrecoverable), Error (handled errors), Info (important events)

### Fiber Handler Pattern
```go
func handler(ctx *fiber.Ctx) error {
    // Process request
    if err != nil {
        return err  // Fiber handles error response
    }

    return ctx.Status(200).JSON(data)
}
```

### Config Validation
- Use struct tags: `validate:"required,gt=0,oneof=GET POST"`
- Validation happens in `config/config.go:Read()` via validator.Struct()
- Custom validators in `cmd/form_builder/validators.go`

## Mock Generation

Mock files are generated using `go.uber.org/mock`:
```bash
make generate-mock
```

Generates:
- `config/reader_mock.go`: Mocks `ReaderStrategy` interface
- `cmd/form_builder/form_runner_mock.go`: Mocks form execution
- `cmd/create_mock.go`: Mocks route creation functions

After changing interfaces, regenerate mocks before running tests.

## Dependencies

- **Fiber v2**: Web framework (fasthttp-based)
- **Koanf v2**: Config management (JSON/TOML/YAML)
- **Validator v10**: Struct validation
- **Zap**: Structured logging
- **Charmbracelet/huh**: Interactive CLI forms
- **Cobra**: CLI command framework
- **Testify**: Testing utilities (assert, require)
- **go.uber.org/mock**: Mock generation

## Common Development Workflows

### Adding a New Route Type
1. Update `config/model.go` with new fields and validation tags
2. Modify `handler/` to implement new behavior
3. Update `router/router.go` route processing logic
4. Add tests for new functionality
5. Update config examples in `examples/`

### Adding a New CLI Command
1. Create command file in `cmd/` (follow `cmd/create.go` pattern)
2. Register command in `cmd/root.go` via `init()`
3. Use charmbracelet/huh for interactive forms
4. Write tests in `cmd/*_test.go`

### Debugging Configuration Issues
- Check config precedence (flag > env var > cwd file > global)
- Verify validation tags in `config/model.go`
- Review `config/config.go:Read()` for defaults and transformations
- Health check route is auto-added if `healthCheckRoute: true`
