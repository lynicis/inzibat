# Graph Report - inzibat  (2026-11-17)

## Corpus Check
- 92 files · ~37,627 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 7 file(s) not represented in the graph (top: (none) 3, .toml 3, .properties 1)

## Summary
- 806 nodes · 2205 edges · 40 communities (30 shown, 10 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 293 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `96bf2a1a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Request Recording Store
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
- port_test.go
- Body and Headers File Loading
- Pre-commit Hooks
- GitHub Action Pins
- testing.T
- Proxying and Mocking Concepts
- JSON Schema Sync
- Bug Reporting Template
- Recorder Release Notes
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
- `CollectBody()` --references--> `HttpBody`  [EXTRACTED]
  cmd/formbuilder/data_collector.go → config/model.go
- `CollectBodyFromForm()` --references--> `HttpBody`  [EXTRACTED]
  cmd/formbuilder/form_collector.go → config/model.go
- `createClientRequestForm()` --references--> `RequestTo`  [EXTRACTED]
  cmd/create.go → config/model.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Agent Guide Architecture Contract** — agents_entrypoints, agents_architecture, agents_config_precedence, agents_cli_commands, agents_circuit_breaker [EXTRACTED 1.00]
- **Inzibat Core Feature Set** — readme_circuit_breaker_proxying, readme_traffic_recording, readme_interactive_route_builder, readme_schema_validated_config [EXTRACTED 1.00]
- **Contribution Workflow Documents** — contributing_pull_request_process, contributing_issue_reporting, _github_pull_request_template_pull_request_template, _github_issue_template_bug_report_bug_report_template, _github_issue_template_feature_request_feature_request_template [INFERRED 0.85]
- **Repository meta documents** — readme_inzibatoverview, contributing_contributionguide, changelog_releasehistory, security_securitypolicy, code_of_conduct_codeofconduct [INFERRED 0.85]

## Communities (40 total, 10 thin omitted)

### Community 1 - "Request Recording Store"
Cohesion: 0.05
Nodes (46): adminURL(), exportAsInzibatConfig(), exportAsJSON(), fetchRecordedEntries(), fetchRecordedSession(), getJSON(), TestExportAsInzibatConfig(), TestExportAsJSON() (+38 more)

### Community 2 - "Route Model and Forms"
Cohesion: 0.08
Nodes (39): collectCircuitBreakerConfig(), collectClientRequestBody(), createClientRequestFormInternal(), createMockResponseFormInternal(), createRouteInternal(), parsePositiveInt(), TestCollectCircuitBreakerConfig(), TestCreateClientRequestFormInternal() (+31 more)

### Community 3 - "GoReleaser Release Packaging"
Cohesion: 0.04
Nodes (47): archives, format_overrides, formats, name_template, before, hooks, builds, env (+39 more)

### Community 4 - "CI Workflow Pipeline"
Cohesion: 0.06
Nodes (38): needs, runs-on, steps, needs, runs-on, steps, env, GO_VERSION (+30 more)

### Community 5 - "Release and Distribution Workflow"
Cohesion: 0.08
Nodes (28): needs, runs-on, steps, CHOCOLATEY_API_KEY, GITHUB_TOKEN, HOMEBREW_TAP_GITHUB_TOKEN, jobs, chocolatey (+20 more)

### Community 6 - "BoolPointer"
Cohesion: 0.26
Nodes (13): BoolPointer(), NewCircuitBreakerStore(), TestCircuitBreaker(), TestCircuitBreaker_HalfOpenGaugeReleased(), TestCircuitBreakerStore_Allow(), TestCircuitBreakerStore_Get(), TestCircuitBreakerStore_OnFailure(), TestCircuitBreakerStore_OnSuccess() (+5 more)

### Community 7 - "Circuit Breaker State Machine"
Cohesion: 0.08
Nodes (19): Client, isIdempotent(), isRetryableError(), IsUnforwardableHeader(), allowHalfOpen(), CircuitBreakerStore, openRecord(), releaseHalfOpen() (+11 more)

### Community 8 - "Project Docs and Templates"
Cohesion: 0.08
Nodes (23): Bug Report Template, Feature Request Issue Template, Feature Request Template, Pull Request Template, Graphify Code Intelligence Index, Graphify Index Replacing GitNexus, JSON Schema for Configuration, Release history (Keep a Changelog) (+15 more)

### Community 9 - "circuitBreaker"
Cohesion: 0.18
Nodes (10): circuitBreaker, enabled, failureThreshold, halfOpenMaxRequests, minimumRequests, openTimeoutMs, successThreshold, routes (+2 more)

### Community 10 - "Homebrew Tap Update Workflow"
Cohesion: 0.08
Nodes (23): jobs, update-tap, name, on, workflow_run, permissions, contents, id (+15 more)

### Community 11 - "Circuit Breaker Config Schema"
Cohesion: 0.09
Nodes (22): testKey, xTestHeader, body, circuitBreaker, headers, host, method, passWithRequestBody (+14 more)

### Community 12 - "Lint Configuration"
Cohesion: 0.09
Nodes (22): max-complexity, generated, paths, presets, formatters, enable, exclusions, lines (+14 more)

### Community 13 - "HTTP Client Tests"
Cohesion: 0.16
Nodes (23): DefaultRetryConfig(), NewHttpClient(), TestClient_buildRequest(), TestClient_calculateBackoff(), TestClient_Delete(), TestClient_executeRequest(), TestClient_Get(), TestClient_handleResponse() (+15 more)

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
Cohesion: 0.40
Nodes (4): setupServer(), TestRunServer_ReturnsListenError(), TestSetupServer_CircuitBreaker(), TestSetupServer_RecordEnabled()

### Community 20 - "Cfg"
Cohesion: 0.05
Nodes (48): routeRows(), TestRouteRows(), TestStartServerCmd_ConfigFile(), addHealthCheckRoute(), DefaultCfg(), Reader, NewLoader(), normalizeRoutes() (+40 more)

### Community 22 - "Body and Headers File Loading"
Cohesion: 0.20
Nodes (12): newFailOpeningError(), newFailReadingError(), TestNewFailOpeningError(), TestNewFailReadingError(), LoadBodyFromFile(), LoadBodyStringFromFile(), LoadHeadersFromFile(), readFile() (+4 more)

### Community 23 - "Pre-commit Hooks"
Cohesion: 0.17
Nodes (11): default_install_hook_types, always_run, entry, id, language, name, pass_filenames, stages (+3 more)

### Community 24 - "GitHub Action Pins"
Cohesion: 0.18
Nodes (10): actions, actions/checkout, actions/download-artifact, actions/setup-go, actions/upload-artifact, golangci/golangci-lint-action, goreleaser/goreleaser-action, sonarsource/sonarqube-quality-gate-action (+2 more)

### Community 26 - "testing.T"
Cohesion: 0.11
Nodes (21): TestCollectBody(), TestCollectBodyString(), TestCollectHeaders(), TestFormPlaceholders(), TestValidateHost(), TestValidatePath(), TestValidatePositiveInt(), TestValidateStatusCode() (+13 more)

### Community 30 - "Proxying and Mocking Concepts"
Cohesion: 0.33
Nodes (7): Circuit Breaker (proxy routes), fakeResponse Mock Handler, requestTo Proxy Handler, Circuit Breaker (CLAUDE guide), Proxying with Circuit Breaker, requestTo Overrides fakeResponse Rule, Circuit Breaker for Clients (shipped)

### Community 33 - "JSON Schema Sync"
Cohesion: 0.50
Nodes (5): joinPath(), modelPaths(), resolveSchemaRefs(), schemaPaths(), TestSchemaMatchesModel()

### Community 34 - "Bug Reporting Template"
Cohesion: 0.50
Nodes (4): Bug Report Issue Template, Bug Reproduction Steps, Issue Reporting Guidance, Private Vulnerability Disclosure

### Community 35 - "Recorder Release Notes"
Cohesion: 0.50
Nodes (4): Request Recorder Feature, Release 0.4.0, Recorder Admin API, Traffic Recording and Replay

### Community 37 - "Contribution and Conduct"
Cohesion: 0.67
Nodes (3): Pull Request Checklist, Contributor Covenant Code of Conduct, Contribution guide

### Community 38 - "PR Checklist and Setup"
Cohesion: 0.67
Nodes (3): PR Verification Checklist, Build and Test Commands, Development Setup

## Knowledge Gaps
- **37 isolated node(s):** `github.com/lynicis/inzibat`, `fakeResponse Mock Handler`, `Circuit Breaker (CLAUDE guide)`, `Bug Reproduction Steps`, `Private Vulnerability Disclosure` (+32 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 221 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **10 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Request Recording Store` to `Config and Handler Core`, `Circuit Breaker State Machine`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `github.com/lynicis/inzibat`, `fakeResponse Mock Handler`, `Circuit Breaker (CLAUDE guide)` to the rest of the system?**
  _37 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Request Recording Store` be split into smaller, more focused modules?**
  _Cohesion score 0.05311676909569798 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `Circuit Breaker State Machine` to `HTTP Client Tests`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Should `Route Model and Forms` be split into smaller, more focused modules?**
  _Cohesion score 0.07764876632801161 - nodes in this community are weakly interconnected._
- **Why does `Cfg` connect `Cfg` to `Config and Handler Core`, `Request Recording Store`, `Route Model and Forms`, `Router and Endpoint Handlers`, `setupServer`?**
  _High betweenness centrality (0.009) - this node is a cross-community bridge._
- **Should `GoReleaser Release Packaging` be split into smaller, more focused modules?**
  _Cohesion score 0.041666666666666664 - nodes in this community are weakly interconnected._