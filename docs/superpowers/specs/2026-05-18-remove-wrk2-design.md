# Design: Complete removal of wrk2 load-generator support

- **Date:** 2026-05-18
- **Branch:** `remove-wrk2`
- **Status:** Approved (design), pending spec review

## 1. Context & Goal

Meshery's performance subsystem exposes a `LoadGenerator` abstraction with two
options: **fortio** and **wrk2**. The goal is to remove `wrk2` support
**completely** — code paths, dependency, git submodule, build wiring, UI
options, and documentation — while keeping fortio fully functional.

We are removing the **wrk2 option**, not the load-generator abstraction.
`FortioLG`, the `LoadGenerator` type, the `--load-generator` flag, the
`LoadGenerators` profile field, and the UI selector all remain (fortio-only).

## 2. Decisions (resolved during brainstorming)

1. **Legacy data → graceful fallback.** Existing performance profiles, saved
   user preferences, and in-flight API requests that still reference `"wrk2"`
   must keep working with **no user-visible error**. They transparently run on
   fortio.
2. **History preserved.** Release notes (`v0.3.1`, `v0.6.x`, `v0.8.56`, …) and
   mesheryctl mock-response / output golden fixtures are an immutable changelog
   and historical test data. They are **not** edited. Only the active e2e bats
   test that invokes `--load-generator wrk2` is removed. Golden files are
   regenerated **only if** CLI output genuinely changes (it will not — display
   logic is untouched).
3. **openapi.yml → flag only.** `docs/data/openapi.yml` is generated/synced from
   the `meshery/schemas` repo. The 4 `"e.g. fortio, wrk2"` description strings
   are **not** edited here; instead the PR description flags that
   `meshery/schemas` needs the upstream wording change (the `LoadGenerators`
   schema field itself is unchanged).

## 3. Architecture / Approach

With fortio as the **sole** load generator, the load-generator value becomes
vestigial on the execution path: every test runs on fortio regardless of the
incoming value. Graceful fallback therefore falls out of the design for free —
there is no code path left that can reject or mishandle a legacy `wrk2` value.

Implementation principle: **inline `FortioLG` and delete the now-vacuous
branching/validation**, with explanatory comments. We deliberately do **not**
introduce a `wrk2 → fortio` mapping helper: a function that ignores its
argument is less honest than a direct assignment plus a comment stating that
legacy/unknown generators run on fortio.

The UI already coerces unknown generator values to `fortio` via
`toLoadGenerator` in `load-test-prefs.ts`; tightening `validLoadGenerators` to
`{fortio}` makes that path downgrade legacy `wrk2` automatically.

## 4. Change Map

### 4.1 Go server — functional removal + graceful fallback

- **`server/models/load-test.go`** — delete the `Wrk2LG` constant. Keep
  `FortioLG`, the `LoadGenerator` type, and `Name()`.
- **`server/helpers/load_test_interface.go`** — delete the `WRK2LoadTest`
  function and the `github.com/layer5io/gowrk2/api` import. Leave
  `FortioLoadTest`, `sharedHTTPOptions`, `HTTPRunnerResults` untouched.
- **`server/handlers/load_test_handler.go`** — collapse all three wrk2
  `switch` blocks:
  - the two input switches (≈ lines 121, 320) become
    `loadTestOptions.LoadGenerator = models.FortioLG` with a comment that
    legacy/unknown generator values (e.g. removed `wrk2`) run on fortio;
  - the execute switch (≈ line 421) becomes a direct
    `helpers.FortioLoadTest(loadTestOptions, h.log)` call.
- **`server/models/validators.go`** — the generator allow-list is now vacuous
  (one option, value ignored at runtime). Remove the `validGenerators` /
  `isValidGenerator` / `ErrLoadgenerator` check with a comment that the
  load-generator field is legacy-tolerant (any value runs on fortio). If
  `ErrLoadgenerator` becomes unreferenced repo-wide, delete its definition;
  otherwise leave it.
- **`server/handlers/user_handler.go`** — remove the
  `[]models.LoadGenerator{models.FortioLG, models.Wrk2LG}` rejection loop and
  its `ErrSavingUserPreference("invalid load generator")` branch. Saving a
  preference with a legacy `gen: "wrk2"` must succeed (runs on fortio).
- **`server/models/validators_test.go`** — update the test case that uses
  `Wrk2LG.Name()` so it reflects fortio-only behavior and asserts that a
  profile carrying a legacy `wrk2` generator still validates (graceful
  fallback regression guard).

### 4.2 Build & dependencies

- **Submodule `server/cmd/gowrk2`** — `git rm` the submodule, delete the
  `[submodule "server/cmd/gowrk2"]` block from `.gitmodules`, remove the
  `submodule.server/cmd/gowrk2.*` keys from `.git/config`, and remove
  `.git/modules/server/cmd/gowrk2` if present.
- **`go.mod` / `go.sum`** — remove `github.com/layer5io/gowrk2 v0.6.1` and its
  now-orphaned transitive entries via `go mod tidy` (run **after** the import
  is removed, else tidy re-adds it).
- **`Makefile`** — delete the `wrk2-setup` target and its `.PHONY` entry.
  Verify no other target depends on `wrk2-setup`; if one does, remove that
  dependency edge too.
- **`install/docker/testing/Dockerfile`** — remove the active
  `FROM base-alpine AS wrk2` stage (the `git clone … wrk2 && make`) and the two
  `COPY --from=wrk2 …` lines.
- **`install/docker/Dockerfile`** — remove the dead commented-out wrk2 block
  (`# FROM ubuntu AS wrk2`, `# RUN git clone … wrk2`, `# COPY --from=wrk2 …`)
  for consistency with a complete removal.

### 4.3 mesheryctl

- **`mesheryctl/internal/cli/root/perf/apply.go`** — keep the
  `--load-generator` flag (fortio is still a valid value). Change the usage
  string `(fortio/wrk2)` → `(fortio)`; remove the wrk2 lines from the long
  help / examples (the `--load-generator wrk2` example, the
  `gowrk2.go` options link, and the wrk2 options example). Default stays
  `fortio`.
- **`mesheryctl/tests/e2e/004-perf/00-perf-apply.bats`** — remove the
  `@test "mesheryctl perf apply uses specified load generator wrk2"` case.
- **Golden files** — **unchanged.** Display logic is untouched, so historical
  `wrk2` rows/fields in `testdata/*.output.golden` and
  `fixtures/*.response.golden` still render correctly. The removed bats test
  asserts on output partials (no golden comparison), so removing it orphans
  nothing.

### 4.4 UI

- **`ui/components/performance/performance-helpers.tsx`** — `loadGenerators`
  → `['fortio']`; remove the wrk2 `<li>` and the Coordinated-Omission / wrk2
  sentence from the `infoloadGenerators` JSX.
- **`ui/components/settings/MesherySettingsPerformanceComponent.tsx`** —
  `loadGenerators` → `['fortio']`; remove the dead `disabled={lg === 'wrk2'}`
  branch.
- **`ui/components/performance/PerformanceForm.tsx`** — remove the dead
  `disabled={lg === 'wrk2'}` branch; render fortio only.
- **`ui/lib/load-test-prefs.ts`** — `validLoadGenerators` →
  `new Set(['fortio'])`. The existing `toLoadGenerator` coercion then maps a
  stored legacy `gen: 'wrk2'` to the `fortio` default — the UI-side graceful
  fallback.
- **Tests — updated to verify coercion, not assert wrk2:**
  - `ui/lib/__tests__/load-test-prefs.test.ts` — change the `gen: 'wrk2'`
    round-trip cases to assert `wrk2` → `fortio` coercion.
  - `ui/components/settings/MesherySettingsPerformanceComponent.test.tsx` —
    the `renders all load generators as radio options` test asserts only
    `fortio` is rendered and `wrk2` is **not**.
  - `ui/rtk-query/__tests__/user.test.ts` — the `gen: 'wrk2'` request/response
    cases assert the normalized `fortio` outcome.
  - `ui/store/slices/__tests__/prefTest.test.ts` — the `gen: 'wrk2'` slice
    cases assert the coerced `fortio` result.

### 4.5 Docs

- **Edited (prose):** `README.md` (line ≈199), `docs/content/en/reference/
  extensibility/load-generators/index.md` (list item + Coordinated-Omission
  paragraph + "and wrk2" phrasing), `docs/content/en/guides/
  performance-management/managing-performance.md` (remove the `### wrk2`
  section **and** its `- [wrk2](#wrk2)` TOC link),
  `docs/content/en/concepts/architecture/adapters/index.md`
  ("using fortio and wrk2" → "using fortio").
- **Edited (generated, kept in sync with the new CLI help):**
  `docs/content/en/reference/mesheryctl/perf/apply.md`,
  `docs/data/mesheryctlcommands/cmds.yml`,
  `docs/data/mesheryctlcommands/es_cmds.yml` (drop wrk2 from description,
  usage, and example strings — including the Spanish translation).
- **Not touched:** release notes (immutable changelog),
  `docs/data/openapi.yml` (decision 3 — flag upstream instead),
  `docs/data/catalog/**/design.yml` wrk2 hits (base64 image blobs — not real
  references), `ui/package-lock.json` &
  `install/docker-extension/ui/package-lock.json` (coincidental substrings).

## 5. Testing Strategy

- **Go:** `go build ./...`; `go vet ./...`; `go mod tidy` then assert
  `git diff --exit-code go.mod go.sum` is clean; targeted
  `go test ./server/models/... ./server/handlers/...`; assert
  `grep -r "gowrk2" --include=*.go .` (excluding vendor) returns nothing.
- **mesheryctl:** `go test ./mesheryctl/internal/cli/root/perf/...` — golden
  tests pass **unchanged**; remaining perf bats unaffected by removing the
  wrk2 case.
- **UI:** jest on the 4 touched suites + project typecheck/lint; new
  assertions positively prove `wrk2 → fortio` coercion on every read path.
- **Final grep gate:** `grep -ri "wrk2"` across source (excluding `.git`,
  `node_modules`, `vendor`) returns **only** the documented allow-list:
  release notes, golden fixtures, base64 catalog blobs, package-lock
  coincidental hits, and the deliberate legacy-compat code comments. Any
  other hit is a defect.

## 6. Cross-Repo Flags (for PR description)

- **`meshery/schemas`** — the `LoadGenerators` field description example text
  (`"… e.g. fortio, wrk2 …"`, 4 occurrences in the generated
  `docs/data/openapi.yml`) should drop `wrk2` upstream; otherwise it
  regenerates. The schema field itself is unchanged.
- **`layer5io/gowrk2`, `layer5io/wrk2`** — now fully unreferenced by Meshery
  (informational; archival is out of scope).
- **Kanvas / meshery-cloud** — verify no UI there surfaces a `wrk2`
  performance-profile choice that would now silently map to fortio
  (informational).

## 7. Acceptance Criteria

1. No `wrk2`/`gowrk2` references remain in Go, TS/TSX/JS source, Makefile,
   Dockerfiles, `.gitmodules`, `go.mod`/`go.sum`, or edited docs — except the
   documented legacy-compat comments.
2. `server/cmd/gowrk2` submodule fully removed (working tree, `.gitmodules`,
   `.git/config`).
3. Existing performance profile / user preference / API request carrying
   `gen|loadGenerator == "wrk2"` runs successfully on fortio with no error
   (verified by Go + UI tests).
4. fortio performance tests function unchanged end-to-end.
5. All listed test suites green; `go mod tidy` and golden files produce no
   unexpected diff.
6. PR description carries the section 6 cross-repo flags.
