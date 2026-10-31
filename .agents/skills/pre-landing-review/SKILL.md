---
name: pre-landing-review
description: Two-pass code quality and safety review gating pull requests on Meshery failure classes and AGENTS.md rules.
---

# Pre-Landing Review

`pre-landing-review` is a two-pass gating review skill for Meshery pull requests. It inspects staged or proposed changes against known failure classes, historical CVE patterns, and architectural rules documented in `AGENTS.md`.

## Workflow Overview

The review executes in two sequential passes:
1. **Pass 1 — Blocking CRITICAL Checks**: Architecture and security requirements that block PR landing if violated.
2. **Pass 2 — Advisory INFORMATIONAL Checks**: Code hygiene, documentation, and maintainability checks.
3. **Automated Verification**: Runs real repo linters (`make golangci`, `make ui-lint`) and reports actual output.

---

## Pass 1: Blocking CRITICAL Checks

Every item in this pass must be verified before approval. If any check fails, the change is blocked until remediated.

### 1. GORM SQL Injection via `.Order(...)`
* **Rule**: Dynamic strings must never be passed directly to GORM's `.Order()` clause.
* **Remediation**: Always pass user-supplied order parameters through `models.SanitizeOrderInput(order, allowedColumns)` using an explicit allow-list of `snake_case` database column names, or use compile-time constants.
* **Reference**: Enforced by the custom analyzer `server/internal/lint/orderby`.

### 2. Provider Resolution at Boot
* **Rule**: Boot-time or background tasks must never resolve a provider by name directly out of `HandlerConfig.Providers` (e.g., `Providers[LocalProviderName]`), as enforced provider deployments drop unselected registrations.
* **Remediation**:
  * Use `HandlerConfig.SystemEventPersister` for system-level events raised outside user request contexts.
  * Index `Providers` only on request paths using the comma-ok idiom (`provider, ok := providers[name]`).
  * Wrap boot-time operations in `models.RunSeedStage` with deferred `recover()` handlers in spawned goroutines.

### 3. UI Access Gating and Query Suppression
* **Rule**: Content-bearing pages must enforce permission checks on page components and halt background data queries when access is denied.
* **Remediation**:
  * Check capabilities via `const canAccess = useHasPermission(Keys.<PermissionKey>)` from `@sistent/sistent`.
  * Render `<DefaultError permissionKey={Keys.<PermissionKey>} />` when permission is false.
  * Pass `{ skip: !canAccess }` to RTK Query hooks so denied sessions do not issue queries that fail with 403/401.

### 4. Schema-First API Discipline
* **Rule**: Any new or changed HTTP API (endpoints, query parameters, request/response models) must be defined in `meshery/schemas`. Hand-rolled endpoints in `ui/rtk-query/*` or ad-hoc `axios`/`fetch` calls are forbidden.
* **Remediation**: Define the OpenAPI contract in `meshery/schemas`, generate client bindings, and consume generated hooks.

### 5. Multi-Layer Identifier Naming
* **Rule**: The ORM is the sole translation boundary:
  * **Database columns**: `snake_case` (`db:"user_id"`)
  * **Wire JSON properties**: `camelCase` (`json:"userId"`)
  * **Go struct fields**: `PascalCase` (`UserID`, `OrgID`)
  * **TypeScript properties**: `camelCase` (`userId`)
* **Remediation**: Never use `ID` (all caps) in JSON tags or URL params. Never match `json:` tag directly to `db:` tag.

### 6. MeshKit Structured Errors
* **Rule**: Handlers must return structured MeshKit errors (`github.com/meshery/meshkit/errors`) containing valid error code, summary, cause, and actionable remediation.
* **Remediation**: Use `writeMeshkitError(w, err, statusCode)`. When adding new error codes, increment `next_error_code` in `helpers/component_info.json` in the same commit and export reference data to `docs/data/errorref/`.

### 7. DCO Commit Sign-off
* **Rule**: Every commit must be signed-off under the Developer Certificate of Origin (`git commit -s`).

---

## Pass 2: Advisory INFORMATIONAL Checks

These checks verify repository hygiene and long-term maintainability:

1. **No Generated Build Artifacts in Git**: Ensure `ui/tsconfig.tsbuildinfo` is not staged or modified.
2. **Credential Secret Isolation**: Do not index credential secret maps with raw string keys; resolve shapes through `server/models/credential_secret.go` or `ui/utils/credentialSecret.ts`.
3. **Regression Tests**: Any bug fix must include a unit test in a corresponding `*_test.go` file replicating the defect.
4. **Docs Sync**: User-facing changes must update associated documentation in `docs/content/en/`.

---

## Automated Verification

Execute the following commands and include results in the review output:

```bash
# Go server & CLI linting
make golangci

# UI linting
make ui-lint
```
