# Graph Report - inzibat  (2026-11-17)

## Corpus Check
- Corpus is ~37,245 words - fits in a single context window. You may not need a graph.

## Summary
- 785 nodes · 2177 edges · 47 communities (38 shown, 9 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 293 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Request Recording Store
- Route Model and Forms
- GoReleaser Release Packaging
- CI Workflow Pipeline
- Release and Distribution Workflow
- CLI Command Tests
- Circuit Breaker State Machine
- Project Docs and Templates
- HTTP Client Retry Logic
- Homebrew Tap Update Workflow
- Circuit Breaker Config Schema
- Lint Configuration
- HTTP Client Tests
- Headless Form Runner
- File Path Source Forms
- Router and Endpoint Handlers
- Cobra CLI Entrypoint
- Route Creation Forms
- Server Bootstrap
- Koanf Config Reader
- Reader Strategy Loader
- Body and Headers File Loading
- Pre-commit Hooks
- GitHub Action Pins
- Headless Terminal Collectors
- Input Validators
- Form Collector Tests
- Circuit Breaker Config Merge
- Config Write and Export
- Proxying and Mocking Concepts
- Reader Strategy Mocks
- Config Path Resolution
- JSON Schema Sync
- Bug Reporting Template
- Recorder Release Notes
- Response Body Capture
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
- `ClientHandler` --references--> `Client`  [EXTRACTED]
  handler/client.go → client/http/http.go
- `TestClientHandler_CreateHandler()` --calls--> `NewHttpClient()`  [EXTRACTED]
  handler/client_test.go → client/http/http.go
- `TestClientHandler_CreateHandler_NoNullBody()` --calls--> `NewHttpClient()`  [EXTRACTED]
  handler/client_test.go → client/http/http.go
- `TestClientHandler_CreateHandler_PassThrough()` --calls--> `NewHttpClient()`  [EXTRACTED]
  handler/client_test.go → client/http/http.go
- `setupServer()` --calls--> `NewHttpClient()`  [EXTRACTED]
  server/server.go → client/http/http.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Repository meta documents** — readme_inzibatoverview, contributing_contributionguide, changelog_releasehistory, security_securitypolicy, code_of_conduct_codeofconduct [INFERRED 0.85]
- **Inzibat Core Feature Set** — readme_circuit_breaker_proxying, readme_traffic_recording, readme_interactive_route_builder, readme_schema_validated_config [EXTRACTED 1.00]
- **Contribution Workflow Documents** — contributing_pull_request_process, contributing_issue_reporting, _github_pull_request_template_pull_request_template, _github_issue_template_bug_report_bug_report_template, _github_issue_template_feature_request_feature_request_template [INFERRED 0.85]
- **Agent Guide Architecture Contract** — agents_entrypoints, agents_architecture, agents_config_precedence, agents_cli_commands, agents_circuit_breaker [EXTRACTED 1.00]

## Communities (47 total, 9 thin omitted)

### Community 1 - "Request Recording Store"
Cohesion: 0.06
Nodes (43): adminURL(), exportAsInzibatConfig(), fetchRecordedEntries(), fetchRecordedSession(), getJSON(), TestExportAsInzibatConfig(), TestFetchRecordedEntries(), TestFetchRecordedSession() (+35 more)

### Community 2 - "Route Model and Forms"
Cohesion: 0.07
Nodes (39): collectCircuitBreakerConfig(), collectClientRequestBody(), createClientRequestForm(), createClientRequestFormInternal(), createMockResponseFormInternal(), createRoute(), createRouteInternal(), defaultCircuitBreakerFormValues() (+31 more)

### Community 3 - "GoReleaser Release Packaging"
Cohesion: 0.04
Nodes (47): archives, format_overrides, formats, name_template, before, hooks, builds, env (+39 more)

### Community 4 - "CI Workflow Pipeline"
Cohesion: 0.06
Nodes (38): needs, runs-on, steps, needs, runs-on, steps, env, GO_VERSION (+30 more)

### Community 5 - "Release and Distribution Workflow"
Cohesion: 0.08
Nodes (28): needs, runs-on, steps, CHOCOLATEY_API_KEY, GITHUB_TOKEN, HOMEBREW_TAP_GITHUB_TOKEN, jobs, chocolatey (+20 more)

### Community 6 - "CLI Command Tests"
Cohesion: 0.12
Nodes (27): TestDefaultCircuitBreakerFormValues(), TestCollectBody(), TestCollectBodyString(), TestCollectHeaders(), TestFormPlaceholders(), TestListCmd(), TestListCmd_Run(), TestRecordClearCmd_ReturnsConnectionError() (+19 more)

### Community 7 - "Circuit Breaker State Machine"
Cohesion: 0.14
Nodes (13): allowHalfOpen(), CircuitBreakerStore, openRecord(), releaseHalfOpen(), resetToClosed(), TestAllowHalfOpen(), CircuitBreakerRecord, CircuitBreakerState (+5 more)

### Community 8 - "Project Docs and Templates"
Cohesion: 0.08
Nodes (23): Bug Report Template, Feature Request Issue Template, Feature Request Template, Pull Request Template, Graphify Code Intelligence Index, Graphify Index Replacing GitNexus, JSON Schema for Configuration, Release history (Keep a Changelog) (+15 more)

### Community 9 - "HTTP Client Retry Logic"
Cohesion: 0.14
Nodes (9): Client, isIdempotent(), isRetryableError(), IsUnforwardableHeader(), TestIsRetryableError(), TestIsUnforwardableHeader(), mergeHeaders(), Response (+1 more)

### Community 10 - "Homebrew Tap Update Workflow"
Cohesion: 0.08
Nodes (23): jobs, update-tap, name, on, workflow_run, permissions, contents, id (+15 more)

### Community 11 - "Circuit Breaker Config Schema"
Cohesion: 0.09
Nodes (22): testKey, circuitBreaker, enabled, failureThreshold, halfOpenMaxRequests, minimumRequests, openTimeoutMs, successThreshold (+14 more)

### Community 12 - "Lint Configuration"
Cohesion: 0.09
Nodes (22): max-complexity, generated, paths, presets, formatters, enable, exclusions, lines (+14 more)

### Community 13 - "HTTP Client Tests"
Cohesion: 0.20
Nodes (19): DefaultRetryConfig(), NewHttpClient(), TestClient_buildRequest(), TestClient_calculateBackoff(), TestClient_Delete(), TestClient_executeRequest(), TestClient_Get(), TestClient_handleResponse() (+11 more)

### Community 14 - "Headless Form Runner"
Cohesion: 0.17
Nodes (17): createMockResponseForm(), breakerDefaults(), clientRequestSegments(), TestCreateClientRequestFormHeadless(), TestCreateClientRequestFormHeadlessBreakerDisabled(), TestCreateMockResponseFormHeadlessBodyObject(), TestCreateMockResponseFormHeadlessBodyString(), TestCreateMockResponseFormHeadlessHeadersFromFile() (+9 more)

### Community 15 - "File Path Source Forms"
Cohesion: 0.16
Nodes (16): BuildFilePathForm(), BuildSourceSelectionForm(), TestBuildFilePathForm(), TestBuildFilePathFormStructure(), TestBuildFilePathFormValidationEdgeCases(), TestBuildSourceSelectionForm(), TestBuildSourceSelectionFormOptions(), TestBuildSourceSelectionFormStructure() (+8 more)

### Community 16 - "Router and Endpoint Handlers"
Cohesion: 0.12
Nodes (9): TestEndpointHandler_BodylessResponse(), TestMockRoute_CreateRoute(), EndpointHandler, Handler, MainRouter, request(), TestRouter_CreateRoutes(), TestRouter_CreateRoutes_BothHandlersSet() (+1 more)

### Community 17 - "Cobra CLI Entrypoint"
Cohesion: 0.18
Nodes (8): Execute(), TestExecute(), TestExecute_ErrorPath(), TestRootCmd_RunsServer(), addServerFlags(), init(), runStartServer(), main()

### Community 18 - "Route Creation Forms"
Cohesion: 0.21
Nodes (12): createRouteForm(), NewHuhForm(), CollectBodyFromForm(), CollectHeadersFromForm(), createBodyForm(), createContinueForm(), createHeaderForm(), TestCollectBodyFromFormPublic() (+4 more)

### Community 19 - "Server Bootstrap"
Cohesion: 0.19
Nodes (11): TestCfg_GetServerAddr(), loadConfig(), runServer(), setupServer(), StartServer(), StartServerWithContext(), contains(), TestRunServer_ReturnsListenError() (+3 more)

### Community 20 - "Koanf Config Reader"
Cohesion: 0.26
Nodes (6): routeRows(), TestRouteRows(), addHealthCheckRoute(), normalizeRoutes(), koanfReader, Cfg

### Community 21 - "Reader Strategy Loader"
Cohesion: 0.20
Nodes (11): Reader, NewLoader(), TestNewLoader(), TestReader_Read_AllowsBodylessFakeResponse(), TestReader_Read_AppliesDefaultsBeforeValidation(), TestReader_Read_HealthRouteWithEmptyRoutes(), NewReaderStrategy(), TestNewReader() (+3 more)

### Community 22 - "Body and Headers File Loading"
Cohesion: 0.20
Nodes (12): newFailOpeningError(), newFailReadingError(), TestNewFailOpeningError(), TestNewFailReadingError(), LoadBodyFromFile(), LoadBodyStringFromFile(), LoadHeadersFromFile(), readFile() (+4 more)

### Community 23 - "Pre-commit Hooks"
Cohesion: 0.17
Nodes (11): default_install_hook_types, always_run, entry, id, language, name, pass_filenames, stages (+3 more)

### Community 24 - "GitHub Action Pins"
Cohesion: 0.18
Nodes (10): actions, actions/checkout, actions/download-artifact, actions/setup-go, actions/upload-artifact, golangci/golangci-lint-action, goreleaser/goreleaser-action, sonarsource/sonarqube-quality-gate-action (+2 more)

### Community 25 - "Headless Terminal Collectors"
Cohesion: 0.36
Nodes (10): CollectBodyString(), CollectHeaders(), headlessSelect(), TestCollectBodyHeadlessFileSource(), TestCollectBodyHeadlessSkipSource(), TestCollectBodyStringHeadlessFileSource(), TestCollectBodyStringHeadlessSkipSource(), TestCollectHeadersHeadlessFileSource() (+2 more)

### Community 26 - "Input Validators"
Cohesion: 0.20
Nodes (9): TestValidateHost(), TestValidateNonEmpty(), TestValidatePath(), TestValidatePositiveInt(), TestValidateStatusCode(), ValidateHost(), ValidatePath(), ValidatePositiveInt() (+1 more)

### Community 27 - "Form Collector Tests"
Cohesion: 0.22
Nodes (9): TestBuildFilePathFormValidationLogic(), CollectBodyStringFromForm(), TestCollectBodyFromForm(), TestCollectBodyStringFromForm(), TestCollectBodyStringFromFormPublic(), TestCollectHeadersFromForm(), TestFormErrorMessages(), TestCollectBodyStringFromFormHeadless() (+1 more)

### Community 28 - "Circuit Breaker Config Merge"
Cohesion: 0.31
Nodes (8): applyCircuitBreakerConfig(), DefaultCircuitBreakerConfig(), CircuitBreakerConfig, MergeCircuitBreakerConfig(), TestMergeCircuitBreakerConfig(), BuildCircuitBreakerRouteKey(), TestBuildCircuitBreakerRouteKey(), seedCircuitBreakers()

### Community 29 - "Config Write and Export"
Cohesion: 0.25
Nodes (8): exportAsJSON(), TestExportAsJSON(), TestStartServerCmd_ConfigFile(), TestHealthCheckRouteRoundTrip(), TestWriteConfig(), WriteConfig(), ResolveAbsolutePath(), TestResolveAbsolutePath()

### Community 30 - "Proxying and Mocking Concepts"
Cohesion: 0.33
Nodes (7): Circuit Breaker (proxy routes), fakeResponse Mock Handler, requestTo Proxy Handler, Circuit Breaker (CLAUDE guide), Proxying with Circuit Breaker, requestTo Overrides fakeResponse Rule, Circuit Breaker for Clients (shipped)

### Community 31 - "Reader Strategy Mocks"
Cohesion: 0.38
Nodes (4): TestReader_Read(), MockReaderStrategy, MockReaderStrategyMockRecorder, NewMockReaderStrategy()

### Community 32 - "Config Path Resolution"
Cohesion: 0.40
Nodes (6): DefaultCfg(), ReadOrCreateConfig(), ResolveConfigPath(), resolveConfigPath(), TestReadOrCreateConfig(), TestResolveConfigPath()

### Community 33 - "JSON Schema Sync"
Cohesion: 0.50
Nodes (5): joinPath(), modelPaths(), resolveSchemaRefs(), schemaPaths(), TestSchemaMatchesModel()

### Community 34 - "Bug Reporting Template"
Cohesion: 0.50
Nodes (4): Bug Report Issue Template, Bug Reproduction Steps, Issue Reporting Guidance, Private Vulnerability Disclosure

### Community 35 - "Recorder Release Notes"
Cohesion: 0.50
Nodes (4): Request Recorder Feature, Release 0.4.0, Recorder Admin API, Traffic Recording and Replay

### Community 36 - "Response Body Capture"
Cohesion: 0.50
Nodes (3): TestClientHandler_CreateHandler_NoNullBody(), captureBody(), TestCaptureBody()

### Community 37 - "Contribution and Conduct"
Cohesion: 0.67
Nodes (3): Pull Request Checklist, Contributor Covenant Code of Conduct, Contribution guide

### Community 38 - "PR Checklist and Setup"
Cohesion: 0.67
Nodes (3): PR Verification Checklist, Build and Test Commands, Development Setup

## Knowledge Gaps
- **37 isolated node(s):** `github.com/lynicis/inzibat`, `Bug Report Template`, `Feature Request Template`, `Pull Request Checklist`, `Contributor Covenant Code of Conduct` (+32 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 208 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Request Recording Store` to `Config and Handler Core`, `Circuit Breaker State Machine`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `github.com/lynicis/inzibat`, `Bug Report Template`, `Feature Request Template` to the rest of the system?**
  _37 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Request Recording Store` be split into smaller, more focused modules?**
  _Cohesion score 0.05706214689265537 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `HTTP Client Retry Logic` to `HTTP Client Tests`, `Circuit Breaker State Machine`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Should `Route Model and Forms` be split into smaller, more focused modules?**
  _Cohesion score 0.07364114552893045 - nodes in this community are weakly interconnected._
- **Why does `Cfg` connect `Koanf Config Reader` to `Config Path Resolution`, `Config and Handler Core`, `Route Model and Forms`, `Request Recording Store`, `Router and Endpoint Handlers`, `Server Bootstrap`, `Circuit Breaker Config Merge`, `Config Write and Export`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `GoReleaser Release Packaging` be split into smaller, more focused modules?**
  _Cohesion score 0.041666666666666664 - nodes in this community are weakly interconnected._