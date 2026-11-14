# Inzibat 🪖

[![Release Version](https://img.shields.io/github/v/release/Lynicis/inzibat?label=version)](https://github.com/Lynicis/inzibat/releases)
[![Quality Gate](https://sonarcloud.io/api/project_badges/measure?project=lynicis_inzibat&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=lynicis_inzibat)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=lynicis_inzibat&metric=coverage)](https://sonarcloud.io/summary/new_code?id=lynicis_inzibat)
[![License](https://img.shields.io/github/license/Lynicis/inzibat)](LICENSE)

Inzibat (Turkish for "military police") is a config-driven HTTP mock and proxy server. Mock your downstream APIs from a config file — no server code to write, nothing extra to keep running. It serves the mocks, proxies the routes you don't mock, and records real traffic to replay as mocks.

**Status:** actively maintained · [v0.4.0](https://github.com/Lynicis/inzibat/releases) · pre-1.0, config format may change between minor versions · [MIT](LICENSE)

Save this as `inzibat.json`:

```json
{
  "serverPort": 8080,
  "routes": [
    {
      "method": "GET",
      "path": "/api/hello",
      "fakeResponse": {
        "statusCode": 200,
        "headers": { "Content-Type": "application/json" },
        "body": { "message": "Hello, World!" }
      }
    }
  ]
}
```

Start it and call it:

```console
$ inzibat start
{"level":"info","timestamp":"2026-11-14T13:07:38.128+0300","caller":"server/server.go:144","msg":"🫡 INZIBAT 🪖","pid":30235,"open_routes":1,"server_port":8080}

$ curl http://localhost:8080/api/hello
{"message": "Hello, World!"}
```

The same file ships as [`examples/quickstart.json`](examples/quickstart.json) — run it with `inzibat start -c examples/quickstart.json`. YAML and TOML configs work the same way.

## Why Inzibat

[WireMock](https://wiremock.org/) needs a JVM, [Prism](https://stoplight.io/open-source/prism/) needs Node, [Mockoon](https://mockoon.com/) is a desktop app. Inzibat is one self-contained Go binary — drop it into a CI image or a `docker-compose` sidecar and it runs.

Beyond static mocks, one config file gives you:

- **Proxying with a circuit breaker** — routes with `requestTo` forward to a real upstream. The breaker trips on network errors and 5xx (not 4xx), configurable globally and per route.
- **Traffic recording** — start with `--record` and Inzibat captures requests passing through, then exports them as a new Inzibat config. Record once against the real service, replay it as mocks.
- **Interactive route builder** — `inzibat create` is a terminal form that writes routes into your config file.
- **Schema-validated config** — set `$schema` to [examples/inzibat.schema.json](examples/inzibat.schema.json) for editor completion and validation, in JSON, TOML, or YAML.

What it deliberately does not do: responses are static — fixed status, headers, and body, no templating or scripted behaviour — and there is no web UI. If you need stateful request matching or response templating, WireMock is the better tool. If you want mocks inline in Go tests, use [`httptest`](https://pkg.go.dev/net/http/httptest).

## Installation

**Homebrew** (macOS, Linux — builds from source):

```bash
brew tap lynicis/tap
brew install inzibat
```

**Chocolatey** (Windows):

```powershell
choco install inzibat
```

**Scoop** (Windows):

```powershell
scoop bucket add inzibat https://github.com/lynicis/scoop-bucket.git
scoop install inzibat
```

**Binaries and Linux packages** — `tar.gz`/`zip` archives and `.deb`/`.rpm` packages for each release are on the [releases page](https://github.com/Lynicis/inzibat/releases). Asset names carry the version and platform, e.g. `inzibat_0.4.0_linux_amd64.deb`.

**From source** — needs Go 1.26+:

```bash
go install github.com/lynicis/inzibat@latest
```

## Usage

```bash
inzibat start                # serve ./inzibat.json
inzibat start -c api.yml     # any JSON, TOML, or YAML file
inzibat start -g             # use ~/.inzibat.config.json
inzibat start --record       # also capture passing traffic
inzibat create               # add a route interactively
inzibat list                 # print the routes a config defines
inzibat record list          # show captured requests
inzibat record export --format inzibat -o mocks.json   # traffic → mock config
inzibat record clear         # drop the capture
```

Commands resolve the config file in this order: `--config`/`-c`, then the `INZIBAT_CONFIG_FILE` environment variable, then `./inzibat.json`. The `-g`/`--global` flag uses `~/.inzibat.config.json` instead. Aliases: `start` = `s` `server` `start-server`, `create` = `c` `create-route`, `list` = `l` `ls` `list-routes`.

While recording, the server exposes an admin API: `GET /_inzibat/recorder/entries`, `GET /_inzibat/recorder/session`, and `POST /_inzibat/recorder/clear`.

## Configuration

The full shape is defined by [examples/inzibat.schema.json](examples/inzibat.schema.json); working configs in all three formats are under [examples/](examples/). The keys most people touch:

| Key | What it does |
|---|---|
| `serverPort` | Port to listen on (required) |
| `healthCheckRoute` | `true` adds `GET /health` returning 200 |
| `routes[].method` | One of `GET`, `POST`, `PUT`, `PATCH`, `DELETE` |
| `routes[].fakeResponse` | Static mock: `statusCode`, `headers`, `body` (object) or `bodyString` (raw string) |
| `routes[].requestTo` | Proxy: `host`, `path`, `method`, `passWithRequestBody`, `passWithRequestHeaders`, `inErrorReturn500` |
| `circuitBreaker` | Global breaker defaults; override per route with `requestTo.circuitBreaker` |

A route takes either `fakeResponse` or `requestTo`. If both are present, `requestTo` is served and the mock is dropped with a warning.

## Project status

Inzibat is under active, single-maintainer development. Releases follow [Semantic Versioning](https://semver.org/); until 1.0 the config format may change between minor versions. See [CHANGELOG.md](CHANGELOG.md) for release history and [ROADMAP.md](ROADMAP.md) for what is planned.

Support: [open an issue](https://github.com/Lynicis/inzibat/issues). Security reports follow [SECURITY.md](SECURITY.md). Contributions: [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
