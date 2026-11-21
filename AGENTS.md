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

**Data flow**: UI → REST/GraphQL (HTTP/WebSocket) → Server (REST :9081, GraphQL, PostgreSQL, NATS) → Kubernetes API (Operator, MeshSync, gRPC adapters) → NATS → MeshSync → GraphQL subscriptions → UI.

## Identifier Naming Conventions — MANDATORY

Authoritative guide: <https://github.com/meshery/schemas/blob/master/docs/identifier-naming-contributor-guide.md>. `meshery/schemas/AGENTS.md` wins on conflict.

**Wire is camelCase; DB is snake_case; Go fields follow Go idiom; the ORM layer is the sole translation boundary.**

| Layer | Form |
|---|---|
| DB column / `db:` tag | `snake_case` — `user_id`, `created_at` |
| Go struct field | `PascalCase` with Go initialisms — `UserID`, `CreatedAt` |
| JSON tag, URL param, TypeScript property, OpenAPI property | `camelCase` — `json:"userId"`, `{orgId}` |
| OpenAPI `operationId` | `lower camelCase verbNoun` — `getWorkspaces` |
| `components/schemas` type name | `PascalCase` — `WorkspacePayload` |

MUST NOT:

- use a `json:` tag that matches the `db:` tag.
- use `ID` (ALL CAPS) in URL params, JSON tags, or TypeScript properties — use `Id`.
- mix casing within one resource; introduce a new API version to change wire format.
- locally redeclare a Go type that exists in `github.com/meshery/schemas/models/...`.
- import deprecated `v1beta1` in new code; use `v1beta3` (or `v1beta2` where v1beta3 is absent).

Every PR: run `cd ../schemas && make validate-schemas && make consumer-audit`; include test updates for casing/tag changes and doc updates for user-visible API changes.

## API Changes — MUST Go Through Schemas — MANDATORY

Any new or changed HTTP API (endpoint, query param, request/response field) MUST be defined in `meshery/schemas` first and consumed via the generated client (`@meshery/schemas/{mesheryApi,cloudApi}`). Schemas is the single source of truth for Go models, TypeScript types and RTK Query; hand-rolling any of them diverges the wire contract with `meshery-cloud`.

Workflow: define the path in the construct's `api.yml` → `../schemas`: `make bundle-openapi generate-rtk generate-golang` → `make validate-schemas && make consumer-audit` → consume the generated hook. Full steps and exceptions: [Consuming schemas](./docs/content/en/project/contributing/contributing-schemas.md).

MUST NOT:

- add a `builder.query`/`builder.mutation` in `ui/rtk-query/*` for an API that is (or should be) in schemas, or hand-write types/structs that duplicate a generated one.
- change wire casing or field names only in this repo — change the schema and regenerate.
- re-declare a generated endpoint to attach a cache tag; use the callback form of `enhanceEndpoints` (`appendInvalidatesTags` from `ui/rtk-query/utils`).
- add a local endpoint whose name a generated one already defines (it is silently dead code). Check: `grep '<name>:t\.' ui/node_modules/@meshery/schemas/dist/mesheryApi.js`.
- keep a local copy of a contract Meshery only decodes from or encodes to a remote provider. A Go type that doubles as a GORM model keeps identical JSON tags and pins the column with `gorm:"column:..."`.
- index a credential's persisted `secret` map directly; use `server/models/credential_secret.go` or `ui/utils/credentialSecret.ts`.

Detail: [RTK Query integration](./docs/content/en/project/contributing/ui/schemas.md), [Connections](./docs/content/en/project/contributing/models/connections.md).

## Build & Development Commands

- Prefer the `gh-axi` CLI over `gh`; use `chrome-devtools-axi` for browser automation; run `quota-axi` before long work; use the `lavish` skill to turn a plan or report into a reviewable HTML artifact.

```bash
make server                   # run the server (:9081)
make server-local             # run with the local provider
make build-server             # build the server binary
make server-skip-compgen      # run without Kubernetes components
make server-without-operator  # run without the operator deployment
make golangci                 # lint Go
make error                    # generate error codes
make ui-setup                 # install UI dependencies
make ui                       # UI dev server (:3000)
make ui-build                 # build and export the UI
make ui-lint                  # lint the UI
make ui-integration-tests     # UI E2E tests
cd mesheryctl && make         # build the CLI
cd mesheryctl && go test --short ./...             # CLI unit tests
cd mesheryctl && go test -run Integration ./...    # CLI integration tests
make docs-mesheryctl          # regenerate CLI docs
make docker-build             # build the container
make docker-cloud             # run with the production Remote Provider
make docker-local-cloud       # run with a local Remote Provider
make docs                     # Hugo docs (:1313)
make docs-build               # build the docs site
make graphql-build            # build the GraphQL schema
make helm-lint                # lint Helm charts
make helm-docs                # CAUTION: overwrites the hand-maintained meshery-operator README; read the note at the top of install/kubernetes/helm/meshery-operator/README.md first
```

- A build failing `compile: version "goX.Y.Z" does not match go tool version` is a stale `GOROOT`, not a code problem.
- `make docs-mesheryctl` bakes your `$HOME` into ~100 pages; revert every page whose only change is that path.
- Docs need a real **Dart Sass** on `PATH` (`make -C docs check-deps` accepts a Ruby `sass` gem, then fails).
- A docs page with images belongs in a leaf bundle (`<page>/index.md` + `<page>/images/`).
- Never commit `ui/tsconfig.tsbuildinfo`; stage explicit paths, never `git add -A`.
- A `@sistent/sistent` bump must cover `ui/`, `provider-ui/` and `install/docker-extension/ui/` plus sistent's peers; re-verify with `npm ci` after any local sibling build.
- Releases: there is no automatic cadence. Follow `.agents/skills/cut-release/SKILL.md`; never hand-author a tag or notes.

More: [build environment gotchas](./docs/content/en/project/contributing/contributing-build-environment.md), [cutting a release](./docs/content/en/project/contributing/build-and-release.md), [contributing to docs](./docs/content/en/project/contributing/contributing-docs/docs.md).

## Code Style & Conventions

### Go

- Format with `gofmt`/`goimports`; lint with `make golangci` (`.github/.golangci.yml` plus analyzers in `server/internal/lint/`).
- MUST NOT pass an unsanitized dynamic string to gorm `.Order(...)` (the sink behind every Meshery CVE). Use `models.SanitizeOrderInput(order, []string{...})` with that query's snake_case allow-list, or a constant, **before** the call. Enforced by the `orderby` analyzer. See [Go lint rules](./docs/content/en/project/contributing/contributing-lint.md).
- Use MeshKit errors (`github.com/meshery/meshkit/errors`); `make error` skips `mesheryctl`. Both components bump `next_error_code` in their own `helpers/component_info.json` **in the same commit**, and `docs/data/errorref/` must be regenerated. See [writing MeshKit errors](./docs/content/en/project/contributing/contributing-error.md).
- MUST NOT resolve a provider by name from `HandlerConfig.Providers` in boot-time code (`PROVIDER` enforcement can make `Providers[LocalProviderName]` nil). Persist off-request events through `HandlerConfig.SystemEventPersister`; index `Providers` only comma-ok and only on the request path. Boot-time work that can fault belongs in `models.RunSeedStage`, and a stage that spawns a goroutine defers its own `recover` at the spawn site. See [Extensibility: Providers](./docs/content/en/reference/extensibility/providers/index.md).
- In `mesheryctl`, `utils.Log.Error(err)` is the only call that renders a MeshKit error's code, cause and remediation: log it for the user *and* return it.
- Tests in `*_test.go`; manage deps with `go mod tidy`.

### JavaScript/React

- ESLint + Prettier (`ui/.eslintrc.js`); functional components with hooks, no classes.
- Use `@sistent/sistent`, falling back to MUI; read its `DESIGN.md` from the UI package you work in (`ui/`, `provider-ui/`, `install/docker-extension/ui/`), matching the resolved version.
- Redux Toolkit for global state; GraphQL via Relay; REST via Axios; Playwright for E2E.
- Every content-bearing page needs an **access** gate, not just gated controls: `useHasPermission(Keys.X)` from `@sistent/sistent`, `<DefaultError permissionKey={Keys.X} />` when denied, and `skip: !canX` on the page's RTK Query hooks. Exception: the dashboard (`/`). Pin the **deny** path in a test.
- Under `ui/pages/`, render `DefaultError` as a branch *inside* `MesheryPage`, never an early return; keep tests out of `ui/pages/` (they become routes) and put them in `ui/__tests__/`. Details and where each layer's gate goes: [Extensibility: Authorization](./docs/content/en/reference/extensibility/authorization/index.md).

### Commits

- Format `[component] descriptive message` (e.g. `[UI] Add workspace filter dropdown`); reference issues with `Fixes #1234`.
- Always `git commit -s`. The Husky `commit-msg` hook (`ui/.husky/`, from `make ui-setup`) is not a safeguard: a clone without it, or a tool that repoints `core.hooksPath`, bypasses it. To repair a pushed commit with no sign-off, see [Contributing](./docs/content/en/project/contributing/_index.md).

## Testing

### Go

- `go test --short ./...`; integration (needs Docker, kind, kubectl, helm): `make server-integration-tests-meshsync-setup`, then `make server-integration-tests-meshsync-run`. Target ≥70% coverage on business logic.
- Golden files (`-args -update`, `fixtures/` vs `testdata/`): [Contributing to Meshery CLI](./docs/content/en/project/contributing/cli/cli.md).
- **Never point viper at `mesheryctl/pkg/utils/TestConfig.yaml` in a test** (parallel packages kill each other via `Log.Fatal`). Use `utils.CopyMeshconfigFixture` or `utils.SetupCustomContextEnv` with your own testdata.
- A rename in `meshery/schemas` renames the gorm column too (derived from the Go field name, not `db:`). After a bump, grep every raw column reference (`Select`, `Where`, `Order`, `Joins`, `Scan`, migrations) and the `mesheryctl` fixtures for the old spelling, and propagate gorm errors.

### UI

- E2E (Playwright): `make ui-integration-tests` or `npm run test:e2e` in `ui/`; setup `make test-setup-ui`. Read [UI End-to-End Tests](./docs/content/en/project/contributing/ui/tests.md) before changing the workflow or setup projects.
- The job gates on `steps.playwright-tests.outcome`, never `.conclusion`. Never re-disarm it to turn a build green; fix the test or `test.fixme` it with the tracking issue.
- The credential is the `REMOTE_PROVIDER_TEST_USER_TOKEN` org secret, spelled identically in caller and reusable workflow. Never alias it.
- A missing credential must fail, never `setup.skip()`.
- A local run needs `make ui-provider-build`, a server on `:9081`, `make ui` on `:3000`, and `MESHERY_SERVER_URL=http://localhost:3000`.
- CI is serial (`workers: process.env.CI ? 1 : 4`); do not add workers.
- Test results feed the [meshery/qa](https://github.com/meshery/qa) Allure dashboard keyed on `testGroup`: CLI via the `[tg=<Test Group>]` `@test` token, UI via `allure.label('testGroup', …)`. Contract: `docs/content/en/project/contributing/{cli/tests.md,build-and-release.md}`.

Before committing: `make golangci` (Go), `make ui-lint` (UI).

## Security & Compliance

- Report vulnerabilities privately to [security@meshery.io](mailto:security@meshery.io), never in a public issue or PR. Reports are acknowledged within 10 business days.
- Never commit secrets; use env vars (`PROVIDER_BASE_URLS`, `KEYS_PATH`) and GitHub Secrets.
- Use parameterized queries; validate/sanitize all user input; Apache 2.0 — check dependency licenses.

## Agent Guardrails

- **Do not modify:** `LICENSE`, `CODE_OF_CONDUCT.md`, `GOVERNANCE.md`, `MAINTAINERS.md`, `.github/copilot-instructions.md`, `.github/agents/`.
- Never hand-edit a lock file (`go.sum`, any `package-lock.json`); regenerate it. `.agents/hooks/block-lockfiles.sh` enforces this by basename.
- **Require human review:** security changes (auth, secrets, encryption), database migrations, API breaking changes, Helm templates (`install/kubernetes/helm/`), CI/CD workflows (`.github/workflows/`).
- **Quality gates:** `make golangci` (runs `golangci-lint` and the repo analyzers, same as CI) and `make ui-lint` must pass; new features need docs, breaking changes need deprecation notices; keep PRs under 500 lines; don't merge on CI failure.

## Extensibility

- **Providers:** interface in `server/models/provider.go`.
- **Adapters (gRPC):** `server/meshes/meshops.proto`; adapters self-register.
- **UI extensions:** loaded via `@paciolan/remote-component`. A bundle MUST expose `module.exports = { default: Component, __esModule: true }` (`output.library.type = "commonjs2"`), and every event literal and capability key comes from `@sistent/sistent`'s `mesheryExtensionContract`, with `ui/utils/eventBus.ts` typed `EventBus<MesheryExtensionEvent>`. See [UI Extensions](./docs/content/en/project/contributing/ui/ui.md).
- **GraphQL:** `server/internal/graphql/schema.graphql`, then `make graphql-build`.
- **Config:** env vars `PLAYGROUND`, `DEBUG`, `SKIP_COMP_GEN`; runtime config `~/.meshery/config.yaml`. NATS topics: `meshsync.request`, `meshery.broker`. Build: extend `Makefile` or `install/Makefile.core.mk`.

## Agent Tooling

`.agents/` is the LLM-agnostic home for agent definitions, skills and hooks ([.agents/README.md](./.agents/README.md) — read it before touching any of them).

- `.agents/skills/` is the single source of truth for skills, one directory each. Add skills only there; never enumerate them here.
- `.claude/skills` is a symlink to `../.agents/skills` and nothing else. Never replace it with copies, and never run the AXI installer against a real checkout.
- Skill content addresses its files as `.agents/skills/...`, never through `.claude/`.
- Hooks in `.agents/hooks/`: `format-frontend.sh` (post-edit Prettier), `block-lockfiles.sh` (pre-edit lock-file guard).

## Worktrees and landing a change

- **Work in a treehouse worktree, never in the main clone.** Other agent sessions and the no-mistakes pipeline share the main clone. Take one with `treehouse get --lease -b <branch>` (`--base <branch>` to start elsewhere), release it with `treehouse return <path>` once merged; `treehouse status` lists the pool. Never use `git worktree add` or a second clone.
- **Land every change through no-mistakes, at most three runs per pull request.** Commit in the worktree, then `no-mistakes axi run --intent "<goal>"` (a clone without a `no-mistakes` remote needs `no-mistakes init` once). `no-mistakes axi status` shows the step and findings; `no-mistakes axi respond` answers a review gate. The pipeline reviews, tests, pushes and opens the pull request. Drive it with the vendored skill at `.agents/skills/no-mistakes/SKILL.md`. After the third run, open an issue for anything still flagged and hand the pull request to a human. A run that stopped on a usage limit counts; one cancelled before doing anything does not. Three runs do not waive the merge gate.
- **A pull request that changes any path outside documentation, the docs website, and CI merges only with a green no-mistakes run on its current head.** Documentation and the docs website are `docs/**` plus the root `CNAME`; CI is `.github/workflows/**`. The check is `PR must be raised via no-mistakes`: `review`, `test`, and `document` must each be `completed` in a head-bound attestation. Detail: [Build & Release (CI)](./docs/content/en/project/contributing/build-and-release.md#no-mistakes-merge-gate).

## Reusable Workflows Consumed by Other Repos

`.github/workflows/*.{yml,yaml}` are called from other `meshery` and `meshery-extensions` repos via `uses: meshery/meshery/.github/workflows/<file>@master`. Two traps that have shipped as silent no-ops:

- **`with:` values are not evaluated by a shell.** Use `${{ github.ref_name }}` / `${{ github.sha }}`, not `${GITHUB_REF/refs\/tags\//}` or a bare `${GITHUB_SHA}`; the shell form is correct inside `run:`, so don't "fix" those.
- **`if: github.repository == 'meshery/meshery'` disables the job for every external caller** (the `github` context is the caller's), and a skipped `needs:` skips dependents. Changing such a guard activates dormant deploys — a product decision, not a repair.

Find a workflow's callers before editing it: `gh-axi api -X GET search/code --field q='<workflow-file> path:.github/workflows'`. A caller with zero runs means CI has not tested your change; verify by reading. Detail: [Build & Release](./docs/content/en/project/contributing/build-and-release.md).

## Further Reading

The rules above are complete on their own. Open the file that matches your work:

| Working on | Read |
|---|---|
| A schema, an API contract, or a consumer of one | [Contributing to Meshery Schemas](./docs/content/en/project/contributing/contributing-schemas.md) |
| RTK Query, generated hooks, cache tags | [Schema-Driven UI Development](./docs/content/en/project/contributing/ui/schemas.md) |
| Playwright E2E, the E2E CI job, its credentials | [UI End-to-End Tests](./docs/content/en/project/contributing/ui/tests.md) |
| A build that fails for an unrelated-looking reason | [Build Environment Gotchas](./docs/content/en/project/contributing/contributing-build-environment.md) |
| MeshKit error codes | [How to write MeshKit compatible errors](./docs/content/en/project/contributing/contributing-error.md) |
| A Go lint rule firing, or adding one | [Go Lint Rules](./docs/content/en/project/contributing/contributing-lint.md) |
| Releases, CI secrets, the QA dashboard, reusable workflows, the no-mistakes merge gate | [Build & Release (CI)](./docs/content/en/project/contributing/build-and-release.md) |
| Connections and credential secrets | [Connections](./docs/content/en/project/contributing/models/connections.md) |
| A permission-gated page, control or key | [Extensibility: Authorization](./docs/content/en/reference/extensibility/authorization/index.md) |
| Providers, `PROVIDER` enforcement, boot-time seeding | [Extensibility: Providers](./docs/content/en/reference/extensibility/providers/index.md) |
| UI extensions, Remote Components | [Contributing to Meshery UI](./docs/content/en/project/contributing/ui/ui.md) |
| `mesheryctl`, golden files | [Contributing to Meshery CLI](./docs/content/en/project/contributing/cli/cli.md) |
| A docs page, its assets or shortcodes | [Contributing to Meshery Docs](./docs/content/en/project/contributing/contributing-docs/docs.md) |
| Commit hooks, repairing a missing sign-off | [Contributing](./docs/content/en/project/contributing/_index.md) |
| Agent definitions, skills, hooks | [.agents/README.md](./.agents/README.md) |

External: [Meshery Documentation](https://docs.meshery.io) · [Architecture](https://docs.meshery.io/concepts/architecture) · [Community Handbook](https://meshery.io/community#handbook) · [Contributing Guide](./CONTRIBUTING.md) · [Security Policy](./SECURITY.md) · [Governance](./GOVERNANCE.md)

## Maintaining this file

Keep this file for knowledge useful to almost every future agent session in this project.

**This file states rules; the files in Further Reading hold the detail.** A session must be able to learn that a constraint exists without opening anything else, so every MANDATORY rule and prohibition stays here, short and flat. Worked examples, causal history, run IDs, dates, evidence and per-endpoint walkthroughs belong in the linked file that owns the subject, with a pointer from here.

Do not repeat what the codebase already shows; point to the authoritative file or command instead. Relocate rather than delete, and prefer updating an existing doc over creating a new one. When you add a pointer, add it to the Further Reading table too.
