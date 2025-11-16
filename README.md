# Inzibat 🪖

[![Release Version](https://img.shields.io/github/v/release/Lynicis/inzibat?label=version)](https://github.com/Lynicis/inzibat/releases)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=lynicis_inzibat&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=lynicis_inzibat)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=lynicis_inzibat&metric=coverage)](https://sonarcloud.io/summary/new_code?id=lynicis_inzibat)
[![Go Version](https://img.shields.io/badge/go-1.25.4-blue)](https://golang.org)
[![License](https://img.shields.io/github/license/Lynicis/inzibat)](LICENSE)

**Inzibat** (from Turkish, meaning "Military Police") is a lightweight, fully-customizable HTTP mock server designed for microservices testing and development. Built in Go and powered by [Fiber](https://gofiber.io/), it provides a fast and simple way to simulate downstream services through declarative configuration files.

Perfect for frontend development, backend integration testing, and CI/CD pipelines—configure your mock responses in JSON, TOML, or YAML without writing a single line of server code.

---

## 🧭 Table of Contents

- [Inzibat 🪖](#inzibat-)
  - [🧭 Table of Contents](#-table-of-contents)
  - [✨ Key Features](#-key-features)
  - [🎯 Why Inzibat?](#-why-inzibat)
  - [🛠️ Installation](#️-installation)
    - [From Releases (Recommended)](#from-releases-recommended)
    - [From Source](#from-source)
      - [Option 1: Quick Install](#option-1-quick-install)
      - [Option 2: Build from Clone](#option-2-build-from-clone)
  - [🚀 Quick Start](#-quick-start)
    - [Step 1: Create a Configuration File](#step-1-create-a-configuration-file)
    - [Step 2: Start the Server](#step-2-start-the-server)
    - [Step 3: Test It](#step-3-test-it)
  - [💻 CLI Commands](#-cli-commands)
    - [Start Server](#start-server)
    - [Create Routes](#create-routes)
    - [List Routes](#list-routes)
    - [Command Aliases](#command-aliases)
  - [🧪 Testing](#-testing)
  - [📝 Configuration](#-configuration)
    - [Basic Configuration Structure](#basic-configuration-structure)
    - [Route Types](#route-types)
  - [🤝 Contributing](#-contributing)
    - [Getting Started](#getting-started)
    - [Guidelines](#guidelines)
    - [Reporting Issues](#reporting-issues)
  - [📜 License](#-license)

## ✨ Key Features

- Lightweight HTTP mock server implemented in Go
- Config-driven (JSON, TOML, YAML) for easy scenario definition
- Fast — built on top of [Fiber](https://gofiber.io/) (which uses fasthttp)
- Simple, declarative API for defining routes and responses

## 🎯 Why Inzibat?

- **For Frontend Teams:** Get predictable API responses for your UI development without waiting for the backend.
- **For Backend Teams:** Isolate your service during integration testing by mocking downstream dependencies.
- **For CI Pipelines:** Run reliable end-to-end tests by simulating third-party APIs.
- **No-Code Scenarios:** Implement complex mock behavior without writing a single line of server code.

## 🛠️ Installation

### From Releases (Recommended)

This is the easiest way to get `inzibat` for most users.

1. Go to the [**Releases Page**](https://github.com/Lynicis/inzibat/releases).
2. Download the archive matching your OS and architecture (e.g., `inzibat_linux_amd64.tar.gz`).
3. Extract the archive and move the `inzibat` binary to a directory in your system's `PATH`.

```bash
# Example for Linux/macOS
tar -xzf inzibat_linux_amd64.tar.gz
sudo mv inzibat /usr/local/bin/
```

### From Source

If you have Go 1.25+ installed, you can build `inzibat` from source.

#### Option 1: Quick Install

```bash
go install github.com/Lynicis/inzibat@latest
```

This installs the latest version to your `$GOPATH/bin` directory.

#### Option 2: Build from Clone

For development or custom builds:

```bash
git clone https://github.com/Lynicis/inzibat.git
cd inzibat
go build -o inzibat .
```

This creates a local `inzibat` binary in the project directory.

---

## 🚀 Quick Start

Get a mock server running in under 30 seconds.

### Step 1: Create a Configuration File

Create a file named `inzibat.json` (or `inzibat.yml`, `inzibat.toml`) in your current directory:

```yaml
# inzibat.yml
port: 8080
routes:
  - path: /api/hello
    method: GET
    response:
      status_code: 200
      headers:
        Content-Type: application/json
      body: '{"message": "Hello, World!"}'
```

### Step 2: Start the Server

Run Inzibat:

```bash
inzibat start
# or use the short alias
inzibat s
```

The server will start on port 8080 (or the port specified in your config).

### Step 3: Test It

In another terminal, send a request:

```bash
curl http://localhost:8080/api/hello
```

You should receive:

```json
{"message": "Hello, World!"}
```

## 💻 CLI Commands

Inzibat provides a powerful CLI for managing routes and starting the server interactively.

### Start Server

Start the mock server with the `start` command:

```bash
# Start with default config (inzibat.json in current directory)
inzibat start

# Use the short alias
inzibat s

# Specify a custom config file
inzibat start --config /path/to/config.yml
inzibat start -c /path/to/config.yml
```

**Configuration Precedence:**

The server reads configuration in the following order:

1. File specified by the `--config` / `-c` flag
2. File specified by the `INZIBAT_CONFIG_FILE` environment variable
3. `inzibat.json` in the current working directory

### Create Routes

Create new routes interactively using the `create` command:

```bash
# Launch interactive route creation
inzibat create

# Or use the short alias
inzibat c
```

The interactive form guides you through:

- Setting the route path and HTTP method
- Choosing between mock responses or proxy routes
- Configuring response status codes, headers, and body
- Setting up proxy targets for client routes

Routes are automatically saved to your global config file (`~/.inzibat.json`).

### List Routes

View all configured routes:

```bash
# List all routes
inzibat list

# Or use the short aliases
inzibat ls
inzibat l
```

This displays all routes from your configuration file in a structured, easy-to-read format.

### Command Aliases

For convenience, all commands have shorter aliases:

| Command | Aliases |
|---------|---------|
| `start` | `start-server`, `server`, `s` |
| `create` | `create-route`, `c` |
| `list` | `list-routes`, `ls`, `l` |

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run tests for a specific package
go test ./handler -v
```

## 📝 Configuration

Inzibat supports configuration files in multiple formats: JSON, TOML, and YAML. The configuration file defines the server port and routes.

### Basic Configuration Structure

```yaml
port: 8080
routes:
  - path: /api/users
    method: GET
    response:
      status_code: 200
      headers:
        Content-Type: application/json
      body: '{"users": []}'
```

### Route Types

- **Mock Routes**: Return predefined responses with custom status codes, headers, and body
- **Proxy Routes**: Forward requests to another service (useful for development)

For more examples and advanced configuration options, see the [documentation](https://github.com/Lynicis/inzibat/wiki).

## 🤝 Contributing

Contributions are welcome! We appreciate your help in making Inzibat better.

### Getting Started

1. **Fork the repository** and clone your fork
2. **Create a feature branch** (`git checkout -b feature/my-new-feature`)
3. **Make your changes** and add tests for any new behavior
4. **Run the tests** to ensure everything passes:

   ```bash
   go test ./... -v
   ```

5. **Commit your changes** with clear, descriptive messages
6. **Push to your fork** and open a Pull Request

### Guidelines

- Follow Go best practices and conventions
- Add tests for new features and bug fixes
- Update documentation as needed
- Keep commits focused and atomic
- Write clear commit messages

### Reporting Issues

For bug reports or feature requests, please [open an issue](https://github.com/Lynicis/inzibat/issues) with:

- A clear description of the problem or feature request
- Steps to reproduce (for bugs)
- Expected vs. actual behavior
- Environment details (OS, Go version, etc.)

## 📜 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
