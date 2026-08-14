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

### Attaching local cache tags to a generated endpoint

`ui/rtk-query/index.ts` re-exports the schemas client itself (`mesheryApi as api`),
so every `ui/rtk-query/*` module injects into that same API instance - generated
hooks are therefore already available from those local modules and need no
re-declaration. To give a generated endpoint a local cache tag, use the **callback
form** of `enhanceEndpoints` via `appendInvalidatesTags` from `ui/rtk-query/utils`,
then re-export the generated hook - see the `importDesign` enhancement in
`ui/rtk-query/design.ts`. Do **not** use the object form
(`{ <operationId>: { invalidatesTags: [...] } }`): `enhanceEndpoints` applies an
object partial with `Object.assign(getEndpointDefinition(...) || {}, partial)`, so
it REPLACES `invalidatesTags` wholesale and every schemas-side tag has to be
hand-relisted - drift the moment schemas adds one. The callback form is handed the
live definition by reference, so the local tag is appended to the generated ones
and cannot drop them; that same lookup has no fallback, so `appendInvalidatesTags`
fails loudly and by name when the operationId is gone from schemas rather than
enhancing a throwaway object. Re-declaring the endpoint with `builder.mutation` to
get a tag is the forbidden path above - it forks the wire contract silently.

On the Go side, schemas ships **models only, no generated HTTP client**, so
`mesheryctl` builds request bodies from the generated structs (and, for `oneOf`
bodies, the `From<Variant>Payload` union builders) rather than a
`map[string]interface{}` — see `mesheryctl/internal/cli/root/design/import.go`.
A hand-written map is how camelCase `fileName` regressed to `file_name`.

### A local endpoint that shadows a schemas operationId is DEAD CODE

`injectEndpoints` without `overrideExisting: true` **silently discards** any
endpoint whose name `@meshery/schemas` already defines (dev-only console warning)
and serves every call from the schemas definition, so a local declaration can sit
there looking authoritative while a different request goes over the wire - which
is how notification delete shipped as `DELETE /api/events/undefined`. Before
adding a `builder.query`/`builder.mutation`, check the name against the generated
client (`grep '<name>:t\.' ui/node_modules/@meshery/schemas/dist/mesheryApi.js`);
if it is there, consume the generated hook - that is the rule above anyway.

The deliberate-override exception and the rule that tests must assert the
*effective* endpoint rather than the declared one are in
`docs/content/en/project/contributing/ui/schemas.md` (Integration Points in UI, A).

### Consumed contracts are the schemas type, not a copy of it

A struct Meshery only *decodes* from a remote provider (or only *encodes* to one)
carries no local freedom: it is the schemas construct, aliased. A local copy is
how `AnonymousFlowResponse` came to read `owner` while meshery-cloud kept sending
`userId`, so every anonymous sign-in wrote its capabilities under the nil UUID.
The same rename hit `PatternResource` and `Preference.selectedOrganizationId`.

When the Go type must stay local because it doubles as a GORM model (the schemas
models carry `db:` tags GORM does not read), keep the **JSON tags** identical to
the schemas construct and pin the column explicitly with `gorm:"column:..."` -
see `server/models/pattern_resource.go`. Guard it with a test that compares the
emitted JSON keys against the schemas type rather than restating them by hand.

### A credential's persisted `secret` has four shapes, and readers must tolerate all of them

Canonical (what `meshery/schemas`
`constructs/v1beta1/credential/forms/*.json` declares, and what Layer5 Cloud is
moving to) is a top-level `name` plus a `secret` object that **is** the payload.
Stored data also holds the Kubernetes `{auth, cluster}` shape, the legacy
double-nested `{credentialName, secret:{...}}` shape Meshery's credential form
still writes, and a legacy `{secret: "<token>"}` string. Legacy rows are never
rewritten, so tolerance - not migration - is what keeps them working.

`server/models/credential_secret.go` and its mirror `ui/utils/credentialSecret.ts`
own that decision for both languages; every read site delegates to them rather
than indexing the map. Reaching into `secret["secret"]` directly is how a
canonical credential's API key silently became an empty `Authorization` header.
The two files must keep the same *resolution rules* - their return types
deliberately differ on the legacy string shape. The shape catalogue and the
rules are documented once in
`docs/content/en/project/contributing/models/connections.md`.

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

A build that fails with `compile: version "goX.Y.Z" does not match go tool version`
on dozens of dependencies is an environment mismatch, not a code problem: a stale
`GOROOT` points at one Go installation while the `go` on `PATH` is a different
one. Any toolchain at or above `go.mod`'s version works - `go 1.26.4` there does
not mean 1.26.5 is wrong - so the fix is to make the two agree rather than to
pin an exact patch release. Drop the stale `GOROOT` and use one installation's
own binary, e.g.
`env -u GOROOT PATH="$HOME/.gvm/gos/go1.26.4/bin:$PATH" go test ...`.

### UI (Next.js/React)

```bash
make ui-setup              # Install dependencies
make ui                    # Dev server (port 3000)
make ui-build              # Build and export
make ui-lint               # Lint UI code
make ui-integration-tests  # Run E2E tests
```

`ui/tsconfig.tsbuildinfo` is a tracked build artifact that is not gitignored, so any local
`tsc --noEmit` leaves it modified and it has repeatedly been committed by accident. Stage
explicit paths rather than `git add -A`, and `git checkout -- ui/tsconfig.tsbuildinfo`
before committing.

When a change depends on an unreleased `@sistent/sistent` (or any sibling-repo package),
`ui/node_modules/@sistent/sistent` is often overwritten in place with a locally-built dist.
A local test run is then green against code that is not published, and CI fails on the same
commit. Re-verify with `npm ci` after any local sibling build before trusting a green run or
declaring a dependency bump done - a local build usually keeps the published version string,
so matching versions are not evidence that the installed contents are the published ones.
A version *mismatch* against `ui/package.json` is a useful tell that this has happened, but
a match proves nothing.

**Three package.json files consume `@sistent/sistent` - `ui/`, `provider-ui/` and
`install/docker-extension/ui/` - and a bump must cover all three.** Sistent's peers
(`@mui/x-date-pickers`, `date-fns`, the `@rjsf/*` set, `xstate`/`@xstate/react`) are not
optional: a consumer that omits one still installs cleanly and fails only at bundle time with
`Module not found` pointing *inside* `@sistent/sistent/dist`, which reads as a sistent bug
rather than a missing peer. `install/docker-extension/ui` installs with `--legacy-peer-deps`
(`@docker/docker-mui-theme` pins MUI <=6 against sistent's MUI 9) - see its Dockerfile.

`@meshery/schemas` deliberately keeps its `latest` dist-tag *below* its highest semver (1.4.0
is stale). That is safe because npm prefers the `latest`-tagged version whenever it satisfies
the range, so the `^1.3.x` carets do not jump to 1.4.0 - but verify the resolved version in
every regenerated lockfile rather than assuming it.

### CLI (mesheryctl)

```bash
cd mesheryctl && make                       # Build binary
cd mesheryctl && go test --short ./...      # Unit tests
cd mesheryctl && go test -run Integration ./...  # Integration tests
make docs-mesheryctl                        # Generate CLI docs
```

`make docs-mesheryctl` (i.e. `cd mesheryctl/doc && go run doc.go`) bakes the machine's
`$HOME` into every generated page's "Options inherited from parent commands" block (the
`--config` default path). Running it locally rewrites all ~100 pages under
`docs/content/en/reference/references/mesheryctl/` with your local home directory even
though only one command changed. CI/committed docs use `/home/runner/...` (the GitHub
Actions runner home). After regenerating, `git diff --stat` the docs dir, `git checkout --`
every file whose only change is that path, and manually fix the path back to
`/home/runner/...` in the pages you actually intended to change.

### Releasing

Meshery has **no automatic release cadence**. Release Drafter keeps exactly one draft
release current on every push to `master`; publishing that draft creates the `v*` tag,
which is what fires `build-and-release-stable.yml` and its fan-out. Follow
`.agents/skills/cut-release/SKILL.md` - never hand-author a tag or notes.

`gh release edit --draft=false` can exit 0 and leave the release a draft (seen cutting
v1.0.65). Publication is proven only by re-reading the release for `draft: false` plus a
non-null `published_at`, and by the release-triggered runs actually appearing.

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
  `make error` skips `mesheryctl` - a new `mesheryctl` code is taken from
  `mesheryctl/helpers/component_info.json` (`next_error_code`) and that value bumped in the
  same commit. `.github/workflows/error-codes-updater.yaml` re-runs errorutil and fails the
  PR if its analysis reports anything.
  The server side has the same contract in `server/helpers/component_info.json`: errorutil
  refuses to run at all ("next_error_code is lower than or equal to highest used code") until
  `next_error_code` is bumped past every code you added, so bump it in the same commit.
  Name each constant `<BuilderFuncName>Code` - errorutil keys the export off that pairing.
  `server/helpers/errorutil_errors_export.json` is gitignored, but the reference data at
  `docs/data/errorref/meshery-server_errors_export.json` is tracked: regenerate it with the
  `jq --slurpfile` wrapper the workflow uses, or the docs reference silently omits the new
  codes. Adding a constant longer than the block's current widest name makes gofmt realign
  the entire `error.go` const block - prefer a shorter name over a 300-line whitespace diff.
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

**A rename in `meshery/schemas` propagates further than the Go field name.** gorm
derives the AutoMigrate column from the *field name* via its naming strategy
(snake_case), not from the `db:` tag - only a `gorm:"column:..."` tag overrides
it. So renaming `UserID` to `Owner` renames the column `user_id` to `owner`
whatever the `db:` tag says, and any hand-written SQL naming the old one breaks -
silently, if the gorm error is dropped. After bumping schemas, grep every raw
column reference for the old spelling (`Select`, `Where`, `Order`, `Joins`,
`Scan`, and migrations - not just `Select`) along with the `mesheryctl` fixtures,
and propagate gorm errors so the next such rename fails loudly. Regression test:
`server/models/performance_profile_persister_test.go`.

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

Remote Components loaded via `@paciolan/remote-component`. Bundle **must** expose `module.exports = { default: Component, __esModule: true }`; a bundle built without `output.library.type = "commonjs2"` resolves to `undefined` with no loader error, so `NavigatorExtension` guards for it explicitly and reports the export shape as the cause. See `ui/components/layout/Navigator/NavigatorExtension.tsx`.

The host <-> extension contract (injected capability keys, event-bus event literals, contract version) is declared once in `@sistent/sistent`'s `mesheryExtensionContract` module and shared by both sides. Derive every event literal from `MESHERY_EXTENSION_EVENT` and every injected key from that module rather than typing strings: hand-duplicated literals are why `OPEN_DESIGN_IN_KANVAS` -> `OPEN_DESIGN_IN_EXTENSION` and `capabilitiesRegistry` -> `providerCapabilities` both shipped as silent runtime no-ops. `ui/utils/eventBus.ts` must stay typed as `EventBus<MesheryExtensionEvent>`; a bare `new EventBus()` widens `T` to its constraint and disables publish-site checking entirely. The `NavigatorExtension` unit test asserts the built `injectProps` bag against the contract, which is the gate that catches a capability rename before merge.

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

`.agents/skills/` is the single source of truth for every packaged workflow in this repo - one
directory per skill, each with a `SKILL.md`. Do not enumerate them here; list the directory. Add a
new skill only there.

Per-tool discovery, so no skill is ever copied per tool:

| Tool | How it finds these skills |
|---|---|
| Codex | Natively scans `$REPO_ROOT/.agents/skills` - nothing to configure ([docs](https://learn.chatgpt.com/docs/build-skills)) |
| OpenCode | Natively scans `.agents/skills`, one of six roots it searches alongside `.opencode/skills` and `.claude/skills` ([docs](https://opencode.ai/docs/skills/)) |
| Claude Code | Reads `.claude/skills`, which is a relative symlink to `../.agents/skills` |

`.claude/skills` is that symlink and nothing else. Never replace it with real directories or copies -
that reintroduces the drift this layout removes.

Skill content must address its own files by their canonical `.agents/skills/...` path, never
through `.claude/`. `iterate-pr` used to invoke its scripts as
`.claude/skills/iterate-pr/scripts/<script>.py`, which resolved only through the symlink and so
broke wherever the symlink was absent; it was corrected to `.agents/` in iterate-pr 2.4.0. The
symlink is therefore a *discovery* path for Claude Code, not a runtime dependency - but it is
still load-bearing for discovery, and the installer-collision hazard below is a further reason
not to touch it.

The four skills tracked in `skills-lock.json` - `chrome-devtools-axi`, `gh-axi`, `lavish`,
`quota-axi` - are installed by the AXI installer, and its layout is skill content at
`.agents/skills/<name>/` *plus* a per-skill symlink at `.claude/skills/<name>`. Those per-skill
symlinks were installer-owned, not hand-made, and this layout removed them as redundant. The next
installer run recreates them, and `.claude/skills/<name>` now resolves *through* the directory
symlink onto `.agents/skills/<name>` - an existing real directory holding the canonical content.
Best case the installer fails with `EEXIST`. Worst case a force-replacing installer destroys that
canonical directory and leaves a self-referential symlink loop. Which of the two occurs is not
established, and must not be determined by running the installer against a real checkout: the
failure mode under test is destruction of the canonical skill content.

Neither `.codex/skills` nor `.opencode/skills` is created: both tools already read `.agents/skills`
natively, so a second copy or link would be redundant. `.opencode/skills` is a real OpenCode search
root, just an unnecessary one here; `.codex/skills` is not a path Codex scans at all.

Windows caveat: on a checkout with `core.symlinks=false` - the default outside developer mode - git
materialises `.claude/skills` as a regular text file containing the literal string
`../.agents/skills`. Claude Code then discovers no project skills. Enable Windows developer mode
or set `git config core.symlinks true`, then re-checkout. (Skill *scripts* still resolve there,
because skill content addresses them via `.agents/` - only discovery breaks.)

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

## Maintaining this file

Keep this file for knowledge useful to almost every future agent session in this project.
Do not repeat what the codebase already shows; point to the authoritative file or command instead.
Prefer rewriting or pruning existing entries over appending new ones.
When updating this file, preserve this bar for all agents and keep entries concise.
