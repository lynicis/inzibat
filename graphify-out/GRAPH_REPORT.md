# Graph Report - inzibat  (2026-11-16)

## Corpus Check
- Corpus is ~37,187 words - fits in a single context window. You may not need a graph.

## Summary
- 744 nodes · 2144 edges · 28 communities (24 shown, 4 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 272 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Interactive Form Internals
- Recorder CLI
- Formbuilder Collectors
- GoReleaser Config
- Config Loading
- CI Workflow
- Release Workflow
- Formbuilder Tests
- HTTP Client Tests
- Circuit Breaker Store
- HTTP Client Core
- Homebrew Tap Workflow
- Example Configs
- Lint Config
- Endpoint Handlers
- Repo Meta Docs
- Record Export Tests
- Root Command
- Server Bootstrap
- Circuit Breaker Tests
- Pre-commit Hooks
- Actup Config
- Validators Tests
- JSON Schema
- Roadmap Items
- Security Policy
- Module Root

## God Nodes (most connected - your core abstractions)
1. `NewHttpClient()` - 22 edges
2. `Store` - 22 edges
3. `Cfg` - 21 edges
4. `FormRunner` - 17 edges
5. `BoolPointer()` - 16 edges
6. `NewCircuitBreakerStore()` - 15 edges
7. `CircuitBreakerStore` - 14 edges
8. `GetFreePort()` - 14 edges
9. `setupServer()` - 14 edges
10. `Client` - 13 edges

## Surprising Connections (you probably didn't know these)
- `TestClientHandler_CreateHandler_NoNullBody()` --calls--> `captureBody()`  [INFERRED]
  handler/client_test.go → recorder/middleware.go
- `ClientHandler` --references--> `Client`  [EXTRACTED]
  handler/client.go → client/http/http.go
- `TestClientHandler_CreateHandler()` --calls--> `NewHttpClient()`  [EXTRACTED]
  handler/client_test.go → client/http/http.go
- `TestClientHandler_CreateHandler_NoNullBody()` --calls--> `NewHttpClient()`  [EXTRACTED]
  handler/client_test.go → client/http/http.go
- `TestClientHandler_CreateHandler_PassThrough()` --calls--> `NewHttpClient()`  [EXTRACTED]
  handler/client_test.go → client/http/http.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Repository meta documents** — agents_agentguide, readme_inzibatoverview, contributing_contributionguide, changelog_releasehistory, security_securitypolicy, code_of_conduct_codeofconduct [INFERRED 0.85]

## Communities (28 total, 4 thin omitted)

### Community 1 - "Interactive Form Internals"
Cohesion: 0.05
Nodes (56): collectCircuitBreakerConfig(), collectClientRequestBody(), createClientRequestForm(), createClientRequestFormInternal(), createMockResponseForm(), createMockResponseFormInternal(), createRoute(), createRouteInternal() (+48 more)

### Community 2 - "Recorder CLI"
Cohesion: 0.06
Nodes (38): fetchRecordedEntries(), fetchRecordedSession(), getJSON(), TestFetchRecordedEntries(), TestFetchRecordedSession(), TestGetJSON_RespectsTimeout(), TestRecordClearCmd_ReturnsConnectionError(), TestRecordExportCmd_RejectsUnknownFormat() (+30 more)

### Community 3 - "Formbuilder Collectors"
Cohesion: 0.07
Nodes (48): createRouteForm(), CollectBody(), CollectBodyString(), CollectHeaders(), BuildFilePathForm(), BuildSourceSelectionForm(), NewHuhForm(), TestBuildFilePathForm() (+40 more)

### Community 4 - "GoReleaser Config"
Cohesion: 0.04
Nodes (47): archives, format_overrides, formats, name_template, before, hooks, builds, env (+39 more)

### Community 5 - "Config Loading"
Cohesion: 0.07
Nodes (36): routeRows(), TestRouteRows(), addHealthCheckRoute(), DefaultCfg(), Reader, NewLoader(), normalizeRoutes(), ReadOrCreateConfig() (+28 more)

### Community 6 - "CI Workflow"
Cohesion: 0.06
Nodes (38): needs, runs-on, steps, needs, runs-on, steps, env, GO_VERSION (+30 more)

### Community 7 - "Release Workflow"
Cohesion: 0.08
Nodes (28): needs, runs-on, steps, CHOCOLATEY_API_KEY, GITHUB_TOKEN, HOMEBREW_TAP_GITHUB_TOKEN, jobs, chocolatey (+20 more)

### Community 8 - "Formbuilder Tests"
Cohesion: 0.09
Nodes (26): TestIsRetryableError(), TestIsUnforwardableHeader(), TestCollectBody(), TestCollectBodyString(), TestCollectHeaders(), TestFormPlaceholders(), TestListCmd(), TestListCmd_Run() (+18 more)

### Community 9 - "HTTP Client Tests"
Cohesion: 0.17
Nodes (20): DefaultRetryConfig(), NewHttpClient(), TestClient_buildRequest(), TestClient_calculateBackoff(), TestClient_Delete(), TestClient_executeRequest(), TestClient_Get(), TestClient_handleResponse() (+12 more)

### Community 10 - "Circuit Breaker Store"
Cohesion: 0.15
Nodes (12): allowHalfOpen(), CircuitBreakerStore, openRecord(), releaseHalfOpen(), resetToClosed(), TestAllowHalfOpen(), CircuitBreakerRecord, CircuitBreakerState (+4 more)

### Community 11 - "HTTP Client Core"
Cohesion: 0.16
Nodes (7): Client, isIdempotent(), isRetryableError(), IsUnforwardableHeader(), mergeHeaders(), Response, RetryConfig

### Community 12 - "Homebrew Tap Workflow"
Cohesion: 0.08
Nodes (23): jobs, update-tap, name, on, workflow_run, permissions, contents, id (+15 more)

### Community 13 - "Example Configs"
Cohesion: 0.09
Nodes (22): testKey, circuitBreaker, enabled, failureThreshold, halfOpenMaxRequests, minimumRequests, openTimeoutMs, successThreshold (+14 more)

### Community 14 - "Lint Config"
Cohesion: 0.09
Nodes (22): max-complexity, generated, paths, presets, formatters, enable, exclusions, lines (+14 more)

### Community 15 - "Endpoint Handlers"
Cohesion: 0.12
Nodes (9): TestEndpointHandler_BodylessResponse(), TestMockRoute_CreateRoute(), EndpointHandler, Handler, MainRouter, request(), TestRouter_CreateRoutes(), TestRouter_CreateRoutes_BothHandlersSet() (+1 more)

### Community 16 - "Repo Meta Docs"
Cohesion: 0.14
Nodes (14): Bug Report Template, Feature Request Template, Pull Request Checklist, Agent Guide (repo conventions), Graphify code intelligence index, Release history (Keep a Changelog), Release 0.5.0 changes, CLAUDE.md (symlink to AGENTS.md) (+6 more)

### Community 17 - "Record Export Tests"
Cohesion: 0.14
Nodes (16): exportAsJSON(), TestExportAsJSON(), newFailOpeningError(), newFailReadingError(), TestNewFailOpeningError(), TestNewFailReadingError(), LoadBodyFromFile(), LoadBodyStringFromFile() (+8 more)

### Community 18 - "Root Command"
Cohesion: 0.19
Nodes (8): Execute(), TestExecute(), TestExecute_ErrorPath(), TestRootCmd_RunsServer(), addServerFlags(), init(), runStartServer(), main()

### Community 19 - "Server Bootstrap"
Cohesion: 0.21
Nodes (10): TestCfg_GetServerAddr(), runServer(), setupServer(), StartServer(), StartServerWithContext(), contains(), TestRunServer_ReturnsListenError(), TestSetupServer_CircuitBreaker() (+2 more)

### Community 20 - "Circuit Breaker Tests"
Cohesion: 0.26
Nodes (13): BoolPointer(), NewCircuitBreakerStore(), TestCircuitBreaker(), TestCircuitBreaker_HalfOpenGaugeReleased(), TestCircuitBreakerStore_Allow(), TestCircuitBreakerStore_Get(), TestCircuitBreakerStore_OnFailure(), TestCircuitBreakerStore_OnSuccess() (+5 more)

### Community 21 - "Pre-commit Hooks"
Cohesion: 0.17
Nodes (11): default_install_hook_types, always_run, entry, id, language, name, pass_filenames, stages (+3 more)

### Community 22 - "Actup Config"
Cohesion: 0.18
Nodes (10): actions, actions/checkout, actions/download-artifact, actions/setup-go, actions/upload-artifact, golangci/golangci-lint-action, goreleaser/goreleaser-action, sonarsource/sonarqube-quality-gate-action (+2 more)

### Community 23 - "Validators Tests"
Cohesion: 0.20
Nodes (9): TestValidateHost(), TestValidateNonEmpty(), TestValidatePath(), TestValidatePositiveInt(), TestValidateStatusCode(), ValidateHost(), ValidatePath(), ValidatePositiveInt() (+1 more)

### Community 24 - "JSON Schema"
Cohesion: 0.50
Nodes (5): joinPath(), modelPaths(), resolveSchemaRefs(), schemaPaths(), TestSchemaMatchesModel()

## Knowledge Gaps
- **9 isolated node(s):** `github.com/lynicis/inzibat`, `Bug Report Template`, `Feature Request Template`, `Contributor Covenant Code of Conduct`, `Roadmap: service registry with KV sidecar` (+4 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 179 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Recorder CLI` to `Route Creation Forms`, `Circuit Breaker Store`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **What connects `github.com/lynicis/inzibat`, `Bug Report Template`, `Feature Request Template` to the rest of the system?**
  _9 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Interactive Form Internals` be split into smaller, more focused modules?**
  _Cohesion score 0.0547945205479452 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `HTTP Client Core` to `HTTP Client Tests`, `Circuit Breaker Store`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Should `Recorder CLI` be split into smaller, more focused modules?**
  _Cohesion score 0.061581920903954805 - nodes in this community are weakly interconnected._
- **Why does `Cfg` connect `Config Loading` to `Route Creation Forms`, `Interactive Form Internals`, `Formbuilder Tests`, `Endpoint Handlers`, `Server Bootstrap`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Should `Formbuilder Collectors` be split into smaller, more focused modules?**
  _Cohesion score 0.06966618287373004 - nodes in this community are weakly interconnected._