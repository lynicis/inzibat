# Inzibat 🪖

[![Release Version](https://img.shields.io/github/v/release/Lynicis/inzibat?label=version)](https://github.com/Lynicis/inzibat/releases)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=lynicis_inzibat&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=lynicis_inzibat)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=lynicis_inzibat&metric=coverage)](https://sonarcloud.io/summary/new_code?id=lynicis_inzibat)
[![Go Version](https://img.shields.io/badge/go-1.25.4-blue)](https://golang.org)
[![License](https://img.shields.io/github/license/Lynicis/inzibat)](LICENSE)

Inzibat (from Turkish, meaning "Military Police") is a small, fully-customizable mock service intended for use as a lightweight HTTP mock server for microservices testing and development.

This repository provides a configurable mock server written in Go. It reads simple configuration files (JSON/TOML/YAML) and serves mock responses, allowing teams to simulate downstream services during development and integration testing.

---

## 🧭 Table of Contents

- [Inzibat 🪖](#inzibat-)
  - [🧭 Table of Contents](#-table-of-contents)
  - [✨ Key Features](#-key-features)
  - [🎯 Why Inzibat?](#-why-inzibat)
  - [🛠️ Installation](#️-installation)
    - [From Releases (Recommended)](#from-releases-recommended)
    - [From Source](#from-source)
  - [🚀 Quick Start (Hello, World!)](#-quick-start-hello-world)
  - [🧪 Testing](#-testing)
  - [🤝 Contributing](#-contributing)
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

1.  Go to the [**Releases Page**](https://github.com/Lynicis/inzibat/releases).
2.  Download the archive matching your OS and architecture (e.g., `inzibat_linux_amd64.tar.gz`).
3.  Extract the archive and move the `inzibat` binary to a directory in your system's `PATH`.

```bash
# Example for Linux/macOS
tar -xzf inzibat_linux_amd64.tar.gz
sudo mv inzibat /usr/local/bin/
```

### From Source

If you have Go (1.25+) installed, you can build `inzibat` from source.

Option 1: go install (quick)

```bash
go install github.com/Lynicis/inzibat@latest
```

Option 2: Build from Clone (for development)

```bash
git clone https://github.com/Lynicis/inzibat.git
cd inzibat
go build -o inzibat .
```

---

## 🚀 Quick Start (Hello, World\!)

Let's get a mock server running in 30 seconds.

1. Create a config file

Create a file named `config.yml`:

```yaml
# config.yml
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

1. Run Inzibat

Start the server by running the binary (ensure your config file is next to the binary):

```bash
./inzibat
```

1. Test it!

In another terminal, use curl to send a request:

```bash
curl http://localhost:8080/api/hello
```

## 🧪 Testing

This project includes unit tests. To run them:

```bash
go test ./... -v
```

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new feature branch (`git checkout -b feature/my-new-feature`).
3. Make your changes and add tests for any new behavior.
4. Run the tests (`go test ./...`).
5. Open a Pull Request describing your changes.

For bug reports or feature requests, please [open an issue](https://github.com/Lynicis/inzibat/issues) with a reproducible example.

## 📜 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
