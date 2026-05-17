# Design: Eliminate residual `@mui/*` imports

**Date:** 2026-05-16
**Status:** Approved
**Author:** hamza-mohd
**PR shape:** Single PR, single commit (per user direction)

---

## 1. Goal

Bring the repo to a state where every `@mui/*` import is either:

1. A consciously-allowlisted "approved wrapper" that exists because Sistent does not yet expose the underlying primitive, OR
2. A consciously-allowlisted "documented exception" in `install/docker-extension/ui/`, where the pinned `@sistent/sistent@0.16.10` lacks `ButtonBase` and `HelpIcon`.

…and where `audit-mui.js` *actually sees* every place an `@mui/*` import can live, so the same regression class cannot hide again.

Concretely: drop the repo-wide `from '@mui/...'` file count from **15 → 7**, where every one of the remaining 7 is intentional and allowlisted with a written reason. `from '@rjsf/mui'` is unchanged at 1 (the existing `RJSFProvider.tsx` wrapper).

## 2. Background

The task as originally framed described 100 files still importing from `@mui/*`. A repo-wide grep for `from '@mui/...'` finds **15**. Five of those are not real work to do:

- 4 `ui/` wrappers that are kept by design (`theme/index.ts`, `TreeView.tsx`, `DatePicker/index.ts`, `MesheryDateTimePicker.tsx`) — Sistent does not yet upstream `GlobalStyles`/`darken`, `@mui/x-tree-view`, or `@mui/x-date-pickers`.
- 1 audit script itself (`ui/scripts/audit-mui.js`) — the `@mui/...` strings live in its regex, not in real imports.

(`RJSFProvider.tsx`, the 5th `ui/` wrapper, doesn't appear in the `@mui/` grep — it imports from `@rjsf/mui`. It's allowlisted separately.)

The remaining 10 files are the surprises that are genuinely worth migrating:

- **4 drawer-icon hover components** in `ui/public/static/img/drawer-icons/` — these are React components, imported by [navigatorComponents.tsx:10-12](../../../ui/components/layout/Navigator/navigatorComponents.tsx). They were invisible to `audit-mui.js` because the script's `SCAN_DIRS` list does not include `public/`. One of the four (`conformance_hover_svg.js`) is dead code: defined but never imported.
- **6 files in `install/docker-extension/ui/src/components/`** — a separate sub-app (own `package.json`, own React-Scripts build, runs embedded in Docker Desktop). The audit script never scanned this sub-app at all.

The two facts above mean the recurring "Eliminate MUI" task has been mis-scoped by missing instrumentation, not by missing work. Closing the instrumentation gap is the long-term-maintainer move that prevents this from recurring.

## 3. Decisions (locked during brainstorm)

| Question | Decision |
|---|---|
| Scope | Drawer icons + `install/docker-extension/ui/` (6 files). The 5 approved `ui/` wrappers stay; the 7 `@sistent/mui-datatables` consumers are acceptable per task brief. |
| Docker theme | Keep `DockerMuiThemeProvider` — Docker Desktop native theme integration is a feature, not collateral damage. `@docker/docker-mui-theme` peers MUI; this implies `@mui/material` and `@mui/icons-material` stay in `install/docker-extension/ui/package.json` regardless. |
| Sistent 0.16.10 gaps (`ButtonBase`, `HelpIcon`) | Leave as documented exceptions. Adding to `install/docker-extension/` allowlist with inline reason. |
| PR shape | Single PR, single commit. |
| Audit hardening | In scope. Closes the regression hole that hid this work. |

## 4. Changes

### 4.1 `ui/` — drawer icons

Three files: replace `@mui/icons-material/ArrowDropDown` import with the Sistent re-export.

- [ui/public/static/img/drawer-icons/configuration_hover_svg.js](../../../ui/public/static/img/drawer-icons/configuration_hover_svg.js)
- [ui/public/static/img/drawer-icons/lifecycle_hover_svg.js](../../../ui/public/static/img/drawer-icons/lifecycle_hover_svg.js)
- [ui/public/static/img/drawer-icons/performance_hover_svg.js](../../../ui/public/static/img/drawer-icons/performance_hover_svg.js)

Each goes from:

```js
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
```

to:

```js
import { ArrowDropDownIcon } from '@sistent/sistent';
```

Verified: `@sistent/sistent@0.21.11` (main UI's pin) exports `ArrowDropDownIcon` and it is already used at [ActionButton.tsx:3](../../../ui/components/designs/patterns/ActionButton.tsx).

Component output and prop interface are unchanged.

### 4.2 `ui/` — dead-code removal

Delete [ui/public/static/img/drawer-icons/conformance_hover_svg.js](../../../ui/public/static/img/drawer-icons/conformance_hover_svg.js). `ConformanceHover` is exported but never imported (verified via repo-wide grep). The "Conformance" navigator entry was removed earlier; the SVG component lingered.

### 4.3 `install/docker-extension/ui/` — primitive migration

Symbol-by-symbol swap to `@sistent/sistent`, verified against the Sistent 0.16.10 dist tarball.

| File | Imports changed | Imports kept (exception) |
|---|---|---|
| [styledComponents.js](../../../install/docker-extension/ui/src/components/ExtensionComponent/styledComponents.js) | `styled`, `Typography` → `@sistent/sistent` | `ButtonBase` from `@mui/material` (Sistent 0.16.10 has no equivalent) |
| [Catalog/style.js](../../../install/docker-extension/ui/src/components/Catalog/style.js) | `styled`, `Box` → `@sistent/sistent` | — |
| [Walkthrough/tourStyledComponents.js](../../../install/docker-extension/ui/src/components/Walkthrough/tourStyledComponents.js) | `styled` → `@sistent/sistent` | — |
| [ExtensionComponent/ExtensionComponent.js](../../../install/docker-extension/ui/src/components/ExtensionComponent/ExtensionComponent.js) | `Typography`, `Button`, `Tooltip`, `Grid`, `Avatar`, `OpenInNewIcon`, `CssBaseline` → `@sistent/sistent` | `DockerMuiThemeProvider` from `@docker/docker-mui-theme` (intentional — Docker Desktop theme integration) |
| [Catalog/CatalogCard.js](../../../install/docker-extension/ui/src/components/Catalog/CatalogCard.js) | `Tooltip` → `@sistent/sistent` | — |
| [Walkthrough/Tour.js](../../../install/docker-extension/ui/src/components/Walkthrough/Tour.js) | `IconButton`, `Tooltip` → `@sistent/sistent` | `HelpIcon` from `@mui/icons-material/Help` (Sistent 0.16.10 has no `HelpIcon`, only `HelpOutlinedIcon` — visual divergence not acceptable) |

**Prop-shape note:** Sistent's `<Button>` accepts both `label` and `children`. All four `<Button>` usages in `ExtensionComponent.js` pass children; no signature change needed.

### 4.4 `install/docker-extension/ui/package.json` — drop unused dep

Remove `"@mui/styles": "^6.4.8"`. Verified zero source imports of `@mui/styles` (no `makeStyles` / `withStyles` calls).

`@mui/material`, `@mui/icons-material`, `@mui/styled-engine-sc` stay — peers of `@docker/docker-mui-theme` and required by the two documented exceptions above.

### 4.5 `ui/scripts/audit-mui.js` — audit hardening

The audit currently misses two entire categories of imports. After this PR it sees them. Specifically:

1. **Multi-root scan.** Refactor the script's single `(UI_ROOT, SCAN_DIRS, APPROVED_WRAPPERS)` triple into a `SCANS` array of `{name, root, dirs, allowlist}` entries. Two entries:
   - `ui` — same as today, plus `public/` added to `dirs`.
   - `docker-extension` — root at `install/docker-extension/ui/`, dirs `[src]`, allowlist:
     ```
     src/components/ExtensionComponent/styledComponents.js  // ButtonBase (Sistent 0.16.10 gap)
     src/components/Walkthrough/Tour.js                     // HelpIcon (Sistent 0.16.10 gap)
     ```
2. **`--fail-on-new` flag.** Exit code 1 if any non-allowlisted file matches the import pattern. Default behavior unchanged (exit 0, informational). The flag is opt-in for now; wiring it into CI is a follow-up once the count is stable.
3. **Output preserved.** The final `AUDIT mui files=<N> matches=<M>` line stays — that's the CI-grep contract. The per-scan breakdown lands in the human-readable section above it.

The header comment block is updated to reflect the new scope and the new flag.

### 4.6 Commit

Single signed commit (`git commit -s`):

```
[UI] Phase 6: Eliminate residual @mui/* imports from ui/public and docker-extension

- ui/: migrate 3 drawer-icon hovers (configuration, lifecycle, performance)
  to @sistent/sistent ArrowDropDownIcon. Delete dead conformance_hover_svg.js.
- install/docker-extension/ui/: migrate 6 files' primitives to @sistent/sistent.
  Keep DockerMuiThemeProvider (Docker Desktop theme integration). Drop unused
  @mui/styles from package.json. Two narrow MUI imports remain as documented
  exceptions (ButtonBase, HelpIcon — Sistent 0.16.10 lacks equivalents).
- ui/scripts/audit-mui.js: scan ui/public/ and install/docker-extension/ui/src
  with separate allowlists. Add --fail-on-new flag (off by default).
```

`-s` appends the `Signed-off-by:` trailer automatically from the git config.

## 5. Verification

| Check | Command | Expected |
|---|---|---|
| Audit (informational) | `node ui/scripts/audit-mui.js` | `AUDIT mui files=0 matches=0` (with allowlist accounting for 5 ui wrappers + 2 docker-extension exceptions) |
| Audit (gate) | `node ui/scripts/audit-mui.js --fail-on-new` | Exit 0 |
| Lint | `cd ui && npm run lint` | Green (no new no-restricted-imports violations) |
| Unit tests | `cd ui && npm test` | Green. Drawer-icon mocks at [navigatorComponents.test.tsx:30-36](../../../ui/components/layout/Navigator/navigatorComponents.test.tsx) are transparent to import path; tests should pass unchanged. |
| Build (main UI) | `cd ui && npm run build` | Green |
| Build (docker-extension) | `cd install/docker-extension/ui && npm install --legacy-peer-deps && npm run build` | Green (matches Dockerfile `client-builder` stage) |
| Repo-wide grep (`@mui/`) | `grep -rln "from ['\"]@mui/" --include="*.{ts,tsx,js,jsx}" --exclude-dir=node_modules .` | Drop from 15 to **7**: 4 ui/ approved wrappers (`theme/index.ts`, `TreeView.tsx`, `DatePicker/index.ts`, `MesheryDateTimePicker.tsx`), 1 audit script (regex string literals, not real imports), 2 docker-extension allowlisted exceptions (`styledComponents.js`, `Tour.js`). |
| Repo-wide grep (`@rjsf/mui`) | `grep -rln "from ['\"]@rjsf/mui" --include="*.{ts,tsx,js,jsx}" --exclude-dir=node_modules .` | Unchanged at **1**: `ui/components/shared/FormFields/RJSFProvider.tsx` (5th approved wrapper). |
| Sistent re-export sanity | Manual: render the navigator drawer at the three hover states; render the docker-extension `ExtensionComponent`, tour, and catalog cards. | No visual regression. |

**Manual visual verification of docker-extension is best-effort only.** The sub-app only runs fully inside Docker Desktop with the extension installed. Build success + symbol-for-symbol swap (no signature changes) is the realistic floor of pre-merge verification. Risk is mitigated by: (a) keeping diffs surgical; (b) the existing nested `SistentThemeProviderWithoutBaseLine` already proves Sistent components render correctly inside the Docker MUI theme tree.

## 6. Risks

- **Sistent's `<Button>` is opinionated.** It accepts `label` and `children`. All current MUI `<Button>` usages in `ExtensionComponent.js` pass children. The component renders an internal `<MuiButton>` so behavior is broadly compatible. Verified mitigation: spot-check rendering during build.
- **`Grid` API drift.** Sistent 0.16.10 re-exports MUI v6 `Grid`. The docker-extension `package.json` also pins `@mui/material@^7.3.7` directly. Today the lockfile resolves to a single MUI version under `--legacy-peer-deps`; this PR does not change that resolution. `npm ls @mui/material` post-install should match the pre-PR state.
- **Audit script refactor scope creep.** The multi-root rewrite is a non-trivial diff. Mitigation: keep the externally-observable behavior (final `AUDIT mui files=N matches=M` line) byte-for-byte compatible so existing CI grep continues to work.

## 7. Explicitly out of scope (documented for follow-up)

These are visible but not addressed in this PR. They're called out so they're not lost:

1. **`@sistent/sistent` version mismatch in docker-extension.** Pinned at `0.16.10` (peers MUI ^6) while the same `package.json` directly pins `@mui/material@^7.3.7`. The version skew is papered over by `--legacy-peer-deps`. Resolving this requires either bumping Sistent to a release that peers MUI ^7, or pinning MUI back to ^6. Either way it's a deps-graph change with its own risk surface; deferred to a follow-up.
2. **The two documented MUI exceptions (`ButtonBase`, `HelpIcon`).** Eliminable by the same Sistent bump above. Listed in the docker-extension allowlist with explicit `// Sistent 0.16.10 gap` comments so the next maintainer knows the exit condition.
3. **Wiring `audit-mui.js --fail-on-new` into CI.** Once the post-PR audit is verified stable at the allowlist baseline, a follow-up PR can add the gate to the `.github/workflows/` lint/test job.
4. **`ui/`'s 5 approved wrappers.** Sistent may have closed some of the gaps that originally justified the wrappers (`GlobalStyles`, `darken`, date pickers, tree view, RJSF theme adapter). Worth a periodic re-audit — but separately, against the current Sistent release, not as part of this PR.

## 8. Spec self-review (post-write)

- **Placeholders:** none — every file path and symbol verified by direct read or grep.
- **Internal consistency:** §4.1–§4.5 changes map one-to-one to the §5 verification table. Pre→post `@mui/` grep math: 15 (today) − 3 migrated drawer icons − 1 deleted (`conformance_hover_svg.js`) − 4 migrated docker-extension files (`Catalog/style.js`, `tourStyledComponents.js`, `ExtensionComponent.js`, `CatalogCard.js`) = **7** remaining, all allowlisted.
- **Scope:** single implementation plan, fits a single commit. No decomposition needed.
- **Ambiguity:** the only judgment call (which symbols stay as MUI exceptions) is named explicitly in §4.3 with a per-symbol reason.
