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
- **/docs** – Documentation site powered by Jekyll; contains user guides, concepts, contributing
  guides, and API references.
- **/install** – Installation artifacts including Dockerfiles, Kubernetes manifests, Helm charts,
  Docker Compose files, and platform-specific deployment scripts.
- **/provider-ui** – Provider-specific UI extensions; isolated React application for remote provider
  integrations.
- **/policies** – Open Policy Agent (OPA) Rego policies for relationship evaluation and design
  validation.
- **/.github** – GitHub Actions workflows, issue templates, Copilot agent definitions, and community
  health files.

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
# Run docs site with live reload (port 4000)
make docs

# Build docs site
make docs-build

# Run docs in Docker container
make docs-docker
```

### API Specifications

```bash
# Build Swagger/OpenAPI specs
make swagger-build

# Serve Swagger UI
make swagger

# Build GraphQL schema
make graphql-build

# Generate GraphQL documentation
make graphql-docs-build
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
│  └─────────────┘  └──────────────┘  └────────────┘ └────────────┐ │
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
