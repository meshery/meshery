# Project Overview

Meshery is a self-service engineering platform and open source cloud native manager for Kubernetes infrastructure. A CNCF project supporting 300+ integrations with visual GitOps, multi-cluster management, and workspace collaboration.

## Repository Structure

| Directory | Purpose |
|-----------|---------|
| `/server` | Go backend — REST/GraphQL APIs, Kubernetes management, PostgreSQL |
| `/ui` | Next.js/React frontend — MUI, Redux Toolkit, Relay GraphQL |
| `/mesheryctl` | Go CLI with Cobra — install, lifecycle, pattern deployment |
| `/docs` | Hugo documentation site |
| `/install` | Dockerfiles, Kubernetes manifests, Helm charts |
| `/provider-ui` | Provider-specific React UI extensions |
| `/.github` | GitHub Actions, issue templates, Copilot agent definitions |

## Identifier Naming Conventions — MANDATORY

Authoritative guide: <https://github.com/meshery/schemas/blob/master/docs/identifier-naming-contributor-guide.md>

**Wire is camelCase; DB is snake_case; Go fields follow Go idiom; the ORM layer is the sole translation boundary.**

### Per-layer canonical forms

| Layer | Form |
|---|---|
| DB column / `db:` tag | `snake_case` — `user_id`, `org_id`, `created_at` |
| Go struct field | `PascalCase` with Go initialisms — `UserID`, `OrgID`, `CreatedAt` |
| JSON tag | `camelCase` — `json:"userId"`, `json:"orgId"` |
| URL query/path param | `camelCase` — `{orgId}`, `?userId=...` |
| TypeScript property | `camelCase` — `response.userId` |
| OpenAPI schema property | `camelCase` |
| OpenAPI `operationId` | `lower camelCase verbNoun` — `getWorkspaces` |
| `components/schemas` type name | `PascalCase` — `WorkspacePayload` |

### Forbidden (MUST NOT)

- MUST NOT use a `json:` tag matching the `db:` tag — wire is camel, DB is snake.
- MUST NOT hand-roll an RTK query endpoint when `@meshery/schemas/{mesheryApi,cloudApi}` provides one.
- MUST NOT locally redeclare a Go type with an equivalent in `github.com/meshery/schemas/models/...`.
- MUST NOT use `ID` (ALL CAPS) in URL params, JSON tags, or TypeScript properties — use `Id`.
- MUST NOT mix casing within a single resource; introduce a new API version to change wire format.
- MUST NOT import deprecated `v1beta1` in new code; use `v1beta3` (or `v1beta2` where v1beta3 absent).

### Required on every PR

- Run schemas validator: `cd ../schemas && make validate-schemas && make consumer-audit`
- Include test updates for casing/tag changes.
- Include doc updates for user-visible API changes.
- Sign off commits: `git commit -s`

> `meshery/schemas/AGENTS.md` is authoritative. On conflicts, schemas wins.

## API Changes — MUST Go Through Schemas — MANDATORY

**Any new or changed HTTP API (new endpoint, new/renamed query param, new
request/response field) MUST be defined in `meshery/schemas` first and consumed
via the generated client. Do NOT hand-roll RTK Query endpoints, response types,
or ad-hoc `fetch`/`axios` calls for an API that can live in schemas.**

Schemas is the single source of truth: one OpenAPI definition drives the Go
models, TypeScript types, and the RTK Query client (`@meshery/schemas/{mesheryApi,cloudApi}`)
consumed here. Hand-rolling any of these silently diverges the wire contract
across `meshery/meshery` and `meshery-cloud`.

### Workflow (adding/updating an endpoint)

1. **Define** the path + schemas in the matching construct's `api.yml`
   (e.g. `../schemas/schemas/constructs/v1beta1/system/api.yml` for `/api/system/*`,
   `.../connection/api.yml` for `/api/integrations/connections*`). Follow the
   schemas conventions: `operationId` = lower-camel `verbNoun`, camelCase wire
   params/properties, `x-internal: ["meshery"]` for Meshery-only endpoints,
   `additionalProperties: false`, `maxLength` on strings.
2. **Regenerate** in `../schemas`: `make bundle-openapi generate-rtk generate-golang`
   (or `make build` for the full dist). Verify the new `useXQuery` /
   `mesheryApi.endpoints.X` hooks appear.
3. **Validate**: `cd ../schemas && make validate-schemas && make consumer-audit`.
4. **Consume** the generated hook in the UI (import from `@meshery/schemas/mesheryApi`;
   wrap in `ui/rtk-query/*` only for thin ergonomics like bare-id args or cache
   tags — never to re-declare the request). Use the generated Go models on the
   server where applicable.
5. **Release coupling**: schemas releases are automated ("do not manually create
   releases"). Until a new `@meshery/schemas` is published and this repo's
   dependency is bumped, a **local link** is used for development
   (`ui/package.json` → `"@meshery/schemas": "file:../schemas"` and the
   `replace github.com/meshery/schemas => ../schemas` directive in `go.mod`).
   Both the version bump and reverting the local link happen as part of the
   normal release/upgrade flow — do not commit the local link as the permanent
   dependency.

### Narrow exceptions (still prefer schemas)

- **Server-Sent Events / streaming**: RTK codegen can't produce a useful hook
  for `text/event-stream`. Still **document** the endpoint in `api.yml`, but
  consume it with a native `EventSource` client under `ui/lib/*`
  (e.g. `ui/lib/controllersStatusSubscription.ts`).
- Truly Meshery-internal endpoints with no cross-repo consumer may skip schemas,
  but must be justified in the PR description.

### Forbidden

- MUST NOT add a `builder.query`/`builder.mutation` in `ui/rtk-query/*` that
  issues a request to an API which is (or should be) defined in schemas.
- MUST NOT hand-write response/param TypeScript types or Go structs that
  duplicate a schemas-generated type.
- MUST NOT change wire casing/field names only in this repo — change the schema
  and regenerate (see the naming conventions above).

## Build & Development Commands

- Use the `gh-axi` CLI tool to interact with GitHub. Prefer `gh-axi` over `gh`.
- Use `chrome-devtools-axi` for browser automation (navigate, snapshot, click, fill forms, run JS, inspect console/network) in place of raw Playwright/chrome-devtools MCP for ad hoc tasks.
- Run `quota-axi` to check local agent-provider quota windows before long-running work.
- Use the `lavish` skill (`lavish-axi` CLI) to turn a plan, comparison, or report into a reviewable HTML artifact.

### Server (Go)

```bash
make server                    # Run server locally (port 9081)
make server-local              # Run with local provider
make build-server              # Build binary
make golangci                  # Lint Go code
make server-skip-compgen       # Run without Kubernetes components
make server-without-operator   # Run without operator deployment
make error                     # Generate error codes
```

### UI (Next.js/React)

```bash
make ui-setup              # Install dependencies
make ui                    # Dev server (port 3000)
make ui-build              # Build and export
make ui-lint               # Lint UI code
make ui-integration-tests  # Run E2E tests
```

### CLI (mesheryctl)

```bash
cd mesheryctl && make                       # Build binary
cd mesheryctl && go test --short ./...      # Unit tests
cd mesheryctl && go test -run Integration ./...  # Integration tests
make docs-mesheryctl                        # Generate CLI docs
```

### Docker

```bash
make docker-build           # Build container
make docker-cloud           # Run with production Remote Provider
make docker-local-cloud     # Run with local Remote Provider
```

### Documentation

```bash
make docs        # Run docs site (port 1313)
make docs-build  # Build docs site
```

### API & Helm

```bash
make graphql-build  # Build GraphQL schema
make helm-lint      # Lint Helm charts
make helm-docs      # Generate Helm chart docs
```

## Code Style & Conventions

### Go

- Format with `gofmt`/`goimports`; lint with `make golangci` (config: `.golangci.yml`).
- Use MeshKit error utilities (`github.com/meshery/meshkit/errors`); run `make error` for codes.
- Tests in `*_test.go`; manage deps with `go mod tidy`.

### JavaScript/React

- ESLint + Prettier (config: `ui/.eslintrc.js`).
- Functional components with hooks; no class components.
- Use `@sistent/sistent` design system; fall back to MUI.
- Redux Toolkit for global state; GraphQL via Relay; REST via Axios.
- Playwright for E2E tests.

### Commits

- Format: `[component] descriptive message` (e.g., `[UI] Add workspace filter dropdown`)
- Sign off: `git commit -s`
- Reference issues: `Fixes #1234`

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                  Meshery UI (Next.js)                │
│   MUI Components │ Redux Toolkit │ Relay + Axios     │
└──────────────────────────┬───────────────────────────┘
                     HTTP/WebSocket
┌──────────────────────────┴───────────────────────────┐
│                 Meshery Server (Go)                  │
│   REST (9081) │ GraphQL │ PostgreSQL │ NATS          │
│         Provider Plugins (gRPC/Remote)               │
└──────────────────────────┬───────────────────────────┘
                  gRPC / Kubernetes API
┌──────────────────────────┴───────────────────────────┐
│              Kubernetes Clusters                     │
│  Meshery Operator │ MeshSync │ Adapters (gRPC)       │
└──────────────────────────────────────────────────────┘
```

**Data flow**: UI → REST/GraphQL → Server → PostgreSQL + Kubernetes API → NATS → MeshSync → GraphQL subscriptions → UI.

## Testing

### Go

- Unit: `go test ./...` or `go test --short ./...`
- Integration setup: `make server-integration-tests-meshsync-setup` (requires Docker, kind, kubectl, helm)
- Integration run: `make server-integration-tests-meshsync-run`
- Target ≥70% coverage on business logic.

### UI

- E2E (Playwright): `make ui-integration-tests` or `npm run test:e2e` in `ui/`
- Setup: `make test-setup-ui`

### Local Validation

```bash
make golangci    # before Go commits
make ui-lint     # before UI commits
```

## Security & Compliance

- Report vulnerabilities: [security@meshery.dev](mailto:security@meshery.dev) — acknowledged in 10 business days.
- Never commit secrets; use env vars (`PROVIDER_BASE_URLS`, `KEYS_PATH`) and GitHub Secrets.
- CodeQL runs on every PR; OpenSSF Scorecard tracks security posture.
- Apache 2.0 license — verify dependency compatibility.
- Use parameterized queries; validate/sanitize all user inputs.

## Agent Guardrails

### Do Not Modify

`LICENSE`, `CODE_OF_CONDUCT.md`, `GOVERNANCE.md`, `MAINTAINERS.md`, `.github/copilot-instructions.md`, `.github/agents/`, `go.sum`, `ui/package-lock.json`, `provider-ui/package-lock.json`

### Require Human Review

- Security changes (auth, secrets, encryption)
- Database migrations
- API breaking changes
- Helm chart templates (`install/kubernetes/helm/`)
- CI/CD workflows (`.github/workflows/`)

### Quality Gates

- Go: `make golangci` must pass
- JS: `make ui-lint` must pass
- New features need docs; breaking changes need deprecation notices
- Keep PRs under 500 lines; don't merge on CI failure

## Extensibility

### Provider Plugins

Interface: `server/models/provider.go` — implement auth, preferences, and sync externally.

### Adapters (gRPC)

Protocol: `server/meshes/meshops.proto` — adapters self-register on startup. Examples: meshery-istio, meshery-linkerd, meshery-consul.

### UI Extensions

Remote Components loaded via `@paciolan/remote-component`. Bundle **must** expose `module.exports = { default: Component, __esModule: true }`. Silent-undefined gotcha: a mis-built bundle renders as `undefined` with no loader error — React throws "Element type is invalid" at render. Check bundle export shape first. See `ui/components/layout/Navigator/NavigatorExtension.tsx`.

### GraphQL

Schema: `server/internal/graphql/schema.graphql`. Add queries/mutations/subscriptions then run `make graphql-build`.

### Feature Flags

Env vars: `PLAYGROUND`, `DEBUG`, `SKIP_COMP_GEN`. Runtime config: `~/.meshery/config.yaml`.

### Event System

NATS topics: `meshsync.request`, `meshery.broker`. MeshSync publishes cluster state changes.

### Hooks & Scripts

- Pre-commit: Husky hooks in `ui/.husky/`
- Build: extend `Makefile` or `install/Makefile.core.mk`

## Coding Agents

Agent definitions in `.agents/` (LLM-agnostic):

| Agent | File | Purpose |
|-------|------|---------|
| Code Reviewer | `.agents/code-reviewer.md` | Parallel review across Go + frontend |
| Security Reviewer | `.agents/security-reviewer.md` | Security audit |
| Meshery Code Contributor | `.agents/meshery-code-contributor.md` | Full-stack contributions |
| Meshery Docs Contributor | `.agents/meshery-docs-contributor.md` | Hugo docs contributions |
| GitHub Actions Engineer | `.agents/github-actions-engineer.md` | CI/CD design and debugging |
| Relationship Fixture Agent | `.agents/relationship-fixture-agent.md` | Relationship test fixtures |

## Skills

Packaged workflows in `.agents/skills/`:

| Skill | Directory | Purpose |
|-------|-----------|---------|
| gen-test | `.agents/skills/gen-test/` | Generate idiomatic Go tests |
| api-doc | `.agents/skills/api-doc/` | Document REST/GraphQL endpoints |
| gen-relationship | `.agents/skills/gen-relationship/` | Generate schema-backed relationships |

## Automation Hooks

Scripts in `.agents/hooks/`:

| Hook | Script | Trigger | Purpose |
|------|--------|---------|---------|
| Format Frontend | `.agents/hooks/format-frontend.sh` | Post-edit | Auto-format JS/TS with Prettier |
| Block Lock Files | `.agents/hooks/block-lockfiles.sh` | Pre-edit | Prevent direct edits to lock files |

## Further Reading

- [Contributing Guide](./CONTRIBUTING.md)
- [Meshery Documentation](https://docs.meshery.io)
- [Architecture Overview](https://docs.meshery.io/concepts/architecture)
- [API Documentation](https://docs.meshery.io/extensibility/api)
- [CLI Guide](https://docs.meshery.io/guides/mesheryctl)
- [Extensibility](https://docs.meshery.io/extensibility)
- [Community Handbook](https://meshery.io/community#handbook)
- [Security Policy](./SECURITY.md)
- [Governance](./GOVERNANCE.md)
