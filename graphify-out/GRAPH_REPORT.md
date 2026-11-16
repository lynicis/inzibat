# Graph Report - inzibat  (2026-11-16)

## Corpus Check
- 92 files · ~37,232 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 7 file(s) not represented in the graph (top: (none) 3, .toml 3, .properties 1)

## Summary
- 745 nodes · 2145 edges · 27 communities (23 shown, 4 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 272 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `04822477`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Interactive Form Internals
- Recorder CLI
- form_builder_test.go
- GoReleaser Config
- Config Loading
- CI Workflow
- Release Workflow
- testing.T
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
- withHeadlessTerminal
- Pre-commit Hooks
- Actup Config
- Validators Tests
- JSON Schema
- Roadmap Items
- Security Policy
- Module Root

## God Nodes (most connected - your core abstractions)
1. `NewHttpClient()` - 22 edges
2. `Store` - 21 edges
3. `Cfg` - 21 edges
4. `FormRunner` - 17 edges
5. `BoolPointer()` - 16 edges
6. `NewCircuitBreakerStore()` - 15 edges
7. `CircuitBreakerStore` - 14 edges
8. `setupServer()` - 14 edges
9. `GetFreePort()` - 14 edges
10. `Client` - 13 edges

## Surprising Connections (you probably didn't know these)
- `TestClientHandler_CreateHandler_NoNullBody()` --calls--> `captureBody()`  [INFERRED]
  handler/client_test.go → recorder/middleware.go
- `exportAsJSON()` --references--> `RecordedSession`  [EXTRACTED]
  cmd/record.go → recorder/model.go
- `exportAsInzibatConfig()` --calls--> `WriteConfig()`  [EXTRACTED]
  cmd/record.go → config/config.go
- `ClientHandler` --references--> `Client`  [EXTRACTED]
  handler/client.go → client/http/http.go
- `buildFakeResponse()` --references--> `FakeResponse`  [EXTRACTED]
  recorder/converter.go → config/model.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Repository meta documents** — agents_agentguide, readme_inzibatoverview, contributing_contributionguide, changelog_releasehistory, security_securitypolicy, code_of_conduct_codeofconduct [INFERRED 0.85]

## Communities (27 total, 4 thin omitted)

### Community 1 - "Interactive Form Internals"
Cohesion: 0.05
Nodes (67): collectCircuitBreakerConfig(), collectClientRequestBody(), createClientRequestForm(), createClientRequestFormInternal(), createMockResponseForm(), createMockResponseFormInternal(), createRoute(), createRouteForm() (+59 more)

### Community 2 - "Recorder CLI"
Cohesion: 0.05
Nodes (45): adminURL(), exportAsInzibatConfig(), fetchRecordedEntries(), fetchRecordedSession(), getJSON(), TestExportAsInzibatConfig(), TestFetchRecordedEntries(), TestFetchRecordedSession() (+37 more)

### Community 3 - "form_builder_test.go"
Cohesion: 0.16
Nodes (16): BuildFilePathForm(), BuildSourceSelectionForm(), TestBuildFilePathForm(), TestBuildFilePathFormStructure(), TestBuildFilePathFormValidationEdgeCases(), TestBuildSourceSelectionForm(), TestBuildSourceSelectionFormOptions(), TestBuildSourceSelectionFormStructure() (+8 more)

### Community 4 - "GoReleaser Config"
Cohesion: 0.04
Nodes (47): archives, format_overrides, formats, name_template, before, hooks, builds, env (+39 more)

### Community 5 - "Config Loading"
Cohesion: 0.05
Nodes (45): routeRows(), TestRouteRows(), exportAsJSON(), TestExportAsJSON(), TestStartServerCmd_ConfigFile(), addHealthCheckRoute(), DefaultCfg(), Reader (+37 more)

### Community 6 - "CI Workflow"
Cohesion: 0.06
Nodes (38): needs, runs-on, steps, needs, runs-on, steps, env, GO_VERSION (+30 more)

### Community 7 - "Release Workflow"
Cohesion: 0.08
Nodes (28): needs, runs-on, steps, CHOCOLATEY_API_KEY, GITHUB_TOKEN, HOMEBREW_TAP_GITHUB_TOKEN, jobs, chocolatey (+20 more)

### Community 8 - "testing.T"
Cohesion: 0.11
Nodes (20): TestCollectBody(), TestCollectBodyString(), TestCollectHeaders(), TestBuildFilePathFormValidationLogic(), TestCollectBodyFromForm(), TestCollectBodyStringFromForm(), TestCollectHeadersFromForm(), TestFormErrorMessages() (+12 more)

### Community 9 - "HTTP Client Tests"
Cohesion: 0.15
Nodes (24): DefaultRetryConfig(), NewHttpClient(), TestClient_buildRequest(), TestClient_calculateBackoff(), TestClient_Delete(), TestClient_executeRequest(), TestClient_Get(), TestClient_handleResponse() (+16 more)

### Community 10 - "Circuit Breaker Store"
Cohesion: 0.08
Nodes (33): applyCircuitBreakerConfig(), BoolPointer(), DefaultCircuitBreakerConfig(), CircuitBreakerConfig, MergeCircuitBreakerConfig(), TestMergeCircuitBreakerConfig(), allowHalfOpen(), BuildCircuitBreakerRouteKey() (+25 more)

### Community 11 - "HTTP Client Core"
Cohesion: 0.15
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
Cohesion: 0.20
Nodes (12): newFailOpeningError(), newFailReadingError(), TestNewFailOpeningError(), TestNewFailReadingError(), LoadBodyFromFile(), LoadBodyStringFromFile(), LoadHeadersFromFile(), readFile() (+4 more)

### Community 18 - "Root Command"
Cohesion: 0.25
Nodes (7): Execute(), TestExecute(), TestExecute_ErrorPath(), addServerFlags(), init(), runStartServer(), main()

### Community 19 - "withHeadlessTerminal"
Cohesion: 0.21
Nodes (15): CollectBody(), CollectBodyString(), CollectHeaders(), headlessSelect(), TestCollectBodyFromFormHeadless(), TestCollectBodyHeadlessFileSource(), TestCollectBodyHeadlessSkipSource(), TestCollectBodyStringFromFormHeadless() (+7 more)

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
- **9 isolated node(s):** `github.com/lynicis/inzibat`, `Roadmap: circuit breaker support (done)`, `Bug Report Template`, `Feature Request Template`, `Contributor Covenant Code of Conduct` (+4 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 179 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Store` connect `Recorder CLI` to `Route Creation Forms`, `Circuit Breaker Store`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `github.com/lynicis/inzibat`, `Roadmap: circuit breaker support (done)`, `Bug Report Template` to the rest of the system?**
  _9 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Route Creation Forms` be split into smaller, more focused modules?**
  _Cohesion score 0.14311310190369542 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `HTTP Client Core` to `HTTP Client Tests`, `Circuit Breaker Store`?**
  _High betweenness centrality (0.012) - this node is a cross-community bridge._
- **Should `Interactive Form Internals` be split into smaller, more focused modules?**
  _Cohesion score 0.051055867415129644 - nodes in this community are weakly interconnected._
- **Why does `Cfg` connect `Config Loading` to `Route Creation Forms`, `Interactive Form Internals`, `Recorder CLI`, `Circuit Breaker Store`, `Endpoint Handlers`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Should `Recorder CLI` be split into smaller, more focused modules?**
  _Cohesion score 0.052884615384615384 - nodes in this community are weakly interconnected._