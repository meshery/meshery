# Project Overview

Meshery is a self-service engineering platform and the open source, cloud native manager enabling the
design and management of all Kubernetes-based infrastructure and applications. As a CNCF project, 
Meshery offers visual and collaborative GitOps, freeing you from the chains of YAML while managing 
Kubernetes multi-cluster deployments. With support for 300+ integrations, Meshery provides 
infrastructure lifecycle management, workspaces for team collaboration, design pattern catalogs,
performance management, and multi-tenancy capabilities across any cloud or on-premises environment.

## Repository Structure

- **/server** – Meshery Server backend written in Go; handles REST/GraphQL APIs, Kubernetes cluster
  management, adapter orchestration, and database operations (PostgreSQL).
- **/ui** – Meshery UI built with Next.js and React; includes Material UI components, Redux Toolkit
  state management, and Relay for GraphQL queries.
- **/mesheryctl** – Meshery CLI built with Go and Cobra framework; provides commands for installation,
  lifecycle management, pattern deployment, and system validation.
- **/docs** – Documentation site powered by Hugo; contains user guides, concepts, contributing
  guides, and API references.
- **/install** – Installation artifacts including Dockerfiles, Kubernetes manifests, Helm charts,
  Docker Compose files, and platform-specific deployment scripts.
- **/provider-ui** – Provider-specific UI extensions; isolated React application for remote provider
  integrations.
- **/policies** – Open Policy Agent (OPA) Rego policies for relationship evaluation and design
  validation.
- **/.github** – GitHub Actions workflows, issue templates, Copilot agent definitions, and community
  health files.

## Identifier Naming Conventions — MANDATORY

Full canonical directory: <https://github.com/meshery/schemas/blob/master/docs/identifier-naming-contributor-guide.md> — the reader-friendly 26-row naming table with before/after and do/don't examples. The inline per-layer forms below remain the repo-scoped authority; the guide is the ecosystem-wide reference.

> **Identifier-naming overhaul status (2026-04-28):** Complete in this repo at `v1.0.14`. Cross-repo cluster: `meshery/schemas` (`v1.2.6`), `meshery/meshkit` (`v1.0.7`), `layer5io/sistent` (`v0.20.1`), and `layer5io/meshery-cloud` (master HEAD, rolling). **`layer5labs/meshery-extensions` is deferred** pending lift of the layer5labs billing block — see `meshery/schemas/docs/identifier-naming-migration.md §21` for the post-completion landed-PR + tagged-release inventory.

This repository adheres to the canonical camelCase-wire identifier-naming contract
defined authoritatively in `meshery/schemas/AGENTS.md § Casing rules at a
glance`. The contract is **not optional**; deviations should be treated as
repository policy and corrected before review or merge. The cross-repo
consumer-audit gate lives in `meshery/schemas` and flags divergence in
this repo's server handlers and UI slices on every PR — the migration
that armed it is **complete (2026-04-23)** per
`meshery/schemas/docs/identifier-naming-impact-report.md`. See
`meshery/schemas/.github/workflows/schema-audit.yml` for the authoritative
CI job; the consumer-audit was promoted from advisory to **blocking** in
Phase 4.B.

### The rule in one sentence

*Wire is camelCase everywhere; DB is snake_case; Go fields follow Go
idiom; the ORM layer is the sole translation boundary.*

### Per-layer canonical forms

| Layer | Form |
|---|---|
| DB column / `db:` tag | `snake_case` — `user_id`, `org_id`, `created_at` |
| Go struct field | `PascalCase` with Go-idiomatic initialisms — `UserID`, `OrgID`, `CreatedAt` |
| JSON tag | `camelCase` — `json:"userId"`, `json:"orgId"`, `json:"createdAt"` |
| URL query/path param | `camelCase` — `{orgId}`, `?userId=...` |
| TypeScript property | `camelCase` — `response.userId`, `queryArg.orgId` |
| OpenAPI schema property | `camelCase` |
| OpenAPI `operationId` | `lower camelCase verbNoun` — `getWorkspaces` |
| `components/schemas` type name | `PascalCase` — `WorkspacePayload` |

### Forbidden (MUST NOT)

- MUST NOT introduce a `json:` tag that matches the `db:` tag on a new
  DB-backed field. Wire is camel; DB is snake; they differ by design.
- MUST NOT declare an RTK query endpoint hand-rolled when
  `@meshery/schemas/{mesheryApi,cloudApi}` provides a canonical equivalent.
- MUST NOT locally redeclare a Go type that has an equivalent in
  `github.com/meshery/schemas/models/...`.
- MUST NOT use `ID` (ALL CAPS) in URL query parameters, JSON tags, or
  TypeScript properties. `Id` (camelCase) is canonical.
- MUST NOT mix casing conventions within a single resource. If wire
  format must change, introduce a new API version per
  `meshery/schemas/AGENTS.md § Dual-Schema Pattern`.
- MUST NOT import deprecated legacy schema versions in new code. When
  importing from `@meshery/schemas`, consume the latest canonical-casing
  version (v1beta3 where present, otherwise v1beta2). Deprecated
  `v1beta1` constructs are kept for backward compatibility only; do not
  reach for them in new server handlers, `mesheryctl` commands, or UI
  RTK slices.

### Required on every PR

- MUST run the schemas validator locally before pushing. This command
  assumes `meshery` and `schemas` are checked out as sibling directories
  (for example, `../meshery` and `../schemas`):
  `cd ../schemas && make validate-schemas && make consumer-audit`. If
  `meshery/schemas` is cloned elsewhere, run the same targets from that
  checkout instead — e.g. `cd /path/to/schemas && make validate-schemas && make consumer-audit`.
- MUST include test updates for any casing or tag change.
- MUST include doc updates for any user-visible API change.
- MUST sign off commits (`git commit -s`).

### Authority

`meshery/schemas/AGENTS.md` is authoritative. On any conflict between
this repo's documentation and the schemas AGENTS.md, schemas wins. File
discrepancies as issues against `meshery/schemas`, not locally. All
wire-format questions are resolved against `meshery/schemas`, not this
downstream repo.

### Migration

The identifier-naming migration is tracked at
`meshery/schemas/docs/identifier-naming-migration.md`. All
contributors — human and AI agents — MUST read this plan before making
any schema-aware change.

**Status: complete (2026-04-23).** Per
`meshery/schemas/docs/identifier-naming-impact-report.md`, all 22
in-scope resources have been migrated to canonical camelCase-on-wire
versions, the consumer-audit CI gate has been promoted to blocking
(Phase 4.B), and consumer-audit TypeScript findings against this repo
are at zero (Phase 2 tail PR #18904 + the per-resource Phase 3 repoint
PRs #18886, #18888, #18889, #18890, #18894, #18900). This repository's
side of the migration is therefore **closed**; new work should not
introduce snake_case wire properties or hand-rolled local Go types
duplicating canonical schemas types.

#### Server-side post-migration audit (2026-04-23)

A retrospective canonical-conformance audit confirmed that this repo's
server side is conformant with the contract:

- **Local Go duplicates of canonical schemas types — none requiring
  retirement.** The named candidates from the migration plan (§2.3) are
  resolved as follows:
  - `MesheryPattern` (`server/models/meshery_pattern.go`) — retained as
    a deliberate **persistence/storage adapter shim**, not a duplicate.
    The canonical `v1beta3/design.MesheryPattern` expresses the design
    body as a typed `*PatternFile` struct; the server's storage model
    persists the design as a YAML/JSON string in a single `pattern_file`
    column. The local struct is the adapter that bridges the canonical
    design representation and the server's persisted storage model and
    cannot be displaced without first restructuring server-side
    persistence — out of scope for an identifier-naming migration. The
    five count fields (`viewCount`, `shareCount`, `downloadCount`,
    `cloneCount`, `deploymentCount`) and `orgId` already wear canonical
    camelCase JSON tags, and `UnmarshalJSON` dual-accepts the legacy
    snake_case spellings for the deprecation window per the Phase 2.K
    cascade contract.
  - `MesheryPatternRequestBody` — not present locally. Of the two
    locally defined wrappers, only `MesheryPatternUPDATERequestBody`
    remains on an active request path; the POST handler uses
    `DesignPostPayload`, which contains `design.PatternFile`, while
    `MesheryPatternPOSTRequestBody` is deprecated. These local request
    shapes dual-accept `patternData` / `pattern_data` wrapper keys per
    the cascade contract.
  - `MesheryFilter` (`server/models/meshery_filter.go`) — canonical
    schemas `v1beta3/design.MesheryFilter` is a placeholder
    (`map[string]interface{}`), not a typed struct. Local retention is
    correct until schemas authors a structured filter resource.
  - `MesheryApplication` (`server/models/meshery_application.go`) — no
    canonical equivalent exists in `meshery/schemas`. Local retention
    is correct.
  - `PerformanceProfile` (`server/models/performance_profiles.go`) — no
    canonical equivalent exists in `meshery/schemas`. Local retention
    is correct.
- **DB-column conformance — all canonical-declared `db:` tags
  present.** Every column declared via `db:` on
  `v1beta3/design.MesheryPattern` (`clone_count`, `created_at`,
  `deployment_count`, `download_count`, `pattern_file`, `share_count`,
  `updated_at`, `user_id`, `view_count`) has a corresponding gorm-mapped
  column in `models.MesheryPattern`'s `AutoMigrate` registration in
  `server/cmd/main.go`. No migration action required.
- **Consumer audit clean.** The `consumer-audit` blocking CI in
  `meshery/schemas` reports zero TypeScript findings against this repo
  post-merge of PR #18904.

**Reviewer guardrail: the residual local types listed above are
intentional and MUST NOT be flagged as canonical-duplicate violations
in code review.** They are documented as "intentional request-wrapper
shims" in the impact report's row 42. If a future change makes
displacement feasible (e.g., schemas adds a structured `MesheryFilter`
or the server's persistence model is refactored), the displacement
should be coordinated cross-repo per `meshery/schemas/AGENTS.md`'s
"Source of Truth" directive — not unilaterally in this repo.

## Build & Development Commands

### Server (Go)

```bash
# Run Meshery Server locally (port 9081)
make server

# Run with local provider for testing
make server-local

# Build server binary
make build-server

# Run server binary
make server-binary

# Lint Go code
make golangci

# Run server without Kubernetes components
make server-skip-compgen

# Run server without operator deployment
make server-without-operator

# Generate error codes and update helpers
make error
```

### UI (Next.js/React)

```bash
# Install UI dependencies
make ui-setup

# Run UI development server (port 3000)
make ui

# Build and export UI
make ui-build

# Build only Meshery UI
make ui-meshery-build

# Lint UI code
make ui-lint

# Run end-to-end tests
make ui-integration-tests
```

### CLI (mesheryctl)

```bash
# Build mesheryctl binary
cd mesheryctl && make

# Run unit tests
cd mesheryctl && go test --short ./...

# Run integration tests
cd mesheryctl && go test -run Integration ./...

# Generate CLI documentation
make docs-mesheryctl
```

### Docker

```bash
# Build Meshery container
make docker-build

# Build playground mode container
make docker-playground-build

# Run Meshery with Docker pointing to production Remote Provider
make docker-cloud

# Run Meshery with Docker pointing to local Remote Provider
make docker-local-cloud
```

### Documentation

```bash
# Run docs site with and listen for changes (port 1313)
make docs

# Build docs site
make docs-build

# Run docs in Docker container
make docs-docker
```

### API Specifications

```bash
# REST API docs consume the published OpenAPI spec at docs/data/openapi.yml

# Build GraphQL schema
make graphql-build

# GraphQL is self-documenting via the introspection endpoint and the GraphQL
# Playground at http://localhost:9081/api/system/graphql/playground.
# A static reference snapshot lives at docs/content/en/reference/graphql-apis.md.
```

### Helm Charts

```bash
# Lint all Helm charts
make helm-lint

# Generate Helm chart documentation
make helm-docs
```

## Code Style & Conventions

### Go Code

- **Formatting**: Use `gofmt` and `goimports`; enforced via `golangci-lint`.
- **Linting**: Run `make golangci` before committing; config in `.golangci.yml`.
- **Error Handling**: Use MeshKit's error utilities (`github.com/meshery/meshkit/errors`); run
  `make error` to analyze error codes.
- **Package Structure**: Follow Go project layout standards; keep packages focused and cohesive.
- **Testing**: Use Go's standard testing library; place tests in `*_test.go` files; run with 
  `go test`.
- **Dependencies**: Manage with Go modules; run `go mod tidy` regularly.

### JavaScript/React Code

- **Formatting**: ESLint with Prettier; config in `ui/.eslintrc.js`.
- **Component Style**: Use functional components with React hooks; avoid class components.
- **Styling**: Prefer Sistent design system components (`@sistent/sistent`); fall back to Material
  UI (MUI) when components unavailable.
- **State Management**: Redux Toolkit for global state; local state with `useState`/`useReducer`.
- **API Integration**: GraphQL via Relay; REST via Axios.
- **Testing**: Playwright for E2E tests; run with `make ui-integration-tests`.
- **Schema-Driven Development**: Use JSON schemas to define component props and validation.

### Commit Messages

- **Format**: `[component] descriptive message` (e.g., `[UI] Add workspace filter dropdown`).
- **Sign-Off**: All commits must be signed with DCO (`git commit -s`).
- **References**: Link to issues/PRs when applicable (`Fixes #1234`, `Relates to #5678`).

### Git Workflow

- Fork the repository and create feature branches from `master`.
- Keep commits atomic and focused on single changes.
- Squash commits before merging unless history provides value.

## Architecture Notes

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Meshery UI (Next.js)                        │
│  ┌───────────────┐  ┌──────────────┐  ┌──────────────────────────┐ │
│  │   Material UI │  │ Redux Toolkit│  │  Relay (GraphQL Client)  │ │
│  │   Components  │  │  State Mgmt  │  │  + REST (Axios)          │ │
│  └───────────────┘  └──────────────┘  └──────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                              │
                         HTTP/WebSocket
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                      Meshery Server (Go)                            │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────┐ ┌────────────┐ │
│  │  REST API   │  │  GraphQL API │  │ PostgreSQL │ │   NATS     │ │
│  │  (Gorilla)  │  │   (gqlgen)   │  │     DB     │ │  Messaging │ │
│  └─────────────┘  └──────────────┘  └────────────┘ └────────────┘ │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │           Provider Plugins (gRPC/Remote Providers)           │ │
│  └──────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                              │
                    gRPC / Kubernetes API
                              │
┌─────────────────────────────────────────────────────────────────────┐
│                     Kubernetes Clusters                             │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────────────────┐ │
│  │   Meshery   │  │   MeshSync   │  │  Adapters (Istio, Linkerd, │ │
│  │  Operator   │  │ (Discovery)  │  │   Consul, NSM, etc.)       │ │
│  └─────────────┘  └──────────────┘  └────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Interactions

**Meshery UI** serves as the primary user interface, communicating with Meshery Server via REST and
GraphQL APIs. The UI uses Relay for GraphQL subscriptions (real-time updates) and Axios for REST
endpoints. Redux Toolkit manages application state including user sessions, workspace data, and
connection configurations.

**Meshery Server** is the core orchestration engine written in Go. It exposes REST (port 9081) and
GraphQL APIs, manages connections to Kubernetes clusters, handles user authentication via provider
plugins, stores data in PostgreSQL, and publishes events to NATS. The server uses MeshKit for common
utilities, error handling, and database abstractions.

**Meshery Operator** runs inside Kubernetes clusters to manage lifecycle of Meshery components
including MeshSync (for cluster discovery and state synchronization) and Meshery Broker (NATS for
event streaming).

**Adapters** are gRPC services that provide integration with specific service meshes and cloud native
infrastructure. They run as separate containers and communicate with Meshery Server over gRPC.

**mesheryctl** is the CLI tool for local development, deployment, and operations. It communicates
directly with Meshery Server APIs and can manage local Docker containers or Kubernetes deployments.

### Data Flow

1. User interacts with Meshery UI (design canvas, configuration forms, performance tests).
2. UI sends GraphQL mutations/queries or REST requests to Meshery Server.
3. Server processes requests, applies business logic, and persists to PostgreSQL.
4. Server communicates with Kubernetes API or adapter gRPC endpoints to deploy resources.
5. MeshSync discovers cluster state changes and publishes to NATS.
6. Server receives NATS events and updates database.
7. UI receives real-time updates via GraphQL subscriptions.

## Testing Strategy

### Unit Tests (Go)

- **Location**: `*_test.go` files alongside source code.
- **Framework**: Go standard library `testing` package.
- **Run**: `go test ./...` or `go test --short ./...` (skip integration tests).
- **Coverage**: Aim for ≥ 70% coverage on business logic.

### Unit Tests (JavaScript)

> TODO: Expand JavaScript unit testing infrastructure.

### Integration Tests (Go)

- **Location**: `server/integration-tests/`
- **Setup**: `make server-integration-tests-meshsync-setup` (creates kind cluster, deploys
  operator).
- **Run**: `make server-integration-tests-meshsync-run`
- **Cleanup**: `make server-integration-tests-meshsync-cleanup`
- **Requirements**: Docker, kind, kubectl, helm.

### End-to-End Tests (UI)

- **Location**: `ui/` (Playwright tests).
- **Framework**: Playwright.
- **Setup**: `make test-setup-ui` (installs Playwright browsers with dependencies).
- **Run**: `make ui-integration-tests` or `npm run test:e2e` in `ui/`.
- **CI**: `make test-e2e-ci` runs in non-interactive mode.

### CI/CD Testing

- **Workflows**: `.github/workflows/` contains GitHub Actions for PR checks, release builds, and
  deployments.
- **Automated Tests**: `build-ui-and-server.yml` runs on PRs; includes linting, builds, and test
  execution.
- **CodeQL**: `codeql-analysis.yml` performs security scanning.

### Local Testing Best Practices

- Run `make golangci` before committing Go code.
- Run `make ui-lint` before committing UI code.
- Test server changes with `make server` and verify on `http://localhost:9081`.
- Test UI changes with `make ui` and verify on `http://localhost:3000`.
- Run integration tests when modifying Kubernetes or MeshSync interactions.

## Security & Compliance

### Reporting Vulnerabilities

- **Email**: [security@meshery.dev](mailto:security@meshery.dev)
- **Response Time**: Acknowledged within 10 business days.
- **Full Policy**: [SECURITY.md](./SECURITY.md)

### Secrets Management

- **Never commit secrets**: API keys, tokens, passwords, certificates.
- **Use environment variables**: Server accepts secrets via env vars (e.g., `PROVIDER_BASE_URLS`,
  `KEYS_PATH`).
- **GitHub Secrets**: Store in repository secrets for CI/CD; access via `${{ secrets.VAR_NAME }}`.

### Dependency Scanning

- **Go**: Dependabot enabled; monitors `go.mod` and `go.sum`.
- **JavaScript**: Dependabot monitors `ui/package.json` and `provider-ui/package.json`.
- **SBOM**: Software Bill of Materials generated via `.github/workflows/bom.yaml`.

### Security Scanning

- **CodeQL**: Automated security analysis runs on every PR and push to `master`.
- **OpenSSF Scorecard**: Tracks security posture; badge in README.
- **License Compliance**: Apache 2.0; verify dependencies are compatible.

### Best Practices

- Run `make golangci` to catch security lints (e.g., SQL injection, command injection).
- Validate all user inputs; sanitize data before rendering in UI.
- Use parameterized queries for database operations.
- Limit permissions on Kubernetes service accounts and RBAC roles.

## Agent Guardrails

### Files That Should Not Be Modified by Agents

- **LICENSE** – Apache 2.0 license text; do not edit.
- **CODE_OF_CONDUCT.md** – Community standards; changes require human review.
- **GOVERNANCE.md** – Project governance; changes require maintainer consensus.
- **MAINTAINERS.md** – Maintainer list; requires verification of identity.
- **.github/copilot-instructions.md** – Base Copilot instructions; changes require careful review.
- **.github/agents/** – Agent definition files; self-modification risks infinite loops.
- **go.sum**, **ui/package-lock.json**, **provider-ui/package-lock.json** – Lock files; only update
  via package managers.

### Required Human Reviews

- **Security changes**: All changes touching authentication, authorization, secrets, or encryption.
- **Database migrations**: Schema changes require careful review to avoid data loss.
- **API breaking changes**: REST/GraphQL API modifications that break backward compatibility.
- **Helm chart templates**: Changes to `install/kubernetes/helm/` require validation in test clusters.
- **CI/CD workflows**: `.github/workflows/` changes may affect release process.

### Rate Limits & Constraints

- **PR size**: Keep PRs focused and under 500 lines changed when possible.
- **Commit frequency**: Batch related changes; avoid excessive micro-commits.
- **Build failures**: Do not merge if CI checks fail; investigate and fix root cause.

### Quality Standards

- All Go code must pass `make golangci`.
- All JavaScript code must pass `make ui-lint`.
- New features require corresponding documentation updates.
- Breaking changes require deprecation notices and migration guides.

## Extensibility Hooks

### Provider Plugins

- **Interface**: `server/models/provider.go` defines the provider contract.
- **Remote Providers**: Implement authentication, preferences, and sync logic externally.
- **Example**: Meshery Cloud provider integrates via HTTPS/gRPC.

### Adapters (gRPC)

- **Protocol**: Defined in `server/meshes/meshops.proto`.
- **Examples**: meshery-istio, meshery-linkerd, meshery-consul (separate repositories).
- **Registration**: Adapters self-register with Meshery Server on startup.

### UI Extensions

- **Remote Components**: Load React components from URLs at runtime via `@paciolan/remote-component`.
- **Provider UI**: `provider-ui/` directory for provider-specific UI extensions.

### GraphQL Subscriptions

- **Location**: `server/internal/graphql/schema.graphql`
- **Extend**: Add new queries, mutations, or subscriptions; run `make graphql-build`.

### Feature Flags

- **Environment Variables**: Control behavior via env vars (e.g., `PLAYGROUND`, `DEBUG`,
  `SKIP_COMP_GEN`).
- **Runtime Config**: Server reads config from `~/.meshery/config.yaml`.

### Event System

- **NATS Topics**: Publish/subscribe to topics like `meshsync.request`, `meshery.broker`.
- **MeshSync Events**: Listen for cluster state changes and trigger workflows.

### Hooks & Scripts

- **Pre-commit**: Husky hooks in `ui/.husky/` run linting before commits.
- **Makefile Targets**: Extend `Makefile` or `install/Makefile.core.mk` for custom build steps.

## Coding Agents

Reusable agent definitions live in `.agents/`. These are LLM-agnostic — any coding assistant can use them.

| Agent | File | Purpose |
|-------|------|---------|
| Code Reviewer | `.agents/code-reviewer.md` | Parallel code review across Go + frontend |
| Security Reviewer | `.agents/security-reviewer.md` | Security audit for infrastructure-managing code |
| Meshery Code Contributor | `.agents/meshery-code-contributor.md` | Full-stack code contributions (Go, React, CLI) |
| Meshery Docs Contributor | `.agents/meshery-docs-contributor.md` | Hugo-based documentation contributions |
| GitHub Actions Engineer | `.agents/github-actions-engineer.md` | CI/CD workflow design and debugging |
| Relationship Fixture Agent | `.agents/relationship-fixture-agent.md` | Create relationship test fixture designs |

## Skills

Packaged, repeatable workflows live in `.agents/skills/`. Each skill has a `SKILL.md` with instructions.

| Skill | Directory | Invocation | Purpose |
|-------|-----------|------------|---------|
| gen-test | `.agents/skills/gen-test/` | User-invoked | Generate idiomatic Go tests |
| api-doc | `.agents/skills/api-doc/` | User or agent | Document REST/GraphQL endpoints |

## Automation Hooks

Standalone hook scripts in `.agents/hooks/` can be wired into any coding agent's hook system or run manually:

| Hook | Script | Trigger | Purpose |
|------|--------|---------|---------|
| Format Frontend | `.agents/hooks/format-frontend.sh` | Post-edit | Auto-format JS/TS files with Prettier |
| Block Lock Files | `.agents/hooks/block-lockfiles.sh` | Pre-edit | Prevent direct edits to lock files |

## Further Reading

- [Contributing Guide](./CONTRIBUTING.md) – Start here for onboarding.
- [Meshery Documentation](https://docs.meshery.io) – Full user and contributor guides.
- [Architecture Overview](https://docs.meshery.io/concepts/architecture) – Detailed component
  diagrams and explanations.
- [API Documentation](https://docs.meshery.io/extensibility/api) – REST and GraphQL API references.
- [CLI Guide](https://docs.meshery.io/guides/mesheryctl) – mesheryctl command reference.
- [Extensibility](https://docs.meshery.io/extensibility) – Provider plugins, adapters, and extension
  points.
- [Community Handbook](https://meshery.io/community#handbook) – Community processes and resources.
- [Roadmap](./ROADMAP.md) – Upcoming features and milestones.
- [Security Policy](./SECURITY.md) – Vulnerability reporting and disclosure process.
- [Governance](./GOVERNANCE.md) – Project decision-making and maintainer responsibilities.
