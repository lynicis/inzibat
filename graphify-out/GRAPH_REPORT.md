# Graph Report - inzibat  (2026-11-16)

## Corpus Check
- 93 files · ~36,454 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 7 file(s) not represented in the graph (top: (none) 3, .toml 3, .properties 1)

## Summary
- 835 nodes · 2242 edges · 54 communities (36 shown, 18 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 272 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2d50cbf1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- middleware.go
- Cfg
- FormRunner
- CircuitBreakerStore
- nfpms
- with
- testing.T
- steps
- http.go
- steps
- http_test.go
- requestTo
- .golangci.yml
- .CreateHandler
- Inzibat — Agent Guide
- Changelog
- form_builder_test.go
- setupServer
- RecordedEntry
- Contributor Covenant Code of Conduct
- readFile
- hooks
- Store
- README.md
- record_test.go
- Inzibat 🪖
- GitNexus — Code Intelligence
- TestSchemaMatchesModel
- PULL_REQUEST_TEMPLATE.md
- Security Policy
- [0.4.0] - 2026-06-19
- [0.1.2] - 2025-11-17
- [0.1.4] - 2025-11-18
- [0.1.5] - 2025-12-03
- [0.1.6] - 2025-12-03
- [0.2.0] - 2026-06-05
- [0.3.0] - 2026-06-07
- [0.3.10] - 2026-06-16
- [0.3.11] - 2026-06-16
- [0.3.3] - 2026-06-08
- [0.3.4] - 2026-06-08
- [0.3.5] - 2026-06-14
- [0.3.6] - 2026-06-14
- [0.3.8] - 2026-06-16
- [Unreleased]
- github.com/lynicis/inzibat
- actions
- BoolPointer
- .CreateHandler
- store_test.go
- setupAdminApp

## God Nodes (most connected - your core abstractions)
1. `Changelog` - 24 edges
2. `Store` - 22 edges
3. `NewHttpClient()` - 22 edges
4. `Cfg` - 21 edges
5. `FormRunner` - 17 edges
6. `BoolPointer()` - 16 edges
7. `NewCircuitBreakerStore()` - 15 edges
8. `CircuitBreakerStore` - 14 edges
9. `setupServer()` - 14 edges
10. `Client` - 13 edges

## Surprising Connections (you probably didn't know these)
- `Entrypoints` --references--> `Execute()`  [INFERRED]
  AGENTS.md → cmd/root.go
- `TestClientHandler_CreateHandler_NoNullBody()` --calls--> `captureBody()`  [INFERRED]
  handler/client_test.go → recorder/middleware.go
- `ClientHandler` --references--> `Client`  [EXTRACTED]
  handler/client.go → client/http/http.go
- `ConvertToInzibatConfig()` --references--> `Cfg`  [EXTRACTED]
  recorder/converter.go → config/model.go
- `MainRouter` --references--> `Cfg`  [EXTRACTED]
  router/router.go → config/model.go

## Import Cycles
- None detected.

## Communities (54 total, 18 thin omitted)

### Community 1 - "middleware.go"
Cohesion: 0.19
Nodes (8): captureBody(), captureHeaderValues(), captureRequestHeaders(), captureResponseHeaders(), NewRecorderMiddleware(), BenchmarkMiddleware(), TestCaptureBody(), TestNewRecorderMiddleware()

### Community 2 - "Cfg"
Cohesion: 0.06
Nodes (37): routeRows(), TestRouteRows(), addHealthCheckRoute(), DefaultCfg(), Reader, NewLoader(), normalizeRoutes(), ReadOrCreateConfig() (+29 more)

### Community 3 - "FormRunner"
Cohesion: 0.08
Nodes (36): collectCircuitBreakerConfig(), collectClientRequestBody(), createClientRequestFormInternal(), createMockResponseFormInternal(), createRouteInternal(), parsePositiveInt(), TestCollectCircuitBreakerConfig(), TestCreateClientRequestFormInternal() (+28 more)

### Community 4 - "CircuitBreakerStore"
Cohesion: 0.22
Nodes (9): allowHalfOpen(), CircuitBreakerStore, openRecord(), releaseHalfOpen(), resetToClosed(), TestAllowHalfOpen(), TestCircuitBreakerStore_Allow(), CircuitBreakerRecord (+1 more)

### Community 5 - "nfpms"
Cohesion: 0.04
Nodes (47): archives, format_overrides, formats, name_template, before, hooks, builds, env (+39 more)

### Community 6 - "with"
Cohesion: 0.06
Nodes (38): needs, runs-on, steps, needs, runs-on, steps, env, GO_VERSION (+30 more)

### Community 7 - "testing.T"
Cohesion: 0.08
Nodes (30): TestCollectBody(), TestCollectBodyString(), TestCollectHeaders(), TestBuildFilePathForm(), TestBuildFilePathFormValidationLogic(), CollectBodyStringFromForm(), TestCollectBodyFromForm(), TestCollectBodyStringFromForm() (+22 more)

### Community 8 - "steps"
Cohesion: 0.08
Nodes (28): needs, runs-on, steps, CHOCOLATEY_API_KEY, GITHUB_TOKEN, HOMEBREW_TAP_GITHUB_TOKEN, jobs, chocolatey (+20 more)

### Community 9 - "http.go"
Cohesion: 0.15
Nodes (7): Client, isIdempotent(), isRetryableError(), IsUnforwardableHeader(), mergeHeaders(), Response, RetryConfig

### Community 10 - "steps"
Cohesion: 0.08
Nodes (23): jobs, update-tap, name, on, workflow_run, permissions, contents, id (+15 more)

### Community 11 - "http_test.go"
Cohesion: 0.16
Nodes (21): DefaultRetryConfig(), NewHttpClient(), TestClient_buildRequest(), TestClient_calculateBackoff(), TestClient_Delete(), TestClient_executeRequest(), TestClient_Get(), TestClient_handleResponse() (+13 more)

### Community 12 - "requestTo"
Cohesion: 0.09
Nodes (22): testKey, circuitBreaker, enabled, failureThreshold, halfOpenMaxRequests, minimumRequests, openTimeoutMs, successThreshold (+14 more)

### Community 13 - ".golangci.yml"
Cohesion: 0.09
Nodes (22): max-complexity, generated, paths, presets, formatters, enable, exclusions, lines (+14 more)

### Community 14 - ".CreateHandler"
Cohesion: 0.12
Nodes (9): TestEndpointHandler_BodylessResponse(), TestMockRoute_CreateRoute(), EndpointHandler, Handler, MainRouter, request(), TestRouter_CreateRoutes(), TestRouter_CreateRoutes_BothHandlersSet() (+1 more)

### Community 15 - "Inzibat — Agent Guide"
Cohesion: 0.08
Nodes (21): Always Do, Architecture, Build & Test, CLI, CLI, Config, Entrypoints, GitNexus — Code Intelligence (+13 more)

### Community 16 - "Changelog"
Cohesion: 0.12
Nodes (17): [0.1.0] - 2025-11-11, [0.1.1] - 2025-11-14, [0.1.3] - 2025-11-18, [0.1.7] - 2025-12-15, [0.3.1] - 2026-06-07, [0.3.2] - 2026-06-07, [0.3.7] - 2026-06-16, [0.3.9] - 2026-06-16 (+9 more)

### Community 17 - "form_builder_test.go"
Cohesion: 0.06
Nodes (58): createClientRequestForm(), createMockResponseForm(), createRoute(), createRouteForm(), defaultCircuitBreakerFormValues(), TestCreateClientRequestFormHeadless(), TestCreateClientRequestFormHeadlessBreakerDisabled(), TestCreateMockResponseFormHeadlessBodyObject() (+50 more)

### Community 18 - "setupServer"
Cohesion: 0.13
Nodes (16): TestStartServerCmd_ConfigFile(), TestHealthCheckRouteRoundTrip(), TestWriteConfig(), WriteConfig(), TestCfg_GetServerAddr(), ResolveAbsolutePath(), TestResolveAbsolutePath(), runServer() (+8 more)

### Community 19 - "RecordedEntry"
Cohesion: 0.19
Nodes (13): exportAsInzibatConfig(), TestExportAsInzibatConfig(), buildFakeResponse(), convertHeaders(), ConvertToInzibatConfig(), TestBuildFakeResponse(), TestConvertHeaders(), TestConvertToInzibatConfig() (+5 more)

### Community 20 - "Contributor Covenant Code of Conduct"
Cohesion: 0.15
Nodes (12): 1. Correction, 2. Warning, 3. Temporary Ban, 4. Permanent Ban, Attribution, Contributor Covenant Code of Conduct, Enforcement, Enforcement Guidelines (+4 more)

### Community 21 - "readFile"
Cohesion: 0.20
Nodes (12): newFailOpeningError(), newFailReadingError(), TestNewFailOpeningError(), TestNewFailReadingError(), LoadBodyFromFile(), LoadBodyStringFromFile(), LoadHeadersFromFile(), readFile() (+4 more)

### Community 22 - "hooks"
Cohesion: 0.17
Nodes (11): default_install_hook_types, always_run, entry, id, language, name, pass_filenames, stages (+3 more)

### Community 23 - "Store"
Cohesion: 0.27
Nodes (5): clearEntriesHandler(), getSessionHandler(), listEntriesHandler(), RegisterAdminRoutes(), Store

### Community 24 - "README.md"
Cohesion: 0.22
Nodes (5): Contributing, Development setup, Pull requests, Reporting issues, Roadmap

### Community 25 - "record_test.go"
Cohesion: 0.18
Nodes (11): exportAsJSON(), fetchRecordedEntries(), fetchRecordedSession(), getJSON(), TestExportAsJSON(), TestFetchRecordedEntries(), TestFetchRecordedSession(), TestGetJSON_RespectsTimeout() (+3 more)

### Community 26 - "Inzibat 🪖"
Cohesion: 0.29
Nodes (7): Configuration, Installation, Inzibat 🪖, License, Project status, Usage, Why Inzibat

### Community 27 - "GitNexus — Code Intelligence"
Cohesion: 0.33
Nodes (5): Always Do, CLI, GitNexus — Code Intelligence, Never Do, Resources

### Community 28 - "TestSchemaMatchesModel"
Cohesion: 0.50
Nodes (5): joinPath(), modelPaths(), resolveSchemaRefs(), schemaPaths(), TestSchemaMatchesModel()

### Community 29 - "PULL_REQUEST_TEMPLATE.md"
Cohesion: 0.40
Nodes (4): Checklist, Description, Type of Change, Verification & Testing

### Community 30 - "Security Policy"
Cohesion: 0.50
Nodes (4): Our Commitment, Reporting a Vulnerability, Security Policy, Supported Versions

### Community 31 - "[0.4.0] - 2026-06-19"
Cohesion: 0.67
Nodes (3): [0.4.0] - 2026-06-19, Added, Fixed

### Community 49 - "actions"
Cohesion: 0.18
Nodes (10): actions, actions/checkout, actions/download-artifact, actions/setup-go, actions/upload-artifact, golangci/golangci-lint-action, goreleaser/goreleaser-action, sonarsource/sonarqube-quality-gate-action (+2 more)

### Community 50 - "BoolPointer"
Cohesion: 0.29
Nodes (11): BoolPointer(), NewCircuitBreakerStore(), TestCircuitBreaker(), TestCircuitBreaker_HalfOpenGaugeReleased(), TestCircuitBreakerStore_Get(), TestCircuitBreakerStore_OnFailure(), TestCircuitBreakerStore_OnSuccess(), TestCircuitBreakerStore_Seed() (+3 more)

### Community 51 - ".CreateHandler"
Cohesion: 0.27
Nodes (5): TestClientHandler_allowRequest(), TestClientHandler_CreateHandler_NoNullBody(), TestClientHandler_CreateHandler_PassThrough(), TestClientHandler_recordFailure(), ClientHandler

### Community 52 - "store_test.go"
Cohesion: 0.50
Nodes (8): NewStore(), newTestEntry(), TestNewStore(), TestStoreAdd(), TestStoreClear(), TestStoreConcurrency(), TestStoreList(), TestStoreSession()

### Community 53 - "setupAdminApp"
Cohesion: 0.40
Nodes (4): setupAdminApp(), TestClearEntriesHandler(), TestGetSessionHandler(), TestListEntriesHandler()

## Knowledge Gaps
- **1 isolated node(s):** `github.com/lynicis/inzibat`
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 240 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **18 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What connects `github.com/lynicis/inzibat` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Cfg` be split into smaller, more focused modules?**
  _Cohesion score 0.06382978723404255 - nodes in this community are weakly interconnected._
- **Should `FormRunner` be split into smaller, more focused modules?**
  _Cohesion score 0.08408163265306122 - nodes in this community are weakly interconnected._
- **Should `nfpms` be split into smaller, more focused modules?**
  _Cohesion score 0.041666666666666664 - nodes in this community are weakly interconnected._
- **Should `with` be split into smaller, more focused modules?**
  _Cohesion score 0.06477732793522267 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.07862903225806452 - nodes in this community are weakly interconnected._
- **Should `steps` be split into smaller, more focused modules?**
  _Cohesion score 0.07881773399014778 - nodes in this community are weakly interconnected._