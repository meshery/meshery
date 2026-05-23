# **Meshery UI Restructure Plan**

**Branch:** `claude/ui-restructure-plan-JuSmU` **Scope:** `meshery/meshery` — `/ui` **Goal:** Collapse the sprawl. Enthrone `@sistent/sistent` as the single source of truth for every component, token, color, and theme. Eliminate Material UI. Break the giant files. Make the directory tree discoverable at a glance.

---

## **0\. Executive Summary**

The Meshery UI is fundamentally sound — 87% TypeScript, RTK Query for data, Sistent already imported in 276 files — but it has accreted three overlapping styling systems, two parallel icon libraries, two parallel workspace concepts, and eight files over 1,000 lines. The task is not a rewrite; it is a **disciplined consolidation** across five phases:

| Phase                       | Theme                                                                | Effort | Risk                        |
| --------------------------- | -------------------------------------------------------------------- | ------ | --------------------------- |
| **P1** Freeze the sprawl    | ESLint rules \+ docs that make the target state _enforceable_        | Small  | None                        |
| **P2** Single design system | Eliminate all 100 `@mui/*` imports → `@sistent/sistent`              | Medium | Medium — visual regressions |
| **P3** Single theme         | Collapse 3 color files \+ 1,500+ hex literals into `theme.palette.*` | Medium | Low                         |
| **P4** Folder reshape       | Rename, flatten, and cluster `/components` into 6 domains            | Medium | Low (mostly mechanical)     |
| **P5** Break the monoliths  | Decompose the 8 giant files \+ dedupe 22 modals                      | Large  | Medium                      |

Nothing in this plan requires coordinated downtime; every phase is incremental and can ship behind small PRs.

---

## **1\. Current State — Measured**

All numbers below came from direct Grep/wc passes against `HEAD` on `claude/ui-restructure-plan-JuSmU`.

### **1.1 Inventory**

| Area                           | Count | Notes                                                           |
| ------------------------------ | ----- | --------------------------------------------------------------- |
| `.tsx` files                   | 451   |                                                                 |
| `.ts` files                    | 130   |                                                                 |
| `.js` / `.jsx` files           | 69    | Mostly legacy SVG-as-JS \+ a few page files                     |
| Files under `ui/components/`   | 418   | Across **30 top-level subfolders**                              |
| Files under `ui/pages/`        | 19    | Next.js pages router                                            |
| Files under `ui/rtk-query/`    | 26    | Clean — leave alone                                             |
| Files under `ui/store/slices/` | 6     | Clean — leave alone                                             |
| Files under `ui/utils/`        | 31    | Needs tidying                                                   |
| Files under `ui/themes/`       | 5     | **Redundant with Sistent**                                      |
| Files under `ui/constants/`    | 4     | `colors.ts` is dead weight                                      |
| Files under `ui/assets/icons/` | 76    | Competes with `ui/assets/new-icons/` and `ui/components/icons/` |

### **1.2 Design‑system sprawl**

| Import source             | Files   | Status                                                                      |
| ------------------------- | ------- | --------------------------------------------------------------------------- |
| `@sistent/sistent`        | **276** | ✅ Target                                                                   |
| `@mui/material`           | 64      | ❌ Eliminate                                                                |
| `@mui/icons-material`     | 89      | ❌ Eliminate                                                                |
| `@mui/x-date-pickers`     | \~6     | ❌ Eliminate (Sistent has its own)                                          |
| `@mui/x-tree-view`        | \~4     | ❌ Eliminate                                                                |
| `@rjsf/mui`               | 1       | ❌ Eliminate (`ui/components/MesheryMeshInterface/PatternService/RJSF.tsx`) |
| `@sistent/mui-datatables` | 6       | ⚠️ Acceptable only if Sistent wraps it; otherwise replace                   |
| `@material-ui/*` (v4)     | 0       | ✅ Already gone                                                             |

**100 files total** still touch `@mui/*` directly. Many files mix Sistent \+ MUI in the same import block — the single worst pattern in the codebase.

### **1.3 Theme & color sprawl**

Three independent color tables exist, none of them derived from Sistent:

// ui/themes/app.ts  
Colors \= { darkJungleGreen: '\#1E2117', caribbeanGreen: '\#00D3a9',  
 keppelGreen: '\#00B39F', charcoal: '\#3C494F' }  
notificationColors \= { error: '\#F91313', warning: '\#F0A303', ... }

// ui/themes/index.ts  
NOTIFICATIONCOLORS \= { ERROR: '\#F91313', WARNING: '\#F0A303',  
 SUCCESS: '\#206D24', INFO: '\#2196F3',  
 SUCCESS_V2: '\#3FC6B6', ERROR_DARK: '\#B32700' }

// ui/constants/colors.ts  
PRIMARY_COLOR \= '\#647881'; SUCCESS_COLOR \= '\#83B71E'; ...

Plus:

- **457 hex literals across 103 source files** (excluding auto‑generated static SVG JS)
- **680 `style={{…}}`** props across **182 files**
- **6 files** using `makeStyles` / `withStyles` / `createTheme` / `ThemeProvider` outside Sistent: `themes/rjsf.ts`, `pages/_app.tsx`, `components/Performance/PerformanceResults.tsx`, `components/SistentWrapper.tsx`, `components/PageContext.tsx`, `components/MesheryMeshInterface/PatternService/RJSF.tsx`
- **22 scattered `*.style(s).tsx` files** — each a private island of styled-components

### **1.4 Giant files (≥1,000 lines)**

| Lines | File                                                      |
| ----- | --------------------------------------------------------- |
| 1,679 | `components/MesheryPatterns/MesheryPatterns.tsx`          |
| 1,453 | `components/Dashboard/resources/configuration/config.tsx` |
| 1,376 | `components/MesheryFilters/Filters.tsx`                   |
| 1,298 | `components/Performance/index.tsx`                        |
| 1,291 | `components/connections/ConnectionTable.tsx`              |
| 1,273 | `components/Dashboard/resources/workloads/config.tsx`     |
| 1,230 | `components/MesheryAdapterPlayComponent.tsx`              |
| 1,142 | `components/Navigator.tsx`                                |
| 917   | `components/Dashboard/resources/network/config.tsx`       |
| 827   | `components/Settings/Registry/Stepper/UrlStepper.tsx`     |
| 750   | `components/Performance/PerformanceResults.tsx`           |
| 698   | `components/telemetry/grafana/GrafanaCustomChart.tsx`     |
| 689   | `components/connections/meshSync/index.tsx`               |
| 677   | `components/UserPreferences/index.tsx`                    |
| 673   | `pages/extensions.tsx`                                    |
| 653   | `components/Lifecycle/Environments/index.tsx`             |
| 605   | `pages/_app.tsx`                                          |

### **1.5 Conceptual duplication**

| Concept              | Scattered across                                                                                                                                                                                                           |
| -------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Modal**            | 22 files across `General/Modals/`, `Dashboard/`, `Registry/`, `Settings/Registry/`, `SpacesSwitcher/`, plus top‑level `ConfirmationModal.tsx`, `ExportModal.tsx`, `ViewInfoModal.tsx`, `TroubleshootingModalComponent.tsx` |
| **Workspaces**       | `SpacesSwitcher/` _and_ `Lifecycle/Workspaces/` _and_ `utils/context/WorkspaceModalContextProvider.tsx`                                                                                                                    |
| **Lifecycle**        | `Lifecycle/` _and_ `DesignLifeCycle/` _and_ `configuratorComponents/MeshModel/hooks/useDesignLifecycle.tsx`                                                                                                                |
| **Designs/Patterns** | `MesheryPatterns/` _and_ `configuratorComponents/` _and_ `pages/configuration/designs/`                                                                                                                                    |
| **Hooks**            | `components/hooks/` _and_ `utils/hooks/` _and_ `themes/hooks.tsx`                                                                                                                                                          |
| **Icons**            | `assets/icons/` (76) _and_ `assets/new-icons/` (25) _and_ `components/icons/` _and_ direct `@mui/icons-material`                                                                                                           |
| **Loading**          | `LoadingComponents/` with 9 files, multiple inline spinners elsewhere                                                                                                                                                      |
| **Empty State**      | `EmptyState/K8sContextEmptyState.tsx` _and_ `connections/meshSync/MeshSyncEmptyState.tsx` _and_ `Lifecycle/General/empty-state/`                                                                                           |
| **Card**             | `MesheryPatternCard`, `PerformanceCard`, `FiltersCard`, `MesheryWorkspaceCard`, `FlipCard`, `environment-card`                                                                                                             |
| **Stepper**          | `DeployStepper`, `RelationshipFormStepper`, `CSVStepper`, `UrlStepper`, `connections/meshSync/Stepper/*`                                                                                                                   |

### **1.6 Naming & layout**

- 27 `index.tsx` \+ 11 `index.ts` files — hard to navigate
- Top‑level folders mix `PascalCase` (`Dashboard`, `Lifecycle`) and lowercase (`connections`, `configuratorComponents`, `graphql`, `hooks`, `icons`, `schemas`, `shapes`, `subscription`, `telemetry`, `extensions`)
- Some files use `kebab-case.tsx` (`view-component.tsx`, `environment-card.tsx`) — 3 found
- `themes/` exists as a peer of `components/` even though Sistent owns theming

---

## **2\. Guiding Principles**

These are load‑bearing. The entire plan enforces these six rules.

1. **Sistent is the only UI kit.** No file outside Sistent's internals may import from `@mui/*`, `@material-ui/*`, or `@rjsf/mui`. Period.
2. **Tokens, never literals.** No hex, rgb, or named color string may appear in any `.tsx`/`.ts` file outside `ui/theme/` (a single _wrapper_ module — see §4). Colors come from `theme.palette.*`, spacing from `theme.spacing()`, breakpoints from `theme.breakpoints.*`.
3. **`styled()` over `style={{}}`.** Inline `style` props are reserved for _dynamic_ values that cannot be expressed in a styled component (e.g. `transform: translate(x,y)` for a draggable). Colors are never dynamic enough to justify inline.
4. **One concept, one home.** Every domain (Designs, Workspaces, Connections, …) lives in exactly one folder. Cross‑cutting primitives live under `ui/components/shared/`.
5. **Size budget: 400 lines.** Any component file \> 400 lines is a refactor candidate; \> 600 lines is blocked by lint warning; \> 1000 lines is a hard error in CI.
6. **Discoverable filenames.** No `index.tsx` as the _only_ file in a folder. No generic `style.tsx`, `utils.tsx`, `helpers.tsx`, `components.tsx` — files are named for what they contain.

---

## **3\. Target Directory Structure**

ui/  
├── pages/ \# Next.js routes — no change in shape, only imports  
│ ├── \_app.tsx \# Slimmed: provider stack only, \< 150 lines  
│ ├── \_document.tsx  
│ ├── index.tsx  
│ ├── designs/ \# was: configuration/designs  
│ ├── filters/ \# was: configuration/filters.tsx (promote)  
│ ├── catalog/ \# was: configuration/catalog.tsx (promote)  
│ ├── performance/  
│ ├── workspaces/ \# was: management/workspaces  
│ ├── environments/ \# was: management/environments  
│ ├── connections/ \# was: management/connections  
│ ├── adapters/ \# was: management/adapter  
│ ├── settings/  
│ ├── user/  
│ └── extensions/  
│  
├── components/  
│ ├── shared/ \# ⭐ NEW — cross‑cutting primitives  
│ │ ├── Modal/ \# 1 canonical Modal \+ Confirm/Info/Export flavors  
│ │ ├── Card/ \# ResourceCard, used by Patterns/Filters/Workspaces  
│ │ ├── Stepper/ \# Single Stepper used by Deploy/URL/CSV/Relationship  
│ │ ├── EmptyState/ \# 1 canonical EmptyState w/ variants  
│ │ ├── LoadingState/ \# Replaces LoadingComponents/  
│ │ ├── ErrorBoundary/  
│ │ ├── PageHeader/ \# Replaces inline headers in every page  
│ │ ├── DataTable/ \# Wraps @sistent/mui-datatables, used everywhere  
│ │ ├── FormFields/ \# RJSF custom widgets (current PatternService/RJSFCustomComponents/)  
│ │ ├── CodeEditor/  
│ │ └── Markdown/  
│ │  
│ ├── layout/ \# ⭐ NEW — app chrome  
│ │ ├── AppShell/ \# StyledRoot \+ StyledMainContent \+ StyledDrawer  
│ │ ├── Header/ \# was: Header.tsx \+ Header.styles.tsx \+ HeaderMenu.tsx  
│ │ ├── Navigator/ \# was: Navigator.tsx (1142 lines — see §7)  
│ │ ├── Footer/  
│ │ └── NotificationCenter/ \# moved here from top level  
│ │  
│ ├── designs/ \# ⭐ CONSOLIDATED — Patterns \+ Configurator  
│ │ ├── list/ \# was: MesheryPatterns/  
│ │ ├── builder/ \# was: configuratorComponents/MeshModel/  
│ │ ├── lifecycle/ \# was: DesignLifeCycle/  
│ │ ├── validation/ \# ValidateDesign, DryRun  
│ │ └── deployment/ \# DeployStepper, DeploymentSummary, finalize  
│ │  
│ ├── filters/ \# was: MesheryFilters/ (1376‑line Filters.tsx split)  
│ │  
│ ├── workspaces/ \# ⭐ CONSOLIDATED — merge SpacesSwitcher \+ Lifecycle/Workspaces  
│ │ ├── list/  
│ │ ├── switcher/ \# was: SpacesSwitcher/  
│ │ ├── content/ \# tabs, views, designs inside a workspace  
│ │ └── share/  
│ │  
│ ├── environments/ \# was: Lifecycle/Environments/  
│ ├── connections/ \# same, but ConnectionTable split  
│ ├── adapters/ \# was: MeshAdapterConfigComponent \+ MesheryAdapterPlayComponent  
│ ├── registry/ \# was: Settings/Registry/ — it never belonged under Settings  
│ │  
│ ├── performance/ \# Performance/index.tsx split by concern  
│ │ ├── profiles/  
│ │ ├── results/  
│ │ ├── dashboard/  
│ │ └── charts/  
│ │  
│ ├── dashboard/ \# was: Dashboard/ — widgets/resources/charts kept but flattened  
│ │ ├── widgets/  
│ │ ├── charts/  
│ │ └── resources/  
│ │ ├── definitions/ \# config.tsx files, renamed per resource  
│ │ ├── configuration.ts  
│ │ ├── workloads.ts  
│ │ ├── network.ts  
│ │ └── security.ts  
│ │  
│ ├── settings/ \# Slimmed: Registry moved out, Performance moved out  
│ ├── user/ \# UserPreferences/ renamed  
│ ├── telemetry/ \# kept — grafana \+ prometheus  
│ ├── extensions/ \# kept  
│ └── subscription/ \# kept (GraphQL subscription wiring)  
│  
├── theme/ \# ⭐ RENAMED & COLLAPSED from themes/  
│ ├── index.ts \# Re‑exports from @sistent/sistent — the ONLY color source  
│ ├── SistentProvider.tsx \# was: SistentWrapper.tsx \+ \_app.tsx provider soup  
│ ├── snackbar.tsx \# ThemeResponsiveSnackbar (uses theme.palette.\*)  
│ └── rjsf.ts \# Sistent‑backed RJSF theme (no Material UI)  
│  
├── store/ \# ✅ keep as‑is  
├── rtk-query/ \# ✅ keep as‑is  
├── machines/ \# ✅ keep as‑is  
│  
├── lib/ \# Third‑party integrations & network helpers  
├── utils/ \# Pure functions only — no JSX  
│ ├── hooks/ \# Generic hooks (useDebounce etc.)  
│ ├── context/ \# React contexts  
│ ├── format/ \# camelCase, PascalCase, kebabCase helpers consolidated  
│ └── k8s/ \# k8s-utils.ts, multi-ctx.ts  
│  
├── assets/ \# Static only — never contains .tsx components  
│ ├── images/  
│ ├── gifs/  
│ └── schemas/  
│  
├── graphql/ \# ⭐ MOVED from components/graphql/ (79 files — not components)  
│ ├── queries/  
│ ├── mutations/  
│ └── subscriptions/  
│  
├── tests/  
└── playwright/

**What goes away:**

- `ui/themes/` → `ui/theme/` (renamed, shrunk, Sistent‑backed)
- `ui/constants/colors.ts` → deleted
- `ui/css/` → absorbed into `ui/theme/` or deleted if unused
- `ui/components/icons/` → deleted (use `@sistent/sistent` icons or `ui/assets/icons/` SVG components)
- `ui/assets/new-icons/` → merged into `ui/assets/icons/` with consistent naming
- `ui/components/hooks/` → merged into `ui/utils/hooks/`
- `ui/components/General/` → split into `shared/` (Modal, Popup, TipsCarousel, ErrorBoundary) and `layout/` (pieces that are chrome)
- `ui/components/shapes/` → merged into `ui/assets/icons/shapes/`
- All `*.style(s).tsx` sibling files → colocated inside their component folder with unique names (`Header.styled.ts`, not `styles.tsx`)

---

## **4\. Theme & Color Consolidation**

### **4.1 Single source of truth**

Create `ui/theme/index.ts` as a **thin re‑export** — never a redefinition:

// ui/theme/index.ts  
export {  
 SistentThemeProvider,  
 useTheme,  
 styled,  
 alpha,  
 lighten,  
 darken,  
 CssBaseline,  
} from '@sistent/sistent';

// Named accessors that wrap theme.palette for clarity.  
// These exist ONLY so consumers never have to remember palette paths.  
export const palette \= {  
 status: {  
 error: (t: Theme) \=\> t.palette.error.main,  
 warning: (t: Theme) \=\> t.palette.warning.main,  
 success: (t: Theme) \=\> t.palette.success.main,  
 info: (t: Theme) \=\> t.palette.info.main,  
 },  
 surface: {  
 page: (t: Theme) \=\> t.palette.background.default,  
 elevated: (t: Theme) \=\> t.palette.background.elevatedComponents,  
 card: (t: Theme) \=\> t.palette.background.card,  
 },  
} as const;

### **4.2 Deletions**

| Delete                                                                        | Where callers move to                                |
| ----------------------------------------------------------------------------- | ---------------------------------------------------- |
| `ui/themes/app.ts` → `Colors`, `notificationColors`, `darkNotificationColors` | `theme.palette.*` (error/warning/success/info)       |
| `ui/themes/index.ts` → `NOTIFICATIONCOLORS`                                   | `theme.palette.*`                                    |
| `ui/constants/colors.ts` → `PRIMARY_COLOR`, etc.                              | `theme.palette.primary.main`, etc.                   |
| `ui/utils/lightenOrDarkenColor.ts`                                            | `import { lighten, darken } from '@sistent/sistent'` |

If Sistent's palette is missing a token the app needs (e.g. `#3FC6B6` "success_v2"), open an **upstream PR to Sistent**, not a local override. This rule is the whole point; without it the sprawl just regrows.

### **4.3 Migration mechanics**

1. Add an ESLint `no-restricted-syntax` rule (see §8) that errors on any hex literal in `.tsx`/`.ts` outside `ui/theme/`.
2. Search‑and‑replace each hex usage with its `theme.palette.*` equivalent, one file per commit.
3. For styled components, the migration is `({ theme }) => ({ color: theme.palette.error.main })`.
4. For inline `style={{ color: '#F91313' }}`, first assess whether it can become a styled component; if not, `style={{ color: theme.palette.error.main }}` using `useTheme()`.
5. Delete the 3 color files **only after** grep shows zero external imports.

### **4.4 RJSF theme**

`themes/rjsf.ts` currently uses Material UI's `createTheme`. Replace with Sistent's theme extension:

// ui/theme/rjsf.ts  
import { extendSistentTheme } from '@sistent/sistent';

export const rjsfTheme \= extendSistentTheme((base) \=\> ({  
 components: {  
 MuiTextField: { defaultProps: { size: 'small', variant: 'outlined' } },  
 // …only overrides that RJSF actually needs  
 },  
}));

`@rjsf/mui` is retained as a peer dep (RJSF requires a MUI adapter under the hood), but it is **wrapped** inside Sistent and never imported by app code. Only `ui/components/shared/FormFields/RJSFProvider.tsx` touches it.

---

## **5\. Material UI Elimination**

### **5.1 Import‑by‑import replacement table**

All mappings below are 1:1 unless otherwise noted. The explore pass confirmed Sistent exports these.

| From `@mui/material`                                                             | To `@sistent/sistent`                                         |
| -------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| `Box`, `Typography`, `Button`, `IconButton`, `Paper`, `Stack`, `Grid`, `Divider` | same names                                                    |
| `Dialog`, `DialogTitle`, `DialogContent`, `DialogActions`                        | same                                                          |
| `Drawer`, `AppBar`, `Toolbar`                                                    | same                                                          |
| `Tabs`, `Tab`, `tabsClasses`                                                     | `Tabs`, `Tab`, `tabsClasses`                                  |
| `TextField`, `Select`, `MenuItem`, `Checkbox`, `Radio`, `Switch`, `FormControl`  | same                                                          |
| `Table`, `TableHead`, `TableRow`, `TableCell`, `TableBody`                       | same                                                          |
| `Tooltip`                                                                        | `CustomTooltip` (Sistent's themed wrapper)                    |
| `Accordion`, `AccordionSummary`, `AccordionDetails`                              | same                                                          |
| `Autocomplete`                                                                   | same                                                          |
| `Menu`, `MenuList`, `Popover`, `Popper`                                          | same                                                          |
| `CircularProgress`, `LinearProgress`, `Skeleton`                                 | same                                                          |
| `Snackbar`                                                                       | `Snackbar` (or use `notistack` via `ThemeResponsiveSnackbar`) |
| `Chip`                                                                           | `Chip`                                                        |
| `Avatar`                                                                         | `Avatar` / `CustomAvatar`                                     |
| `styled`, `useTheme`, `alpha`, `lighten`, `darken`                               | same                                                          |
| `Hidden`, `NoSsr`                                                                | same                                                          |
| `Breadcrumbs`                                                                    | same                                                          |

| From `@mui/icons-material`                                                 | To                                                                                           |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| Any icon used ≥3×                                                          | add to `ui/assets/icons/` as a typed SVG component, exported from `ui/assets/icons/index.ts` |
| Any icon used once or twice                                                | replace with the Sistent equivalent if it exists                                             |
| Icons Sistent already ships (`Error`, `Warning`, `Info`, `CheckCircle`, …) | use `@sistent/sistent`                                                                       |

| From `@mui/x-date-pickers`              | To                                                                                                                                                                                  |
| --------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `LocalizationProvider`, `AdapterMoment` | Sistent's `SistentLocalizationProvider` if available; otherwise keep as a single wrapper in `ui/components/shared/DatePicker/` and don't import `@mui/x-date-pickers` anywhere else |
| `DateTimePicker`                        | `MesheryDateTimePicker.tsx` (already exists) — promote to `shared/DatePicker/` and remove direct MUI imports                                                                        |

| From `@mui/x-tree-view` | To                                                                                                                              |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `TreeView`, `TreeItem`  | `ui/components/shared/TreeView/` — a single wrapper. All `MesheryTreeView*.tsx` files (5 of them) consolidate into this folder. |

| From `@rjsf/mui` | To                                                                                            |
| ---------------- | --------------------------------------------------------------------------------------------- |
| `Form`           | `ui/components/shared/FormFields/RJSFProvider.tsx` — the _only_ file that imports `@rjsf/mui` |

### **5.2 Execution order (lowest risk first)**

1. **Leaf components** (no children importing MUI): small atoms like `ViewSwitch`, `FlipCard`, `PromptComponent`, `ReactSelectWrapper`. \~20 files.
2. **Icons-only files** (the 89 that only import from `@mui/icons-material`): pure mechanical replacement.
3. **Styled components files** (`*.style(s).tsx` — 22 files): migrate `styled` import source from `@mui/material` to `@sistent/sistent`.
4. **Feature containers** (Dashboard/index.tsx, Header.tsx, Performance/index.tsx): requires breakdown first (see §7).
5. **Form widgets** (`MesheryMeshInterface/PatternService/RJSFCustomComponents/*`): 8 files, must be done together because they share the RJSF form context.
6. **App shell** (`pages/_app.tsx` \+ `themes/App.styles.tsx`): last, because a break here blanks the screen.

### **5.3 What stays (justified exceptions)**

- **`notistack`** — no Sistent equivalent. Keep, but wrap in `ui/theme/snackbar.tsx`.
- **`@emotion/react`, `@emotion/styled`, `@emotion/cache`, `@emotion/server`** — Sistent _is_ Emotion. These stay as transitive infrastructure.
- **`billboard.js`**, **`react-grid-layout`**, **`react-big-calendar`**, **`react-select`**, **`@uiw/react-codemirror`** — not UI kits; these are specialized widgets. Each gets exactly one wrapper in `shared/` (e.g. `shared/Chart/`, `shared/GridLayout/`).

After P2, `package.json` should have:

\- "@mui/material": "^7.3.9",  
\- "@mui/icons-material": "^7.3.9",  
\- "@mui/x-date-pickers": "^8.27.2",  
\- "@mui/x-tree-view": "^8.27.2",  
\- "@rjsf/mui": "^6.4.1",

`@rjsf/mui` may remain as a transitive dep of Sistent — that is fine. What matters is zero direct imports from app code.

---

## **6\. Deduplication Plan**

### **6.1 Modal system**

**Today:** 22 modal files, each reinventing open/close, header, footer, confirm button styling.

**Target:** `ui/components/shared/Modal/` with exactly these files:

shared/Modal/  
├── Modal.tsx \# Base: children, onClose, title, size, actions  
├── ConfirmModal.tsx \# Wraps Modal: title, message, confirm/cancel  
├── InfoModal.tsx \# Wraps Modal: title, body (Markdown)  
├── FormModal.tsx \# Wraps Modal: integrates react-hook-form / RJSF  
├── useModal.ts \# Hook for imperative open/close  
└── index.ts

Migration map:

| Current file                                       | New home                                                       |
| -------------------------------------------------- | -------------------------------------------------------------- |
| `ConfirmationModal.tsx`                            | `ConfirmModal` (delete file)                                   |
| `ExportModal.tsx`                                  | `FormModal` instance in `designs/export/ExportDesignModal.tsx` |
| `ViewInfoModal.tsx`                                | `InfoModal` instance in `workspaces/ViewInfoModal.tsx`         |
| `General/Modals/GenericModal.tsx`                  | `Modal` (delete file)                                          |
| `General/Modals/Modal.tsx`                         | Delete (consolidates into base `Modal`)                        |
| `General/Modals/InfoModal.tsx`                     | `InfoModal` (delete)                                           |
| `General/Modals/EnvironmentModal.tsx`              | `environments/EnvironmentFormModal.tsx`                        |
| `General/Modals/ConnectionModal.tsx`               | `connections/ConnectionFormModal.tsx`                          |
| `General/Modals/PublishModal.tsx`                  | `designs/PublishDesignModal.tsx`                               |
| `General/Modals/ImportModal.tsx`                   | `designs/ImportDesignModal.tsx`                                |
| `Dashboard/UnsavedChangesModal.tsx`                | `shared/Modal/UnsavedChangesModal.tsx`                         |
| `RelationshipBuilder/CreateRelationshipModal.tsx`  | `registry/CreateRelationshipModal.tsx`                         |
| `Registry/RegistryModal.tsx`                       | `registry/RegistryModal.tsx`                                   |
| `Settings/Registry/ImportModelModal.tsx`           | `registry/ImportModelModal.tsx`                                |
| `Settings/Registry/CreateModelModal.tsx`           | `registry/CreateModelModal.tsx`                                |
| `TroubleshootingModalComponent.tsx`                | `shared/Troubleshooting/TroubleshootingModal.tsx`              |
| `SpacesSwitcher/WorkspaceModal.tsx`                | `workspaces/WorkspaceFormModal.tsx`                            |
| `SpacesSwitcher/ShareModal.tsx`                    | `workspaces/ShareWorkspaceModal.tsx`                           |
| `connections/meshSync/RegisterConnectionModal.tsx` | `connections/RegisterConnectionModal.tsx`                      |
| `pages/extension/AccessMesheryModal.tsx`           | `layout/AccessMesheryModal.tsx`                                |

Every one of these becomes either `<Modal>`, `<ConfirmModal>`, `<InfoModal>`, or `<FormModal>`. The top‑level `GenericModal` / `Modal.tsx` / `InfoModal.tsx` duplicates _disappear_.

### **6.2 Card system**

Create `shared/Card/ResourceCard.tsx` — a single flippable, menu‑enabled, status‑badge card. Today's `MesheryPatternCard`, `FiltersCard`, `PerformanceCard`, `MesheryWorkspaceCard`, `environment-card`, `FlipCard` all become thin instances passing content slots. Expected LOC savings: \~800.

### **6.3 Stepper system**

Create `shared/Stepper/Stepper.tsx` that takes a `steps` array. Today's `DeployStepper`, `RelationshipFormStepper`, `CSVStepper`, `UrlStepper`, `connections/meshSync/Stepper/*` all become data‑driven instances. Expected LOC savings: \~600.

### **6.4 EmptyState, LoadingState, PageHeader**

- `EmptyState` — replace `LoadingComponents/`, `EmptyState/K8sContextEmptyState.tsx`, `connections/meshSync/MeshSyncEmptyState.tsx`, `Lifecycle/General/empty-state/`, and the ad‑hoc empty‑states inside `SpacesSwitcher/*Content.tsx` files.
- `LoadingState` — absorbs `LoadingComponents/Animations/*` (6 animation files). Each animation is a variant.
- `PageHeader` — every `pages/*` currently hand‑rolls a title bar. Replace with `<PageHeader title={…} actions={…} breadcrumbs={…} />`.

### **6.5 Icons**

One folder: `ui/assets/icons/`. Rules:

1. Delete `ui/components/icons/` (1 file, an index.ts).
2. Merge `ui/assets/new-icons/*.svg` (25 files) into `ui/assets/icons/` as `.tsx` components (SVGs become React components for consistency with the existing 76 icons).
3. Delete `ui/components/shapes/` (2 files) → `ui/assets/icons/shapes/`.
4. Any `@mui/icons-material` import is replaced by a Sistent icon or a typed SVG under `ui/assets/icons/`.

### **6.6 Hooks consolidation**

- `ui/components/hooks/` (3 files) \+ `ui/utils/hooks/` → merged under `ui/utils/hooks/`.
- `ui/themes/hooks.tsx` (theme toggler) → `ui/theme/useThemePreference.ts`.

### **6.7 Workspaces consolidation**

`SpacesSwitcher/` (18 files) and `Lifecycle/Workspaces/` (10 files) are two halves of the same feature. Merge into `ui/components/workspaces/`:

workspaces/  
├── switcher/ \# was: SpacesSwitcher/SpaceSwitcher, MenuComponent, MobileViewSwitcher  
├── list/ \# was: Lifecycle/Workspaces/Workspace{Grid,Data}View  
├── content/ \# was: SpacesSwitcher/\*Content.tsx (5 files — should be 1 component \+ variants)  
├── card/ \# was: MesheryWorkspaceCard (now instance of shared/Card)  
├── modals/ \# WorkspaceFormModal, ShareModal  
├── hooks/ \# was: SpacesSwitcher/hooks.tsx  
└── styles/ \# was: SpacesSwitcher/styles.tsx \+ Lifecycle/Workspaces/styles.tsx merged

The 5 `*Content.tsx` files in `SpacesSwitcher/` (MainDesignsContent, MainViewsContent, MyDesignsContent, MyViewsContent, RecentContent, SharedContent, WorkspaceContent) are all the same component with different props — collapse to 1 file with a `variant` prop.

### **6.8 Lifecycle consolidation**

`Lifecycle/` and `DesignLifeCycle/` are unrelated despite the name collision.

- `Lifecycle/Environments/` → `components/environments/`
- `Lifecycle/Workspaces/` → `components/workspaces/` (merged above)
- `Lifecycle/General/` → absorbed into `shared/`
- `DesignLifeCycle/` → `components/designs/lifecycle/`
- Delete `Lifecycle/` top‑level folder.

---

## **7\. Breaking the Giants**

### **7.1 `MesheryPatterns.tsx` (1,679 lines)**

Split by responsibility into `components/designs/list/`:

designs/list/  
├── DesignsPage.tsx \# The component; \~200 lines of composition  
├── DesignsTable.tsx \# Table view  
├── DesignsGrid.tsx \# Grid/card view  
├── DesignsToolbar.tsx \# Search, filters, view toggle  
├── DesignsBulkActions.tsx \# Selection toolbar  
├── DesignsEmptyState.tsx \# Instance of shared/EmptyState  
├── hooks/  
│ ├── useDesigns.ts \# RTK Query \+ local filter state  
│ ├── useDesignActions.ts \# Clone, delete, share, publish  
│ └── useDesignSelection.ts  
└── DesignsPage.styled.ts \# styled() components

### **7.2 `Dashboard/resources/configuration/config.tsx` (1,453 \+ siblings)**

These are _table column definitions_, not React components. They should not be `.tsx`:

dashboard/resources/definitions/  
├── configuration.ts \# was: configuration/config.tsx  
├── workloads.ts \# was: workloads/config.tsx (1273 lines)  
├── network.ts \# was: network/config.tsx (917 lines)  
├── security.ts \# was: security/config.tsx (582 lines)  
├── nodes.ts  
├── types.ts \# Shared ResourceDefinition type  
└── formatters.tsx \# Cell renderers (the only .tsx file)

A `ResourceDefinition` type centralizes the pattern. Then `resources-table.tsx` takes a `definition` prop. Expected LOC reduction: \~2,000 lines net across the split.

### **7.3 `MesheryFilters/Filters.tsx` (1,376 lines)**

Same pattern as `MesheryPatterns.tsx`. Split into `components/filters/`:

filters/  
├── FiltersPage.tsx  
├── FiltersTable.tsx  
├── FiltersGrid.tsx \# already exists — absorbs FiltersCard  
├── FiltersToolbar.tsx  
├── hooks/useFilters.ts  
└── FiltersPage.styled.ts

### **7.4 `Performance/index.tsx` (1,298 lines)**

Split into `components/performance/`:

performance/  
├── PerformancePage.tsx \# Tab container  
├── profiles/  
│ ├── ProfilesList.tsx \# was: PerformanceProfiles.tsx (591 lines) split  
│ ├── ProfilesGrid.tsx \# was: PerformanceProfileGrid.tsx  
│ └── ProfileCard.tsx \# was: PerformanceCard.tsx — instance of shared/Card  
├── results/  
│ ├── ResultsTable.tsx \# was: PerformanceResults.tsx (750 lines) split  
│ └── ResultDetail.tsx  
├── dashboard/  
│ ├── PerformanceDashboard.tsx \# was: Dashboard.tsx  
│ └── MesheryMetrics.tsx  
├── calendar/  
│ └── PerformanceCalendar.tsx  
└── hooks/

### **7.5 `connections/ConnectionTable.tsx` (1,291 lines)**

connections/  
├── ConnectionsPage.tsx  
├── ConnectionsTable.tsx \# the data table  
├── ConnectionRow.tsx \# row renderer (\~150 lines)  
├── ConnectionActions.tsx \# menu, bulk actions  
├── ConnectionFilters.tsx \# status/kind filters  
├── ConnectionDetail.tsx \# drawer content  
├── ConnectionChip.tsx \# already exists  
├── hooks/  
│ ├── useConnections.ts  
│ └── useConnectionActions.ts  
├── meshSync/ \# kept, but \~689 lines of index.tsx split  
└── ConnectionsPage.styled.ts

### **7.6 `MesheryAdapterPlayComponent.tsx` \+ `MeshAdapterConfigComponent.tsx`**

Merge into `components/adapters/`:

adapters/  
├── AdaptersPage.tsx  
├── AdapterConfigForm.tsx \# was: MeshAdapterConfigComponent  
├── AdapterOpsPanel.tsx \# "Play" — the ops runner  
├── AdapterCard.tsx \# instance of shared/Card  
├── hooks/useAdapter.ts  
└── AdaptersPage.styled.ts

### **7.7 `Navigator.tsx` (1,142 lines)**

Split into `components/layout/Navigator/`:

layout/Navigator/  
├── Navigator.tsx \# \<200 lines: composition  
├── NavigatorItem.tsx \# single menu item  
├── NavigatorGroup.tsx \# expandable section  
├── NavigatorHeader.tsx \# logo \+ collapse button  
├── NavigatorFooter.tsx \# org switcher, help  
├── useNavigatorState.ts \# collapse, active route, permissions  
├── navigationSchema.ts \# data‑driven menu definition (was: constants/navigator.ts)  
└── Navigator.styled.ts

### **7.8 `pages/_app.tsx` (605 lines) & `themes/App.styles.tsx`**

The current `_app.tsx` mixes: provider stack, theme setup, snackbar config, subscription wiring, RTK Query bootstrap, route guards, connection metadata fetching. Split into:

pages/\_app.tsx \# \<150 lines — only provider composition  
components/layout/AppShell/  
 ├── AppShell.tsx \# StyledRoot \+ Drawer \+ Main (was: App.styles.tsx)  
 ├── AppProviders.tsx \# Sistent, Redux, Relay, Notifications  
 ├── AppBootstrap.tsx \# Effect‑heavy: fetchSystemSync, fetchK8s, fetchOrgs  
 └── AppShell.styled.ts  
theme/  
 ├── SistentProvider.tsx \# was: SistentWrapper.tsx  
 └── snackbar.tsx \# was: ThemeResponsiveSnackbar in App.styles.tsx

### **7.9 `UserPreferences/index.tsx` (677 lines)**

Split into tabs (each was already an implicit section inside the file):

user/preferences/  
├── UserPreferencesPage.tsx  
├── tabs/  
│ ├── GeneralTab.tsx  
│ ├── NotificationsTab.tsx  
│ ├── ExtensionsTab.tsx  
│ ├── ThemeTab.tsx  
│ └── KeysTab.tsx  
├── hooks/useUserPreferences.ts  
└── UserPreferencesPage.styled.ts

### **7.10 `pages/extensions.tsx` (673 lines) & `Lifecycle/Environments/index.tsx` (653 lines)**

Same template as above: split into `Page.tsx` \+ domain‑specific tabs/tables \+ hooks \+ styled file.

---

## **8\. Enforcement — Make Regression Impossible**

Lint rules added to `ui/eslint.config.js`. Each rule ships with a codemod or has an empty allowlist so CI is green on the day it lands, then the allowlist shrinks PR by PR.

### **8.1 Ban Material UI imports**

// eslint.config.js  
{  
 rules: {  
 'no-restricted-imports': \['warn', {  
 paths: \[  
 { name: '@mui/material', message: 'Use @sistent/sistent' },  
 { name: '@mui/icons-material', message: 'Use @sistent/sistent icons or ui/assets/icons' },  
 { name: '@mui/x-date-pickers', message: 'Use components/shared/DatePicker' },  
 { name: '@mui/x-tree-view', message: 'Use components/shared/TreeView' },  
 { name: '@rjsf/mui', message: 'Use components/shared/FormFields/RJSFProvider' },  
 \],  
 patterns: \['@mui/\*', '@material-ui/\*'\],  
 }\],  
 },  
}

### **8.2 Ban hex literals outside `ui/theme/`**

{  
 files: \['\*\*/\*.{ts,tsx}'\],  
 ignores: \['ui/theme/\*\*', 'ui/assets/icons/\*\*'\],  
 rules: {  
 'no-restricted-syntax': \['error', {  
 selector: "Literal\[value=/^\#\[0-9a-fA-F\]{3,8}$/\]",  
 message: 'Hex colors forbidden. Use theme.palette.\*',  
 }, {  
 selector: "Literal\[value=/rgba?\\\\(/\]",  
 message: 'Use theme.palette.\* (or alpha() from @sistent/sistent).',  
 }\],  
 },  
}

### **8.3 Ban `style={{ color / backgroundColor }}` in components**

{  
 rules: {  
 'react/forbid-dom-props': \['warn', {  
 forbid: \[{  
 propName: 'style',  
 message: 'Use styled() from @sistent/sistent; inline style is reserved for dynamic geometry.',  
 }\],  
 }\],  
 },  
}

Applied as `warn` initially with a file‑level allowlist; promoted to `error` once allowlist is drained.

### **8.4 File size budget**

{  
 rules: {  
 'max-lines': \['error', { max: 1000, skipComments: true, skipBlankLines: true }\],  
 },  
}  
// \+ 'warn' at 600 for proactive nudging

### **8.5 Ban theme/color legacy imports**

{  
 rules: {  
 'no-restricted-imports': \['error', {  
 paths: \[  
 { name: '@/theme/index', message: 'Use @/theme; do not deep-import the local theme entry point.' },  
 { name: '@/themes', message: 'Use @/theme, the approved Phase 1 theme entry point.' },  
 { name: '@/themes/app', message: 'Use @/theme and theme.palette.\*' },  
 { name: '@/themes/index', message: 'Use @/theme and theme.palette.\*' },  
 { name: '@/constants/colors', message: 'Use @/theme and theme.palette.\*' },  
 \],  
 }\],  
 },  
}

### **8.6 CI gates**

- **`lint:ui-kit`** — runs the rules above on CI. Failing rule \= failing build.
- **`audit:mui`** — nightly script that greps for `@mui/` imports and posts a GitHub comment with the file count. Must trend to zero.
- **`audit:hex`** — same for hex literals.
- **Stat tracker** — posts a "giant files" count to the PR summary on every PR.

---

## **9\. Documentation**

Add under `ui/docs/`:

ui/docs/  
├── ARCHITECTURE.md \# This plan's target state, prose‑ified  
├── COMPONENT_GUIDE.md \# How to write a new component (folder layout, file naming, Sistent patterns)  
├── THEMING.md \# How to use theme.palette.\*, when to add tokens upstream in Sistent  
├── MIGRATION_NOTES.md \# Codemods, gotchas encountered per phase  
└── templates/  
 ├── NewPage.template.tsx  
 ├── NewComponent.template.tsx  
 └── NewFeatureFolder/ \# Example folder structure

The README at `ui/README.md` grows from 6 lines to \~50, pointing at `ui/docs/` and summarizing the six principles from §2.

---

## **10\. Phased Roadmap**

Each phase ships in small PRs, not one monolith. Estimates below are PR counts, not time.

### **Phase 1 — Guardrails (3–5 PRs)**

- Add ESLint rules in "warn \+ allowlist" mode
- Add CI audit scripts
- Land `ui/docs/ARCHITECTURE.md` \+ `THEMING.md`
- Add `ui/theme/index.ts` (re‑exporting Sistent) — but do not delete old files yet
- **Exit criteria:** CI reports baseline counts (100 MUI files, 457 hex literals, 8 giant files). No behavior change.

### **Phase 2 — Theme & Color Collapse (8–12 PRs)**

- Migrate `themes/rjsf.ts` to Sistent extension
- Replace all imports of `themes/app`, `themes/index`, `constants/colors` with `theme.palette.*` — one folder per PR
- Delete `themes/app.ts`, `themes/index.ts`, `constants/colors.ts`, `utils/lightenOrDarkenColor.ts`
- Rename `ui/themes/` → `ui/theme/`, move `SistentWrapper.tsx` inside
- Drain hex‑literal allowlist; flip rule to `error`
- **Exit criteria:** `audit:hex` reports 0 in `.tsx/.ts` outside `ui/theme/` and `ui/assets/icons/`.

### **Phase 3 — MUI Elimination (12–18 PRs)**

- PR batches in the order of §5.2 (leaves → icons → styles → containers → forms → shell)
- Create the small `shared/` wrappers needed for date picker, tree view, RJSF, data table as each batch requires them
- Remove `@mui/*` and `@rjsf/mui` from `package.json` in the final PR of the phase
- **Exit criteria:** `audit:mui` reports 0; `npm ls` shows MUI only as a transitive dep.

### **Phase 4 — Folder Reshape (6–10 PRs)**

One PR per move. Every PR uses `git mv` so history is preserved. Each PR:

- Moves files
- Updates imports (including `tsconfig.json` path aliases)
- Runs `npm run lint`, `npm run test`, `npm run build`

Move order (dependencies first):

1. `ui/components/graphql/` → `ui/graphql/`
2. `ui/components/hooks/` → `ui/utils/hooks/`
3. `ui/components/icons/` → deleted; `ui/assets/new-icons/` merged into `ui/assets/icons/`
4. `ui/components/shapes/` → `ui/assets/icons/shapes/`
5. Create `ui/components/shared/` and `ui/components/layout/`; move primitives
6. `ui/components/SpacesSwitcher/` \+ `ui/components/Lifecycle/Workspaces/` → `ui/components/workspaces/`
7. `ui/components/Lifecycle/Environments/` → `ui/components/environments/`
8. `ui/components/Settings/Registry/` → `ui/components/registry/`
9. `ui/components/MesheryPatterns/` \+ `ui/components/configuratorComponents/` \+ `ui/components/DesignLifeCycle/` → `ui/components/designs/`
10. `ui/components/MesheryFilters/` → `ui/components/filters/`
11. Rename remaining lowercase folders to PascalCase (or settle on lowercase — pick one and enforce)
12. Delete `ui/components/Lifecycle/`, `ui/components/General/`, `ui/components/Settings/` top‑level after their contents are redistributed

- **Exit criteria:** 30 top‑level component folders → \~15 (see §3 target). Max nesting depth ≤ 6\.

### **Phase 5 — Break the Giants & Dedupe (12–20 PRs)**

One PR per giant file in §7. Then one PR per duplicated primitive in §6 (Modal, Card, Stepper, EmptyState, LoadingState, PageHeader).

- **Exit criteria:** No file \> 600 lines outside `ui/graphql/` (generated). Modal/Card/Stepper count in app code drops to 1 canonical each. `max-lines` rule promoted to `error` at 600\.

---

## **11\. Acceptance Criteria**

The restructure is "done" when all of these are true:

- \[ \] `grep -r '@mui/' ui/ --include='*.tsx' --include='*.ts'` returns **0**
- \[ \] `grep -r '@rjsf/mui' ui/components ui/pages` returns **0**
- \[ \] `grep -rE '#[0-9a-fA-F]{6}' ui/components ui/pages ui/utils --include='*.tsx' --include='*.ts'` returns **0**
- \[ \] `ui/themes/app.ts`, `ui/themes/index.ts`, `ui/constants/colors.ts`, `ui/utils/lightenOrDarkenColor.ts` deleted
- \[ \] `ui/theme/index.ts` exists and is the sole export point for theme utilities
- \[ \] Every file under `ui/components/` is ≤ 600 lines (≤ 1000 hard cap enforced in CI)
- \[ \] `ui/components/` has ≤ 16 top‑level folders with consistent casing
- \[ \] Exactly one `Modal` / `Card` / `Stepper` / `EmptyState` / `LoadingState` primitive, all under `ui/components/shared/`
- \[ \] `pages/_app.tsx` is ≤ 150 lines
- \[ \] `ui/docs/ARCHITECTURE.md`, `THEMING.md`, `COMPONENT_GUIDE.md` merged
- \[ \] ESLint rules from §8 are `error`, not `warn`, with zero allowlisted files
- \[ \] `npm run build` succeeds; `npm run test` passes; `npm run test:e2e:ci` passes
- \[ \] Visual regression check passes on Dashboard, Designs, Performance, Connections, Workspaces, Settings (the 6 routes most affected)
- \[ \] `package.json` no longer lists `@mui/material`, `@mui/icons-material`, `@mui/x-date-pickers`, `@mui/x-tree-view`, `@rjsf/mui` as direct deps

---

## **12\. Risks & Mitigations**

| Risk                                                | Mitigation                                                                                                                                                                                          |
| --------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Sistent is missing a component MUI provides         | File upstream PR to Sistent first. If blocker, add a one‑file wrapper in `shared/` that _only_ it imports MUI — allowlisted in ESLint with a `TODO(sistent-upstream)` comment and a tracking issue. |
| `@sistent/mui-datatables` is really MUI             | Already scoped. Wrap it once in `shared/DataTable/` and never import it again.                                                                                                                      |
| RJSF requires an MUI adapter                        | `@rjsf/mui` stays as a transitive dep inside Sistent's RJSF wrapper; app code only sees Sistent.                                                                                                    |
| Giant‑file splits change behavior                   | Each split PR must include screenshot diffs of the affected page and a link to the playwright e2e run.                                                                                              |
| Theme migration breaks dark mode                    | Test both modes on every PR in Phase 2\. The `useTheme` hook already exposes `palette.mode`.                                                                                                        |
| Folder moves break IDE go‑to‑definition on open PRs | Sequence Phase 4 so each move PR is small (1–3 folders) and merges fast. Announce in \#meshery‑ui Slack before each move PR.                                                                        |
| Visual regressions                                  | Add a Playwright visual-snapshot job before Phase 2 (takes baseline screenshots of 10 key routes). Run it on every PR in Phases 2, 3, 5\.                                                           |
| Contributor PRs in flight                           | Land Phase 1 (guardrails only) first so new code doesn't add regressions. Communicate phase boundaries in CONTRIBUTING.                                                                             |

---

## **13\. What This Plan Deliberately Does Not Do**

To keep scope honest:

- **No upgrade of Sistent itself.** If `@sistent/sistent` needs a new version (`0.18.4` → latest), that's a separate coordinated upgrade.
- **No change to `rtk-query/`, `store/`, `machines/`.** They are clean.
- **No change to the Next.js pages router setup.** The plan reshapes `components/`, not routing.
- **No TypeScript strictness push.** Tempting, but orthogonal.
- **No test‑coverage campaign.** Also orthogonal; current vitest setup is fine.
- **No design‑token rename in Sistent.** We consume what Sistent exports; we don't rename its tokens.

---

## **14\. Single‑Line Summary**

Consolidate 30 folders into \~15, delete 3 color files and \~100 MUI imports, split 8 giant files, dedupe 22 modals into 1 primitive, and enforce the target state with ESLint — in five phased PR batches that are each shippable in isolation.
