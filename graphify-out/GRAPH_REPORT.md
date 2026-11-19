# Graph Report - inzibat  (2026-11-19)

## Corpus Check
- 93 files · ~37,862 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 7 file(s) not represented in the graph (top: (none) 3, .toml 3, .properties 1)

## Summary
- 840 nodes · 2234 edges · 47 communities (38 shown, 9 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 293 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `ada8ebc1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Config and Handler Core
- Store
- Route Model and Forms
- GoReleaser Release Packaging
- CI Workflow Pipeline
- Release and Distribution Workflow
- BoolPointer
- Circuit Breaker State Machine
- Project Docs and Templates
- circuitBreaker
- Homebrew Tap Update Workflow
- Circuit Breaker Config Schema
- Lint Configuration
- HTTP Client Tests
- Headless Form Runner
- form_builder_test.go
- Router and Endpoint Handlers
- Cobra CLI Entrypoint
- quickstart.json
- setupServer
- Cfg
- steps
- Body and Headers File Loading
- Pre-commit Hooks
- GitHub Action Pins
- go_pkg_time
- compileConfigSchema
- server.go
- record.go
- record_test.go
- Proxying and Mocking Concepts
- NewLoader
- CircuitBreakerConfig
- store_test.go
- Bug Reporting Template
- Recorder Release Notes
- admin_test.go
- Contribution and Conduct
- PR Checklist and Setup
- Mock Regeneration Rule
- Agent Guide Docs
- Architecture Data Flow
- Interactive Route Builder CLI
- Config Precedence Rule
- Roadmap RPC and Registry
- Security Policy
- Go Module Identity

## God Nodes (most connected - your core abstractions)
1. `Store` - 22 edges
2. `NewHttpClient()` - 22 edges
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
- `createMockResponseForm()` --references--> `FakeResponse`  [EXTRACTED]
  cmd/create.go → config/model.go
- `buildFakeResponse()` --references--> `FakeResponse`  [EXTRACTED]
  recorder/converter.go → config/model.go
- `CollectBody()` --references--> `HttpBody`  [EXTRACTED]
  cmd/formbuilder/data_collector.go → config/model.go
- `CollectBodyFromForm()` --references--> `HttpBody`  [EXTRACTED]
  cmd/formbuilder/form_collector.go → config/model.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Agent Guide Architecture Contract** — agents_entrypoints, agents_architecture, agents_config_precedence, agents_cli_commands, agents_circuit_breaker [EXTRACTED 1.00]
- **Inzibat Core Feature Set** — readme_circuit_breaker_proxying, readme_traffic_recording, readme_interactive_route_builder, readme_schema_validated_config [EXTRACTED 1.00]
- **Contribution Workflow Documents** — contributing_pull_request_process, contributing_issue_reporting, _github_pull_request_template_pull_request_template, _github_issue_template_bug_report_bug_report_template, _github_issue_template_feature_request_feature_request_template [INFERRED 0.85]
- **Repository meta documents** — readme_inzibatoverview, contributing_contributionguide, changelog_releasehistory, security_securitypolicy, code_of_conduct_codeofconduct [INFERRED 0.85]

## Communities (47 total, 9 thin omitted)

### Community 0 - "Config and Handler Core"
Cohesion: 0.16
Nodes (13): TestValidateHost(), TestValidatePath(), TestValidatePositiveInt(), TestValidateStatusCode(), ValidateHost(), ValidatePath(), ValidatePositiveInt(), ValidateStatusCode() (+5 more)

### Community 1 - "Store"
Cohesion: 0.25
Nodes (5): clearEntriesHandler(), getSessionHandler(), listEntriesHandler(), RegisterAdminRoutes(), Store

### Community 2 - "Route Model and Forms"
Cohesion: 0.07
Nodes (39): collectCircuitBreakerConfig(), collectClientRequestBody(), createClientRequestFormInternal(), createMockResponseFormInternal(), createRouteInternal(), parsePositiveInt(), TestCollectCircuitBreakerConfig(), TestCreateClientRequestFormInternal() (+31 more)

### Community 3 - "GoReleaser Release Packaging"
Cohesion: 0.04
Nodes (47): archives, format_overrides, formats, name_template, before, hooks, builds, env (+39 more)

### Community 4 - "CI Workflow Pipeline"
Cohesion: 0.06
Nodes (38): needs, runs-on, steps, needs, runs-on, steps, env, GO_VERSION (+30 more)

### Community 5 - "Release and Distribution Workflow"
Cohesion: 0.11
Nodes (18): GITHUB_TOKEN, HOMEBREW_TAP_GITHUB_TOKEN, jobs, release, name, on, push, contents (+10 more)

### Community 6 - "BoolPointer"
Cohesion: 0.15
Nodes (22): TestIsRetryableError(), TestIsUnforwardableHeader(), TestCollectBody(), TestCollectBodyString(), TestCollectHeaders(), TestFormPlaceholders(), TestListCmd(), TestListCmd_Run() (+14 more)

### Community 7 - "Circuit Breaker State Machine"
Cohesion: 0.07
Nodes (19): Client, isIdempotent(), isRetryableError(), IsUnforwardableHeader(), allowHalfOpen(), CircuitBreakerStore, openRecord(), releaseHalfOpen() (+11 more)

### Community 8 - "Project Docs and Templates"
Cohesion: 0.08
Nodes (23): Bug Report Template, Feature Request Issue Template, Feature Request Template, Pull Request Template, Graphify Code Intelligence Index, Graphify Index Replacing GitNexus, JSON Schema for Configuration, Release history (Keep a Changelog) (+15 more)

### Community 9 - "circuitBreaker"
Cohesion: 0.18
Nodes (10): circuitBreaker, enabled, failureThreshold, halfOpenMaxRequests, minimumRequests, openTimeoutMs, successThreshold, routes (+2 more)

### Community 10 - "Homebrew Tap Update Workflow"
Cohesion: 0.06
Nodes (33): EVENT_NAME, INPUT_TAG, WORKFLOW_SHA, tag, jobs, update-tap, name, on (+25 more)

### Community 11 - "Circuit Breaker Config Schema"
Cohesion: 0.09
Nodes (22): testKey, xTestHeader, body, circuitBreaker, headers, host, method, passWithRequestBody (+14 more)

### Community 12 - "Lint Configuration"
Cohesion: 0.09
Nodes (22): max-complexity, generated, paths, presets, formatters, enable, exclusions, lines (+14 more)

### Community 13 - "HTTP Client Tests"
Cohesion: 0.15
Nodes (22): DefaultRetryConfig(), NewHttpClient(), TestClient_buildRequest(), TestClient_calculateBackoff(), TestClient_Delete(), TestClient_executeRequest(), TestClient_Get(), TestClient_handleResponse() (+14 more)

### Community 14 - "Headless Form Runner"
Cohesion: 0.14
Nodes (21): createClientRequestForm(), createMockResponseForm(), createRoute(), defaultCircuitBreakerFormValues(), breakerDefaults(), clientRequestSegments(), TestCreateClientRequestFormHeadless(), TestCreateClientRequestFormHeadlessBreakerDisabled() (+13 more)

### Community 15 - "form_builder_test.go"
Cohesion: 0.07
Nodes (49): createRouteForm(), CollectBody(), CollectBodyString(), CollectHeaders(), BuildFilePathForm(), BuildSourceSelectionForm(), NewHuhForm(), TestBuildFilePathForm() (+41 more)

### Community 16 - "Router and Endpoint Handlers"
Cohesion: 0.12
Nodes (9): TestEndpointHandler_BodylessResponse(), TestMockRoute_CreateRoute(), EndpointHandler, Handler, MainRouter, request(), TestRouter_CreateRoutes(), TestRouter_CreateRoutes_BothHandlersSet() (+1 more)

### Community 17 - "Cobra CLI Entrypoint"
Cohesion: 0.18
Nodes (8): Execute(), TestExecute(), TestExecute_ErrorPath(), TestRootCmd_RunsServer(), addServerFlags(), init(), runStartServer(), main()

### Community 18 - "quickstart.json"
Cohesion: 0.50
Nodes (3): routes, $schema, serverPort

### Community 19 - "setupServer"
Cohesion: 0.13
Nodes (17): TestStartServerCmd_ConfigFile(), TestHealthCheckRouteRoundTrip(), TestWriteConfig(), WriteConfig(), TestCfg_GetServerAddr(), ResolveAbsolutePath(), TestResolveAbsolutePath(), loadConfig() (+9 more)

### Community 20 - "Cfg"
Cohesion: 0.14
Nodes (14): routeRows(), TestRouteRows(), addHealthCheckRoute(), DefaultCfg(), normalizeRoutes(), ReadOrCreateConfig(), TestReadOrCreateConfig(), koanfReader (+6 more)

### Community 21 - "steps"
Cohesion: 0.06
Nodes (33): if, runs-on, steps, CHOCOLATEY_API_KEY, EVENT_NAME, INPUT_TAG, TAG, VERSION (+25 more)

### Community 22 - "Body and Headers File Loading"
Cohesion: 0.20
Nodes (12): newFailOpeningError(), newFailReadingError(), TestNewFailOpeningError(), TestNewFailReadingError(), LoadBodyFromFile(), LoadBodyStringFromFile(), LoadHeadersFromFile(), readFile() (+4 more)

### Community 23 - "Pre-commit Hooks"
Cohesion: 0.17
Nodes (11): default_install_hook_types, always_run, entry, id, language, name, pass_filenames, stages (+3 more)

### Community 24 - "GitHub Action Pins"
Cohesion: 0.18
Nodes (10): actions, actions/checkout, actions/download-artifact, actions/setup-go, actions/upload-artifact, golangci/golangci-lint-action, goreleaser/goreleaser-action, sonarsource/sonarqube-quality-gate-action (+2 more)

### Community 25 - "go_pkg_time"
Cohesion: 0.15
Nodes (15): exportAsInzibatConfig(), exportAsJSON(), TestExportAsInzibatConfig(), TestExportAsJSON(), buildFakeResponse(), convertHeaders(), ConvertToInzibatConfig(), TestBuildFakeResponse() (+7 more)

### Community 26 - "compileConfigSchema"
Cohesion: 0.60
Nodes (4): compileConfigSchema(), TestSchemaAcceptsExamples(), TestSchemaRejectsInvalidConfigs(), validateConfigDocument()

### Community 27 - "server.go"
Cohesion: 0.17
Nodes (8): captureBody(), captureHeaderValues(), captureRequestHeaders(), captureResponseHeaders(), NewRecorderMiddleware(), BenchmarkMiddleware(), TestCaptureBody(), TestNewRecorderMiddleware()

### Community 29 - "record_test.go"
Cohesion: 0.19
Nodes (10): adminURL(), fetchRecordedEntries(), fetchRecordedSession(), getJSON(), TestFetchRecordedEntries(), TestFetchRecordedSession(), TestGetJSON_RespectsTimeout(), TestRecordClearCmd_ReturnsConnectionError() (+2 more)

### Community 30 - "Proxying and Mocking Concepts"
Cohesion: 0.33
Nodes (7): Circuit Breaker (proxy routes), fakeResponse Mock Handler, requestTo Proxy Handler, Circuit Breaker (CLAUDE guide), Proxying with Circuit Breaker, requestTo Overrides fakeResponse Rule, Circuit Breaker for Clients (shipped)

### Community 31 - "NewLoader"
Cohesion: 0.24
Nodes (9): Reader, NewLoader(), ResolveConfigPath(), resolveConfigPath(), TestNewLoader(), TestReader_Read_AllowsBodylessFakeResponse(), TestReader_Read_AppliesDefaultsBeforeValidation(), TestReader_Read_HealthRouteWithEmptyRoutes() (+1 more)

### Community 32 - "CircuitBreakerConfig"
Cohesion: 0.27
Nodes (9): applyCircuitBreakerConfig(), DefaultCircuitBreakerConfig(), CircuitBreakerConfig, MergeCircuitBreakerConfig(), TestMergeCircuitBreakerConfig(), BuildCircuitBreakerRouteKey(), TestBuildCircuitBreakerRouteKey(), seedCircuitBreakers() (+1 more)

### Community 33 - "store_test.go"
Cohesion: 0.50
Nodes (8): NewStore(), newTestEntry(), TestNewStore(), TestStoreAdd(), TestStoreClear(), TestStoreConcurrency(), TestStoreList(), TestStoreSession()

### Community 34 - "Bug Reporting Template"
Cohesion: 0.50
Nodes (4): Bug Report Issue Template, Bug Reproduction Steps, Issue Reporting Guidance, Private Vulnerability Disclosure

### Community 35 - "Recorder Release Notes"
Cohesion: 0.50
Nodes (4): Request Recorder Feature, Release 0.4.0, Recorder Admin API, Traffic Recording and Replay

### Community 36 - "admin_test.go"
Cohesion: 0.43
Nodes (4): setupAdminApp(), TestClearEntriesHandler(), TestGetSessionHandler(), TestListEntriesHandler()

### Community 37 - "Contribution and Conduct"
Cohesion: 0.67
Nodes (3): Pull Request Checklist, Contributor Covenant Code of Conduct, Contribution guide

### Community 38 - "PR Checklist and Setup"
Cohesion: 0.67
Nodes (3): PR Verification Checklist, Build and Test Commands, Development Setup

## Knowledge Gaps
- **37 isolated node(s):** `github.com/lynicis/inzibat`, `fakeResponse Mock Handler`, `Circuit Breaker (CLAUDE guide)`, `Bug Reproduction Steps`, `Private Vulnerability Disclosure` (+32 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 243 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Store` to `store_test.go`, `admin_test.go`, `Circuit Breaker State Machine`, `go_pkg_time`, `server.go`, `record.go`, `record_test.go`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `github.com/lynicis/inzibat`, `fakeResponse Mock Handler`, `Circuit Breaker (CLAUDE guide)` to the rest of the system?**
  _37 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Route Model and Forms` be split into smaller, more focused modules?**
  _Cohesion score 0.07337662337662337 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `Circuit Breaker State Machine` to `HTTP Client Tests`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `GoReleaser Release Packaging` be split into smaller, more focused modules?**
  _Cohesion score 0.041666666666666664 - nodes in this community are weakly interconnected._
- **Why does `Cfg` connect `Cfg` to `Config and Handler Core`, `CircuitBreakerConfig`, `Route Model and Forms`, `Router and Endpoint Handlers`, `setupServer`, `go_pkg_time`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Should `CI Workflow Pipeline` be split into smaller, more focused modules?**
  _Cohesion score 0.06477732793522267 - nodes in this community are weakly interconnected._