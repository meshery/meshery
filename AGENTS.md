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

Workflow: define the path in the construct's `api.yml` → regenerate in `../schemas`
(`make bundle-openapi generate-rtk generate-golang`) → validate
(`make validate-schemas && make consumer-audit`) → consume the generated hook.
The full steps, the SSE/streaming exception, the Go models-only rule, and the
release/local-link coupling are in
[Consuming schemas from meshery/meshery](./docs/content/en/project/contributing/contributing-schemas.md).

### Forbidden

- MUST NOT add a `builder.query`/`builder.mutation` in `ui/rtk-query/*` that
  issues a request to an API which is (or should be) defined in schemas.
- MUST NOT hand-write response/param TypeScript types or Go structs that
  duplicate a schemas-generated type.
- MUST NOT change wire casing/field names only in this repo — change the schema
  and regenerate (see the naming conventions above).
- MUST NOT re-declare a generated endpoint to attach a cache tag — use the
  callback form of `enhanceEndpoints` (`appendInvalidatesTags` from
  `ui/rtk-query/utils`). The object form replaces schemas-side tags wholesale.
- MUST NOT add a local endpoint whose name a generated one already defines — it
  is silently discarded dead code and a different request goes over the wire.
  Check first:
  `grep '<name>:t\.' ui/node_modules/@meshery/schemas/dist/mesheryApi.js`.
- MUST NOT keep a local copy of a contract Meshery only decodes from (or encodes
  to) a remote provider — it is the schemas construct, aliased. Where a Go type
  must stay local because it doubles as a GORM model, keep the JSON tags
  identical and pin the column with `gorm:"column:..."`.
- MUST NOT index a credential's persisted `secret` map directly — it has four
  shapes and legacy rows are never rewritten. Delegate to
  `server/models/credential_secret.go` or `ui/utils/credentialSecret.ts`, which
  must keep the same resolution rules.

Detail behind these rules:
[RTK Query integration](./docs/content/en/project/contributing/ui/schemas.md)
(cache tags, the deliberate-override exception, testing the *effective*
endpoint), [consuming schemas](./docs/content/en/project/contributing/contributing-schemas.md)
(consumed contracts, the Go side, rename propagation), and
[connections](./docs/content/en/project/contributing/models/connections.md)
(the credential shape catalogue).

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

A build failing `compile: version "goX.Y.Z" does not match go tool version` across many
dependencies is a stale `GOROOT`, not a code problem — see
[build environment gotchas](./docs/content/en/project/contributing/contributing-build-environment.md).

### UI (Next.js/React)

```bash
make ui-setup              # Install dependencies
make ui                    # Dev server (port 3000)
make ui-build              # Build and export
make ui-lint               # Lint UI code
make ui-integration-tests  # Run E2E tests
```

- Never commit `ui/tsconfig.tsbuildinfo` — it is tracked and not gitignored, so any
  local `tsc --noEmit` dirties it. Stage explicit paths, never `git add -A`.
- A `@sistent/sistent` bump must cover all three consuming manifests (`ui/`,
  `provider-ui/`, `install/docker-extension/ui/`) and sistent's peers.
- After any local sibling-package build, re-verify with `npm ci` before trusting a
  green run — matching version strings do not prove the installed contents are the
  published ones.

### CLI (mesheryctl)

```bash
cd mesheryctl && make                       # Build binary
cd mesheryctl && go test --short ./...      # Unit tests
cd mesheryctl && go test -run Integration ./...  # Integration tests
make docs-mesheryctl                        # Generate CLI docs
```

`make docs-mesheryctl` rewrites ~100 pages with your local `$HOME` baked in — revert
every page whose only change is that path. See
[build environment gotchas](./docs/content/en/project/contributing/contributing-build-environment.md).

### Releasing

Meshery has **no automatic release cadence**. Release Drafter keeps one draft release
current on every push to `master`; publishing it creates the `v*` tag that fires the
stable build fan-out. Follow `.agents/skills/cut-release/SKILL.md` — never hand-author
a tag or notes, and never trust `gh release edit --draft=false`'s exit code as proof of
publication. See
[cutting a release](./docs/content/en/project/contributing/build-and-release.md).

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
  `make error` skips `mesheryctl`. Both components require bumping `next_error_code` in their
  own `helpers/component_info.json` **in the same commit**, and the tracked docs reference at
  `docs/data/errorref/` must be regenerated or the new codes are silently omitted. Full
  contract: [writing MeshKit errors](./docs/content/en/project/contributing/contributing-error.md).
- Only `utils.Log.Error(err)` renders a MeshKit error's code, cause and remediation; cobra's
  default print shows just the message. In `mesheryctl` commands, log the structured error
  for the user *and* return it for the exit path.
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

Golden-file workflow (`-args -update`, the `fixtures/` vs `testdata/` split, and
the rule that a regenerated golden must still encode *intended* behavior) is
documented in `docs/content/en/project/contributing/cli/cli.md`.

**A rename in `meshery/schemas` renames the gorm column too**, because gorm derives
it from the Go *field name*, not the `db:` tag. After bumping schemas, grep every raw
column reference (`Select`, `Where`, `Order`, `Joins`, `Scan`, migrations) and the
`mesheryctl` fixtures for the old spelling, and propagate gorm errors so the next such
rename fails loudly. Detail and regression test:
[consuming schemas](./docs/content/en/project/contributing/contributing-schemas.md).

### UI

- E2E (Playwright): `make ui-integration-tests` or `npm run test:e2e` in `ui/`
- Setup: `make test-setup-ui`

Four E2E invariants. Each one's reasoning, evidence, run IDs and history live in
`docs/content/en/project/contributing/ui/tests.md` - read it before changing the workflow
or the setup projects.

- **The job gates on the Playwright verdict.** The final step of
  `.github/workflows/test-e2e.yml` keys on `steps.playwright-tests.outcome`, never
  `.conclusion` (`continue-on-error` pins `conclusion` to `success`). Never re-disarm it to
  turn a red build green; fix the test or `test.fixme` it with the tracking issue in the
  annotation.
- **The remote-provider credential is the `REMOTE_PROVIDER_TEST_USER_TOKEN` org secret**,
  spelled identically in caller and reusable workflow and asserted non-empty before any
  other work. Never alias it between the two - that is what hid a secret which existed
  nowhere for three months.
- **A missing credential must fail, never `setup.skip()`.** Playwright collapses a
  dependent project only when its setup *fails*; a skip leaves every dependent scheduled to
  die on a storage state file that was never written.
- **A local run needs three things** - `make ui-provider-build`, a server on `:9081` plus
  `make ui` on `:3000`, and `MESHERY_SERVER_URL=http://localhost:3000` on the Playwright
  run - and the failure when one is missing never names it.

### Local Validation

```bash
make golangci    # before Go commits
make ui-lint     # before UI commits
```

### QA Allure traceability (Test-Group-keyed reports)

Test results feed the [meshery/qa](https://github.com/meshery/qa) Allure dashboard
(qa.meshery.io), where each report is a filtered view over one shared result pool
**keyed on the Test Plan Test Group** (Latest tab, col B) via the `testGroup`
label. Emit it from both lanes: CLI via the `[tg=<Test Group>]` `@test` title
token (parser + full token→label mapping in `mesheryctl/bats-to-allure.js`), UI
via `allure.label('testGroup', …)` (see `ui/tests/e2e/connections.spec.ts` and
the sheet↔code map `ui/tests/e2e/connections.testmap.ts`). Contract docs:
`docs/content/en/project/contributing/{cli/tests.md,build-and-release.md}`.
"Connection Lifecycle" is the first such report.

## Security & Compliance

- Report vulnerabilities: [security@meshery.dev](mailto:security@meshery.dev) — acknowledged in 10 business days.
- Never commit secrets; use env vars (`PROVIDER_BASE_URLS`, `KEYS_PATH`) and GitHub Secrets.
- CodeQL runs on every PR; OpenSSF Scorecard tracks security posture.
- Apache 2.0 license — verify dependency compatibility.
- Use parameterized queries; validate/sanitize all user inputs.

## Agent Guardrails

### Do Not Modify

`LICENSE`, `CODE_OF_CONDUCT.md`, `GOVERNANCE.md`, `MAINTAINERS.md`, `.github/copilot-instructions.md`, `.github/agents/`

Never hand-edit a generated lock file (`go.sum`, or any of the several `package-lock.json`
files). Regenerate it with the package manager - a dependency bump legitimately rewrites
every lock file it touches. `.agents/hooks/block-lockfiles.sh` enforces this by basename,
so it covers lock files this list does not enumerate.

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

Remote Components loaded via `@paciolan/remote-component`. Two rules:

- A bundle **must** expose `module.exports = { default: Component, __esModule: true }`;
  built without `output.library.type = "commonjs2"` it resolves to `undefined` with no
  loader error.
- Derive every event literal and injected capability key from `@sistent/sistent`'s
  `mesheryExtensionContract` module rather than typing strings, and keep
  `ui/utils/eventBus.ts` typed as `EventBus<MesheryExtensionEvent>` — a bare
  `new EventBus()` disables publish-site checking entirely.

Detail: [UI Extensions](./docs/content/en/project/contributing/ui/ui.md).

### GraphQL

Schema: `server/internal/graphql/schema.graphql`. Add queries/mutations/subscriptions then run `make graphql-build`.

### Feature Flags

Env vars: `PLAYGROUND`, `DEBUG`, `SKIP_COMP_GEN`. Runtime config: `~/.meshery/config.yaml`.

### Event System

NATS topics: `meshsync.request`, `meshery.broker`. MeshSync publishes cluster state changes.

### Hooks & Scripts

- Pre-commit: Husky hooks in `ui/.husky/`
- Build: extend `Makefile` or `install/Makefile.core.mk`

## Agent Tooling

`.agents/` is the LLM-agnostic home for agent definitions, packaged skills and
automation hooks. The agent table, per-tool skill discovery, and the reasoning behind
the rules below are in [.agents/README.md](./.agents/README.md) — read it before
touching any of them.

- **`.agents/skills/` is the single source of truth for skills**, one directory per
  skill. Add a new skill only there; never enumerate them in this file, list the
  directory.
- **`.claude/skills` is a symlink to `../.agents/skills` and nothing else.** Never
  replace it with real directories or copies, and never run the AXI installer against
  a real checkout to see what happens — the failure mode under test destroys the
  canonical skill content.
- **Skill content addresses its own files as `.agents/skills/...`**, never through
  `.claude/`, which resolves only where the symlink exists.
- Hooks in `.agents/hooks/`: `format-frontend.sh` (post-edit Prettier) and
  `block-lockfiles.sh` (pre-edit lock-file guard).

## Further Reading

The rules above are complete on their own. These files hold the reasoning, evidence and
worked detail behind them — open the one that matches what you are working on.

| Working on | Read |
|---|---|
| A schema, an API contract, or a consumer of one | [Contributing to Meshery Schemas](./docs/content/en/project/contributing/contributing-schemas.md) |
| RTK Query, generated hooks, cache tags | [Schema-Driven UI Development](./docs/content/en/project/contributing/ui/schemas.md) |
| Playwright E2E, the E2E CI job, its credentials | [UI End-to-End Tests](./docs/content/en/project/contributing/ui/tests.md) |
| A build that fails for an unrelated-looking reason | [Build Environment Gotchas](./docs/content/en/project/contributing/contributing-build-environment.md) |
| MeshKit error codes | [How to write MeshKit compatible errors](./docs/content/en/project/contributing/contributing-error.md) |
| Releases, CI secrets, the QA dashboard | [Build & Release (CI)](./docs/content/en/project/contributing/build-and-release.md) |
| Connections and credential secrets | [Connections](./docs/content/en/project/contributing/models/connections.md) |
| UI extensions, Remote Components | [Contributing to Meshery UI](./docs/content/en/project/contributing/ui/ui.md) |
| `mesheryctl`, golden files | [Contributing to Meshery CLI](./docs/content/en/project/contributing/cli/cli.md) |
| Agent definitions, skills, hooks | [.agents/README.md](./.agents/README.md) |

External: [Meshery Documentation](https://docs.meshery.io) ·
[Architecture](https://docs.meshery.io/concepts/architecture) ·
[Extensibility](https://docs.meshery.io/extensibility) ·
[Community Handbook](https://meshery.io/community#handbook) ·
[Contributing Guide](./CONTRIBUTING.md) · [Security Policy](./SECURITY.md) ·
[Governance](./GOVERNANCE.md)

## Maintaining this file

Keep this file for knowledge useful to almost every future agent session in this project.

**This file states rules; the files in Further Reading hold the detail.** A session must
be able to learn that a constraint exists without opening anything else, so every
MANDATORY rule and prohibition stays here, short and flat. Worked examples, causal
history, run IDs, dates, evidence and per-endpoint walkthroughs belong in the linked
file that owns the subject, with a pointer from here.

Do not repeat what the codebase already shows; point to the authoritative file or command
instead. Relocate rather than delete, and prefer updating an existing doc over creating a
new one. When you add a pointer, add it to the Further Reading table too.
