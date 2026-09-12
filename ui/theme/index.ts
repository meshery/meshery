import type { Theme } from '@sistent/sistent';

/**
 * Meshery UI theme entry point.
 *
 * This module is a thin wrapper around Sistent,
 * the Meshery design system. It exists so that every consumer in the
 * Meshery UI imports theme primitives (`useTheme`, `styled`, `alpha`,
 * `lighten`, ...) from a single, project-local path.
 *
 * Phase 1 treats `@/theme` as the approved, project-local import path for
 * theme primitives. Prefer it over importing those primitives directly from
 * Sistent.
 *
 *   import { useTheme, styled, alpha } from '@/theme';
 *
 * Preferred theming conventions for new or touched UI code are:
 *   - Prefer colors from `theme.palette.*` instead of introducing new hex literals.
 *   - Prefer spacing from `theme.spacing()` instead of hard-coded pixel values where possible.
 *   - Prefer breakpoints from `theme.breakpoints.*`.
 *
 * Sibling modules hold theme plumbing that is not part of the entry point:
 * `@/theme/hooks` (theme-preference state), `@/theme/rjsf` (RJSF form themes),
 * and `@/theme/snackbar` (the theme-responsive notification renderer). Treat
 * remaining literal color and spacing values in the UI as migration
 * candidates rather than precedent for new code.
 *
 * If Sistent is missing a token the app needs, open an issue or PR upstream
 * rather than redefining it here. This file must remain a thin wrapper:
 * re-exports plus palette accessors that read directly from `theme.palette.*`.
 */

export {
  // Hooks
  useTheme,

  // CSS-in-JS
  styled,

  // Color helpers
  alpha,
  lighten,

  // Providers & global primitives
  SistentThemeProvider,
  SistentThemeProviderWithoutBaseLine,
  CssBaseline,
  NoSsr,
} from '@sistent/sistent';

// `darken` is not currently re-exported by `@sistent/sistent` (only `lighten`
// is). Until Sistent re-exports it, pull it from the same upstream module
// `lighten` ultimately comes from so callers can still go through `@/theme`
// as the project-local front door.
// eslint-disable-next-line no-restricted-imports
export { darken } from '@mui/material';

// Bridged from MUI: Sistent doesn't yet export GlobalStyles. Routed through
// `@/theme` so the rest of the app can stay off `@mui/material` directly;
// this is the project's approved single bridge.
// eslint-disable-next-line no-restricted-imports
export { GlobalStyles } from '@mui/material';

export type { Theme };

export const palette = {
  status: {
    error: (theme: Theme) => theme.palette.error.main,
    warning: (theme: Theme) => theme.palette.warning.main,
    success: (theme: Theme) => theme.palette.success.main,
    info: (theme: Theme) => theme.palette.info.main,
  },
  surface: {
    page: (theme: Theme) => theme.palette.background.default,
    elevated: (theme: Theme) => theme.palette.background.elevatedComponents,
    card: (theme: Theme) => theme.palette.background.card,
  },
} as const;
