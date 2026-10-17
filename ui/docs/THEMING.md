# Theming

This document is the contributor-facing guide to styling Meshery UI components.

The short version:

- Theme primitives (colors, spacing, breakpoints, the `styled()` factory, the
  `useTheme()` hook) come from `@sistent/sistent`, the Meshery design system.
- The project-local front door is `@/theme`. Import everything theme-related
  from there, not from `@sistent/sistent` directly.
- Hardcoded colors (hex literals, `rgb()`, `rgba()`) and inline `style={{}}`
  props for static values are banned in component code. Use `theme.palette.*`
  via `styled()` instead.

This guide is part of the Phase 1 UI restructure tracked in the parent epic
[meshery/meshery#18656](https://github.com/meshery/meshery/issues/18656). For
the full migration plan, see [`../restructure-plan.md`](../restructure-plan.md),
especially §4 ("Theme & Color Consolidation") and §8 ("Lint rules").

---

## The `@/theme` entry point

`ui/theme/index.ts` is a thin wrapper around `@sistent/sistent`. It exists so
every consumer in the Meshery UI imports theme primitives from a single,
project-local path. The wrapper is deliberately small: re-exports plus a
handful of named accessors that read directly from `theme.palette.*`.

### What it re-exports today

From `@sistent/sistent`:

- **Hooks** — `useTheme`
- **Styling** — `styled`
- **Color helpers** — `alpha`, `lighten`
- **Providers and global primitives** — `SistentThemeProvider`,
  `SistentThemeProviderWithoutBaseLine`, `CssBaseline`, `NoSsr`
- **Types** — `Theme`

Bridged from MUI (until Sistent re-exports them) so callers can still go
through the project-local front door:

- **Color helpers** — `darken`
- **Global primitives** — `GlobalStyles` (used for one-off escape hatches
  like cross-portal z-index overrides; routed through `@/theme` so app code
  stays off `@mui/material` directly)

### What it adds locally

A `palette` object with named accessors for the palette paths the app reaches
for most often:

- `palette.status.error`, `palette.status.warning`, `palette.status.success`,
  `palette.status.info` — return `theme.palette.<status>.main`.
- `palette.surface.page`, `palette.surface.elevated`, `palette.surface.card`
  — return the corresponding `theme.palette.background.*` token.

Each accessor takes a `Theme` argument:

```tsx
import { styled, palette } from '@/theme';

const ErrorBadge = styled('span')(({ theme }) => ({
  color: palette.status.error(theme),
  background: palette.surface.card(theme),
}));
```

### How to import

Prefer `@/theme` as the front door. Never deep-import `@/theme/index`, and
avoid importing from `@sistent/sistent` for primitives that `@/theme`
already re-exports — `@/theme` is the project-local entry point that future
phases will extend with additional palette accessors and helpers.

```tsx
// Correct
import { useTheme, styled, alpha, palette } from '@/theme';

// Incorrect — deep import of the local entry point (lint-blocked)
import { styled } from '@/theme/index';

// Discouraged — bypasses the project-local front door (convention, not
// lint-blocked today)
import { styled } from '@sistent/sistent';
```

Today's `no-restricted-imports` rule in
[`ui/eslint.config.js`](../eslint.config.js) blocks the `@/theme/index`
deep import and direct `@mui/*` / `@material-ui/*` / `@rjsf/mui` imports.
Routing
theme primitives through `@/theme` rather than `@sistent/sistent` is
project convention; a later phase may tighten the rule.

---

## Token-based theming

All static styling values should come from the theme:

- **Colors** — `theme.palette.*` (e.g. `theme.palette.text.primary`,
  `theme.palette.primary.main`, `theme.palette.error.main`).
- **Spacing** — `theme.spacing(n)` where `n` is a multiplier of the design
  system's base unit. `theme.spacing(2)` is preferred over `'16px'`.
- **Breakpoints** — `theme.breakpoints.up('md')`, `theme.breakpoints.down('sm')`,
  etc.

### Right way

```tsx
import { styled } from '@/theme';

const StyledCard = styled('div')(({ theme }) => ({
  color: theme.palette.text.primary,
  background: theme.palette.background.card,
  padding: theme.spacing(2),
  [theme.breakpoints.up('md')]: {
    padding: theme.spacing(4),
  },
}));
```

```tsx
import { styled, alpha } from '@/theme';

const HoverOverlay = styled('div')(({ theme }) => ({
  background: alpha(theme.palette.primary.main, 0.08),
  borderColor: theme.palette.divider,
}));
```

### Wrong way

```tsx
// Don't: hex literals, hard-coded pixels, no theme awareness
<div style={{ color: '#1E1E1E', padding: '16px', background: '#FFFFFF' }}>
```

The wrong example will:

- Break dark mode (the literal `#FFFFFF` doesn't react to theme switches).
- Trip the `no-restricted-syntax` lint rule (hex literals are forbidden
  outside `ui/theme/`).
- Trip the `react/forbid-dom-props` rule on `style` (see below).

---

## The ban on color literals

Hex literals (`#RRGGBB`, `#RGB`, with optional alpha) and `rgb()`/`rgba()`
function-call strings are forbidden in component code. The rule is enforced
by `no-restricted-syntax` in
[`ui/eslint.config.js`](../eslint.config.js).

### What's allowed, and where

Literal colors are only allowed in modules that define the theme or ship
non-themed assets:

- `ui/theme/**` — the theme module itself.
- `ui/assets/**` — SVG icons encoded as React components.
- `ui/lib/**` — third-party integration helpers.
- `ui/public/**` — static assets.

Everywhere else, use `theme.palette.*` (composed if needed with `alpha` or
`lighten` from `@/theme`).

### Current mode and the path to `error`

The rule ships in **warn** mode with a file-level allowlist for legacy
offenders that haven't been migrated yet (the `legacyLiteralColorOffenders`
list in `eslint.config.js`). The plan is to:

1. Drain the allowlist file-by-file as components are touched.
2. Promote the rule from `warn` to `error` once the list is empty.

When you touch a file on the allowlist, please remove it from the list as
part of your change.

---

## The ban on inline `style` props

`react/forbid-dom-props` forbids the `style` prop in components for the
same reason: static styling belongs in `styled()`, not on the element.

### What inline `style` is reserved for

**Truly dynamic geometry that can't be expressed in CSS-in-JS at definition
time.** Examples:

```tsx
// Draggable element: x/y change at 60fps, can't go through styled().
<div style={{ transform: `translate(${x}px, ${y}px)` }} />

// Resize observer driven layout.
<div style={{ width: measuredWidth }} />

// CSS variable driving an animation.
<div style={{ '--progress': progress } as React.CSSProperties} />
```

### What inline `style` is not for

Static colors, paddings, dimensions, borders, font weights. All of these
go in `styled()`.

### Before and after

Before:

```tsx
function StatusPill({ status, label }: Props) {
  return (
    <span
      style={{
        color: status === 'error' ? '#F91313' : '#1E1E1E',
        backgroundColor: '#F5F5F5',
        padding: '4px 8px',
        borderRadius: '12px',
      }}
    >
      {label}
    </span>
  );
}
```

After:

```tsx
import { styled, palette } from '@/theme';

const Pill = styled('span', {
  shouldForwardProp: (prop) => prop !== 'isError',
})<{ isError: boolean }>(({ theme, isError }) => ({
  color: isError ? palette.status.error(theme) : theme.palette.text.primary,
  backgroundColor: palette.surface.elevated(theme),
  padding: theme.spacing(0.5, 1),
  borderRadius: theme.shape.borderRadius,
}));

function StatusPill({ status, label }: Props) {
  return <Pill isError={status === 'error'}>{label}</Pill>;
}
```

The `palette.status.*` and `palette.surface.*` accessors are the
project-local helpers from `@/theme`; they read directly from
`theme.palette.*` but document the supported palette paths in one place.

Like the hex-literal rule, this rule ships in `warn` mode with a file-level
allowlist; future phases drain the allowlist and promote it to `error`.

---

## When to use `styled()` vs `useTheme()` vs inline `style`

Use this decision rule when writing or refactoring a component:

- **`styled()`** — the default for new components and for any component
  whose style does not depend on a value computed at render time. Static,
  theme-aware, and the resulting components are reusable.
- **`useTheme()`** — reach for this only when you need to read a theme
  value at render time and there's no way to express the dependency through
  `styled()` props. This is rare; if you find yourself here, double-check
  whether a prop on a `styled()` component would do the job.
- **Inline `style={{}}`** — only for truly dynamic geometry (the `transform`,
  `width`, CSS-variable cases listed above). Never for colors or static
  spacing.

In practice, 90%+ of components in this codebase should be `styled()`.

---

## Legacy imports → new entry point

A migration table for the legacy modules that the lint rules and the
restructure plan call out for deletion:

| Legacy import                            | Replace with                             |
| ---------------------------------------- | ---------------------------------------- |
| `import { styled } from '@/theme/index'` | `import { styled } from '@/theme'`       |
| `import { ... } from '@mui/material'`    | `import { ... } from '@sistent/sistent'` |

For the underlying ESLint configuration that enforces these mappings, see
the `no-restricted-imports` block in
[`ui/eslint.config.js`](../eslint.config.js).

Once a file is migrated off a legacy import, it should also be removed from
the `legacyRestrictedImportOffenders` allowlist in the same file.

---

## When to add a token upstream to Sistent

If a color, spacing scale, typography setting, or other token you need is
missing from Sistent's palette, **open an upstream PR to
[Sistent](https://npmjs.com/package/@sistent/sistent) rather than overriding
locally.**

Local overrides are how the sprawl regrew last time. The whole point of
consolidating on Sistent is to make the design system the single source of
truth; that only works if missing tokens go back upstream.

Examples of "open it upstream":

- A new status color (e.g. an additional "info" variant) that the app needs
  in multiple places.
- A typography scale step the design team has agreed on but isn't yet in
  Sistent.
- A border-radius or elevation token used by a shared component.

Examples of "no upstream PR needed":

- A one-off pixel offset for a specific animation frame.
- A `transform` value driven by component state.

When in doubt, ask in `#meshery-ui` whether the token is reusable enough to
belong upstream.

---

## RJSF theme

`ui/theme/rjsf.ts` currently builds plain MUI themes with `createTheme`
(imported via its `@sistent/sistent` re-export) to configure
the [react-jsonschema-form](https://github.com/rjsf-team/react-jsonschema-form)
adapter. A later phase (see §4.4 of [`../restructure-plan.md`](../restructure-plan.md))
migrates it to a Sistent-backed theme via `extendSistentTheme`.

Until then:

- `@rjsf/mui` remains as a transitive dependency only. The shared
  `RJSFProvider` wrapper is the only place in the codebase that's allowed
  to import from `@rjsf/mui` directly — the lint rule
  (`no-restricted-imports`) blocks the import everywhere else.
- App code should consume the shared RJSF wrapper, not `@rjsf/mui` and not
  `ui/theme/rjsf.ts` directly.

---

## Related docs and references

- [`../restructure-plan.md`](../restructure-plan.md) — the full Phase 1 plan.
  Sections §4 ("Theme & Color Consolidation") and §8 ("Lint rules") are the
  most relevant to this document.
- [`./ARCHITECTURE.md`](./ARCHITECTURE.md) — companion architecture doc
  describing target file layout, component conventions, and the broader
  Sistent migration. Created alongside this doc in Phase 1.
- Parent epic — [meshery/meshery#18656](https://github.com/meshery/meshery/issues/18656).
- Sistent design system — [@sistent/sistent](https://npmjs.com/package/@sistent/sistent) on npm.
- Lint rules that enforce the conventions in this document live in
  [`ui/eslint.config.js`](../eslint.config.js)
  (`no-restricted-imports`, `no-restricted-syntax`, and the
  `react/forbid-dom-props` configuration).
