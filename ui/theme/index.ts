/**
 * Meshery UI theme entry point.
 *
 * This module is a thin re‑export from {@link https://github.com/layer5io/sistent Sistent},
 * the Meshery design system. It exists so that every consumer in the
 * Meshery UI imports theme primitives (`useTheme`, `styled`, `alpha`,
 * `lighten`, `darken`, …) from a single, project‑local path.
 *
 *   import { useTheme, styled, alpha } from '@/theme';
 *
 * The core theming rules are:
 *   - Colors come from `theme.palette.*` — never from a hex literal.
 *   - Spacing comes from `theme.spacing()` — never from a hard‑coded pixel.
 *   - Breakpoints come from `theme.breakpoints.*`.
 *
 * If Sistent is missing a token the app needs, open an issue or PR upstream
 * rather than redefining it here. This file must remain a re‑export only.
 */

export {
  // Hooks
  useTheme,

  // CSS-in-JS
  styled,

  // Color helpers
  alpha,
  lighten,
  darken,

  // Providers & global primitives
  SistentThemeProvider,
  SistentThemeProviderWithoutBaseline,
  CssBaseline,
  NoSsr,
} from '@sistent/sistent';
