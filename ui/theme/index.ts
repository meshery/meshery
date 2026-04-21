import type { Theme } from '@sistent/sistent';

/**
 * Meshery UI theme entry point.
 *
 * This module is a thin wrapper around {@link https://github.com/layer5io/sistent Sistent},
 * the Meshery design system. It exists so that every consumer in the
 * Meshery UI imports theme primitives (`useTheme`, `styled`, `alpha`,
 * `lighten`, ...) from a single, project-local path.
 *
 * Phase 1 treats `@/theme` as the approved, project-local import path for
 * theme primitives. Prefer it over importing those primitives directly from
 * Sistent or from legacy `@/themes*` modules.
 *
 *   import { useTheme, styled, alpha } from '@/theme';
 *
 * Preferred theming conventions for new or touched UI code are:
 *   - Prefer colors from `theme.palette.*` instead of introducing new hex literals.
 *   - Prefer spacing from `theme.spacing()` instead of hard-coded pixel values where possible.
 *   - Prefer breakpoints from `theme.breakpoints.*`.
 *
 * Some legacy theme files in the UI still contain literal color and spacing
 * values; treat those as migration candidates rather than precedent for new code.
 *
 * `@/themes/hooks` remains for theme-preference plumbing until a later phase,
 * but new theme-entrypoint imports should start from `@/theme`.
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
