# Meshery UI Architecture

This document describes the architecture of the Meshery UI as it exists during and after the multi-phase restructure tracked in [meshery/meshery#18656](https://github.com/meshery/meshery/issues/18656). It is the contributor's map: a short reference for deciding where new code goes, which UI kit and tokens to use, and when a file is too big. For the full long-form spec — the inventory of what exists today, the migration mechanics, ESLint rules, and the phased PR roadmap — read [`../restructure-plan.md`](../restructure-plan.md).

This is a living document. The target architecture is being landed incrementally over five phases; the principles and directory model below are already binding for new code, even where the existing tree has not yet been reshaped.

---

## Why this restructure exists

The Meshery UI is fundamentally sound — 87% TypeScript, RTK Query for data, and `@sistent/sistent` already imported in 276 files — but it has accreted three overlapping styling systems, two parallel icon libraries, two parallel workspace concepts, and eight files over 1,000 lines. The restructure is **not a rewrite**. It is a disciplined consolidation that:

- Enthrones `@sistent/sistent` as the single source of truth for components, tokens, and theme.
- Eliminates direct `@mui/*` imports from application code.
- Collapses three color tables and ~1,500 hex literals into `theme.palette.*`.
- Reshapes 30 top-level `components/` folders into ~15 discoverable domains.
- Breaks the eight giant (>1,000-line) files into focused units and dedupes 22 modals into a single primitive.

See the parent epic [#18656](https://github.com/meshery/meshery/issues/18656) for the full motivation.

---

## The six guiding principles

These are load-bearing. Every rule below answers a single question: **what does this mean when I write a new component?**

### 1. Sistent is the only UI kit

Import every component you need from `@sistent/sistent` — never from `@mui/material`, `@mui/icons-material`, `@material-ui/*`, or `@rjsf/mui`. Hooks and styling utilities (`styled`, `useTheme`, `alpha`, palette accessors) come from `@/theme`, the project-local front door that re-exports the supported Sistent surface; see [`THEMING.md`](./THEMING.md) for the full theming contract. If Sistent is missing something you need, file an upstream PR to Sistent rather than reaching past it. The single exception is RJSF: `@rjsf/mui` is allowed only inside Sistent's wrapper (`ui/components/shared/FormFields/`) and is invisible to the rest of the app.

```tsx
// Good
import { Button, Dialog } from '@sistent/sistent';
import { useTheme, styled } from '@/theme';

// Forbidden in application code
import Button from '@mui/material/Button';
```

### 2. Tokens, never literals

No hex, rgb, or named color string belongs in any `.tsx` or `.ts` file outside `ui/theme/`. Colors come from `theme.palette.*`, spacing from `theme.spacing()`, and breakpoints from `theme.breakpoints.*`. If Sistent's palette is missing a token your design calls for, open an upstream PR to Sistent — do not invent a local override.

```tsx
// Good
const Banner = styled('div')(({ theme }) => ({
  color: theme.palette.error.main,
  padding: theme.spacing(2),
}));

// Forbidden
const Banner = styled('div')({
  color: '#F91313',
  padding: '16px',
});
```

### 3. `styled()` over `style={{}}`

Inline `style` props are reserved for genuinely dynamic values that cannot be expressed in a styled component — e.g. `transform: translate(x, y)` for a draggable element. Colors, spacing, and typography are never dynamic enough to justify inline. Reach for `styled()` (re-exported from `@/theme`) first, and only fall back to `style={{}}` when the value changes per render in a way the CSS engine can't express.

### 4. One concept, one home

Every domain — Designs, Workspaces, Connections, Environments, Performance, Dashboard, and the rest — lives in exactly one folder under `ui/components/`. Cross-cutting primitives (Modal, Card, Stepper, EmptyState, …) live under `ui/components/shared/`. If you find yourself building a second Workspaces-ish component, stop: extend the canonical one or move it to `shared/`.

### 5. Size budget: 400 / 600 / 1000

Files have a budget:

- **400 lines** — soft target; start thinking about a split.
- **600 lines** — lint warning; an explicit refactor is overdue.
- **1000 lines** — hard error in CI; the PR will not merge.

Splits should be **by concern** (extract a hook, lift a sub-component into its own file, move types into a sibling `*.types.ts`), not by arbitrary line count. The budget exists to make the codebase navigable, not to encourage shallow file shuffling.

### 6. Discoverable filenames

A folder containing only `index.tsx` is hard to navigate in an IDE; name files for what they contain. Generic names like `style.tsx`, `utils.tsx`, `helpers.tsx`, and `components.tsx` are also forbidden. Prefer `Header.tsx`, `Header.styled.ts`, `Header.types.ts`, `useHeaderState.ts` — names that survive a file-search.

```text
# Good
components/layout/Header/
  Header.tsx
  Header.styled.ts
  Header.types.ts
  useHeaderState.ts

# Bad
components/layout/Header/
  index.tsx
  styles.tsx
  utils.tsx
```

---

## Target directory structure

The post-restructure `ui/` layout. Areas marked **NEW** or **CONSOLIDATED** are introduced by the restructure; areas marked **keep as-is** are already in good shape.

```text
ui/
├── pages/                          # Next.js routes — shape unchanged, only imports
│   ├── _app.tsx                    # Slimmed to provider stack only, < 150 lines
│   ├── _document.tsx
│   ├── index.tsx
│   ├── designs/                    # was: configuration/designs
│   ├── filters/                    # was: configuration/filters.tsx (promoted)
│   ├── catalog/                    # was: configuration/catalog.tsx (promoted)
│   ├── performance/
│   ├── workspaces/                 # was: management/workspaces
│   ├── environments/               # was: management/environments
│   ├── connections/                # was: management/connections
│   ├── adapters/                   # was: management/adapter
│   ├── settings/
│   ├── user/
│   └── extensions/
│
├── components/
│   ├── shared/                     # NEW — cross-cutting primitives
│   │   ├── Modal/                  # 1 canonical Modal + Confirm/Info/Export flavors
│   │   ├── Card/                   # ResourceCard for Patterns/Filters/Workspaces
│   │   ├── Stepper/                # Single Stepper used by Deploy/URL/CSV/Relationship
│   │   ├── EmptyState/             # 1 canonical EmptyState with variants
│   │   ├── LoadingState/           # Replaces LoadingComponents/
│   │   ├── ErrorBoundary/
│   │   ├── PageHeader/             # Replaces inline page headers
│   │   ├── DataTable/              # Wraps @sistent/mui-datatables
│   │   ├── FormFields/             # RJSF custom widgets
│   │   ├── CodeEditor/
│   │   └── Markdown/
│   │
│   ├── layout/                     # NEW — app chrome
│   │   ├── AppShell/               # StyledRoot + StyledMainContent + StyledDrawer
│   │   ├── Header/                 # was: Header.tsx + Header.styles.tsx + HeaderMenu.tsx
│   │   ├── Navigator/              # was: Navigator.tsx (1142-line giant, split)
│   │   ├── Footer/
│   │   └── NotificationCenter/
│   │
│   ├── designs/                    # CONSOLIDATED — Patterns + Configurator
│   │   ├── list/                   # was: MesheryPatterns/
│   │   ├── builder/                # was: configuratorComponents/MeshModel/
│   │   ├── lifecycle/              # was: DesignLifeCycle/
│   │   ├── validation/
│   │   └── deployment/
│   │
│   ├── filters/                    # was: MesheryFilters/ (1376-line Filters.tsx split)
│   │
│   ├── workspaces/                 # CONSOLIDATED — SpacesSwitcher + Lifecycle/Workspaces
│   │   ├── list/
│   │   ├── switcher/
│   │   ├── content/
│   │   └── share/
│   │
│   ├── environments/               # was: Lifecycle/Environments/
│   ├── connections/                # ConnectionTable split (was 1291 lines)
│   ├── adapters/                   # was: MeshAdapterConfigComponent + MesheryAdapterPlayComponent
│   ├── registry/                   # was: Settings/Registry/ (lifted out of Settings)
│   │
│   ├── performance/                # Performance/index.tsx split by concern
│   │   ├── profiles/
│   │   ├── results/
│   │   ├── dashboard/
│   │   └── charts/
│   │
│   ├── dashboard/                  # was: Dashboard/ — flattened
│   │   ├── widgets/
│   │   ├── charts/
│   │   └── resources/
│   │
│   ├── settings/                   # Slimmed — Registry and Performance moved out
│   ├── user/                       # was: UserPreferences/
│   ├── telemetry/                  # kept — grafana + prometheus
│   ├── extensions/                 # kept
│   └── subscription/               # kept — GraphQL subscription wiring
│
├── theme/                          # RENAMED & COLLAPSED from themes/
│   ├── index.ts                    # Re-exports from @sistent/sistent — the ONLY source for theme & styling
│   ├── SistentProvider.tsx         # was: SistentWrapper.tsx + provider plumbing
│   ├── snackbar.tsx                # ThemeResponsiveSnackbar (uses theme.palette.*)
│   └── rjsf.ts                     # Sistent-backed RJSF theme — no Material UI
│
├── store/                          # keep as-is — RTK store + slices
├── rtk-query/                      # keep as-is — RTK Query endpoints
├── machines/                       # keep as-is — XState machines
│
├── lib/                            # Third-party integrations & network helpers
├── utils/                          # Pure functions only — no JSX
│   ├── hooks/                      # Generic hooks (useDebounce, etc.)
│   ├── context/                    # React contexts
│   ├── format/                     # camelCase, PascalCase, kebabCase helpers
│   └── k8s/                        # k8s-utils.ts, multi-ctx.ts
│
├── assets/                         # Static only — never contains .tsx components
│   ├── images/
│   ├── gifs/
│   ├── icons/                      # SVG components, the one icon home
│   └── schemas/
│
├── graphql/                        # MOVED from components/graphql/ — not components
│   ├── queries/
│   ├── mutations/
│   └── subscriptions/
│
├── tests/
└── playwright/
```

### What each top-level area owns

- **`pages/`** — Next.js route entry points. Pages are thin: they assemble components and wire data. No business logic, no styled components beyond layout glue.
- **`components/shared/`** — Cross-cutting primitives that any domain may use. If two domain folders both need it, it belongs here. The canonical home for Modal, Card, Stepper, EmptyState, LoadingState, PageHeader, DataTable, FormFields, CodeEditor, and Markdown.
- **`components/layout/`** — Application chrome — the Shell, Header, Navigator, Footer, and NotificationCenter. These are mounted once per app, not per route.
- **`components/<domain>/`** — One folder per product domain (Designs, Workspaces, Connections, Environments, Performance, Dashboard, Adapters, Registry, Settings, User, Telemetry, Extensions, Subscription). Each domain owns its own list, detail, and form views.
- **`theme/`** — The single entrypoint for theme utilities. Always import styling utilities (`styled`, `useTheme`, `alpha`, `lighten`, `darken`, and palette accessors) from `@/theme`, rather than importing them directly from `@sistent/sistent` or `@mui/*` in application code. UI components themselves still come from `@sistent/sistent` per principle 1. See [`THEMING.md`](./THEMING.md) for the full theming guide.
- **`utils/`** — Pure functions only. No JSX, no React components. Hooks live in `utils/hooks/`. Generic helpers (`format/`, `k8s/`, etc.) live in subfolders by category.
- **`assets/`** — Static-only, plus typed SVG icons. Icons, images, GIFs, and schemas. The one allowed `.tsx` flavor is the typed SVG components under `assets/icons/`; any other `.tsx` file in `assets/` is a bug.
- **`graphql/`** — All GraphQL queries, mutations, and subscriptions. These are not components — they live at the top level, not under `components/`.
- **`store/`, `rtk-query/`, `machines/`** — Already clean. Untouched by the restructure.

### What's going away

- `ui/themes/` — renamed to `ui/theme/` and shrunk to a Sistent re-export.
- `ui/constants/colors.ts` — deleted; callers move to `theme.palette.*`.
- `ui/components/icons/` and `ui/assets/new-icons/` — merged into `ui/assets/icons/` with consistent naming.
- `ui/components/hooks/` — merged into `ui/utils/hooks/`.
- `ui/components/General/` — split into `shared/` (primitives) and `layout/` (chrome).
- `ui/components/shapes/` — merged into `ui/assets/icons/shapes/`.
- All sibling `*.style(s).tsx` files — colocated inside their component folder with unique names (`Header.styled.ts`, not `styles.tsx`).

---

## Where do I put my code?

A quick decision guide. When in doubt, follow the rule below; when the rule disagrees with the existing code, the rule wins for new code and the existing code is a candidate for a follow-up cleanup.

### Building a UI component

- **Need a Modal?** Use `components/shared/Modal/` — never write a new one. If you need a variant, add it inside `shared/Modal/` (e.g. a new flavor alongside Confirm/Info/Export).
- **Need a Card?** Use `components/shared/Card/`. If your card has domain-specific content, compose it inside your domain folder using the shared `Card` primitive.
- **Need a Stepper, EmptyState, LoadingState, PageHeader, or DataTable?** Same rule: reach into `shared/` first.
- **Building something specific to a domain (e.g. a Designs deploy form)?** It lives in that domain folder: `components/designs/deployment/`.
- **Building app chrome (a new piece of the shell, a global notification)?** It lives in `components/layout/`.

### Styling

- **Theming a new component?** Import from `@/theme`, not from `@mui/*`, not from `@sistent/sistent` directly in application code. The `@/theme` indirection keeps the codebase forward-compatible as Sistent evolves.
- **Need a color?** Use `theme.palette.*` inside a `styled()` callback or via `useTheme()`. No hex, no rgb, no named-color strings.
- **Need spacing?** Use `theme.spacing(n)`. No raw `'8px'`, no `'1rem'`.
- **Need a breakpoint?** Use `theme.breakpoints.up('md')` etc. inside styled callbacks.
- **Tempted to write `style={{ color: ... }}`?** Refactor to a `styled()` component first. Inline `style` is reserved for values that genuinely change per render (drag transforms, dynamic positioning).

### File size

- **File approaching 400 lines?** Plan a refactor. Extract a hook, lift a sub-component, move types out.
- **File past 600 lines?** A lint warning will fire and review will block the PR. Refactor before adding more.
- **File past 1000 lines?** CI hard-blocks the merge. Split the file as part of the same PR.

### Hooks, utilities, contexts

- **Adding a generic hook (`useDebounce`, `useLocalStorage`)?** It lives in `utils/hooks/`, not `components/hooks/`.
- **Adding a domain-specific hook (`useDesignLifecycle`)?** It lives in its domain folder: `components/designs/lifecycle/useDesignLifecycle.ts`.
- **Adding a pure formatting helper?** `utils/format/`.
- **Adding a React context?** `utils/context/` for cross-cutting contexts; otherwise colocate with its consumer.

### Icons and assets

- **Need an icon?** Check `@sistent/sistent` first. If absent, add a typed SVG component to `ui/assets/icons/` and export it from `ui/assets/icons/index.ts`. Do not import from `@mui/icons-material`.
- **Adding an image, GIF, or schema?** It lives under `ui/assets/`, never under `ui/components/`.

### GraphQL

- **Writing a new query, mutation, or subscription?** It lives under `ui/graphql/`, not under `ui/components/graphql/`. Components import from `ui/graphql/`.

### Filenames

- **Naming a new component file?** Name it for what it contains — `DesignBuilder.tsx`, not `index.tsx` (unless `index.tsx` is a folder's public re-export alongside other named files).
- **Splitting styled-components out?** Use `<Component>.styled.ts`, not `styles.tsx`.
- **Splitting types out?** Use `<Component>.types.ts`.

---

## What this document is intentionally NOT

To keep this document honest and the restructure scope contained, the following are out of scope for both this doc and the restructure itself:

- **Not a rewrite.** Existing components are migrated, not rebuilt. The plan is consolidation, not greenfield design.
- **Not a framework swap.** Next.js, RTK, RTK Query, XState, Relay, and the rest stay as-is. The pages router setup is untouched.
- **Not a Sistent version upgrade.** If `@sistent/sistent` needs a new release, that is a separate coordinated upgrade, not part of this restructure.
- **Not a TypeScript strictness push.** Tempting, but orthogonal — handled separately.
- **Not a test-coverage campaign.** Also orthogonal; the current Vitest setup is fine.
- **Not a Sistent token rename.** We consume what Sistent exports; we don't rename its tokens locally.
- **Not coordinated downtime.** Every phase ships in small, independent PRs. Nothing in the restructure requires a freeze.

If your change touches any of the above, it belongs in a separate issue.

---

## Where Phase 1 fits

The restructure is split into five phases, each shipping as small PRs rather than one monolith:

- **Phase 1 — Guardrails.** ESLint rules, CI audit scripts, and the docs you are reading right now (`ARCHITECTURE.md` and `THEMING.md`). Adds the `ui/theme/index.ts` re-export but deletes nothing. **No runtime behavior change.**
- **Phase 2 — Theme & Color Collapse.** Migrate the three color tables and ~1,500 hex literals into `theme.palette.*`. Delete `themes/app.ts`, `themes/index.ts`, `constants/colors.ts`.
- **Phase 3 — MUI Elimination.** Replace all 100 `@mui/*` imports with Sistent. Drop `@mui/*` and `@rjsf/mui` from `package.json` direct deps.
- **Phase 4 — Folder Reshape.** Move folders into the target structure above, one `git mv` PR at a time so history is preserved.
- **Phase 5 — Break the Giants & Dedupe.** Split the eight 1,000+ line files; dedupe the 22 modals (and similar) into one shared primitive each.

This document is part of Phase 1. The directory model it describes is the **target**; the current tree is in the middle of being moved into it. When existing code disagrees with this document, the document wins for new code, and the existing code is a follow-up cleanup candidate.

---

## Related docs and references

- [`../restructure-plan.md`](../restructure-plan.md) — the full long-form restructure plan, including the current-state inventory, ESLint rules, deduplication tables, and the per-file split plan for the giant files.
- [`./THEMING.md`](./THEMING.md) — the theming guide, covering tokens, palette accessors, dark mode, and the `@/theme` entrypoint. (Landing in parallel with this document.)
- Parent epic: [meshery/meshery#18656](https://github.com/meshery/meshery/issues/18656) — the umbrella issue tracking all five phases.
- Phase 1 sub-issue (this doc): [meshery/meshery#18723](https://github.com/meshery/meshery/issues/18723).
- Contributor guide: [Contributing to Meshery UI](https://docs.meshery.io/project/contributing/contributing-ui) — environment setup, build commands, and PR conventions.
