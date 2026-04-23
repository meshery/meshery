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

/**
 * Modernization Brand Colors
 * Centralized here as the approved location for color constants.
 */
export const BrandColors = {
  KEPPEL_GREEN: '#66C2A5',
  AMBER: '#FFD35C',
  RED: '#E75248',
  BLUE: '#3C494F', // Meshery Blue
  FACEBOOK: '#1877F2',
  TWITTER: '#1DA1F2',
  LINKEDIN: '#0077B5',
  GITHUB: '#333333',
  YOUTUBE: '#FF0000',
  SLACK: '#4A154B',
  DOCKER: '#2496ED',
} as const;

/**
 * JSON Viewer Theme Tokens (Modernized)
 */
export const RE_JSON_Theme = {
  BASE00: '#1e2117',
  BASE01: '#323232',
  BASE02: '#353535',
  BASE03: '#444',
  BASE04: '#b0b0b0',
  BASE05: '#d0d0d0',
  BASE06: '#e0e0e0',
  BASE07: '#f5f5f5',
  BASE08: '#66C2A5', // Keppel
  BASE09: '#7BA05B', // Olivine
  BASE0A: '#FFD35C', // Amber
  BASE0B: '#99CC99',
  BASE0C: '#66CCCC',
  BASE0D: '#6699CC',
  BASE0E: '#CC99CC',
  BASE0F: '#D27B53',
} as const;

/**
 * Legacy Colors (Moved from themes/app.ts to comply with linting)
 */
export const Colors = {
  darkJungleGreen: '#1E2117',
  caribbeanGreen: '#00D3a9',
  keppelGreen: '#00B39F',
  charcoal: '#3C494F',
  blueGrey: '#607D8B',
} as const;

export const CHART_COLORS = [
  '#14232A', // Gunmetal
  '#2E5261',
  '#3B697D', // Paynes' Gray
  '#4A839C',
  '#5996B1',
  '#74A8BE',
  '#90B9CB',
  '#CBDEE6',
] as const;

export const notificationColors = {
  error: '#F91313',
  warning: '#F0A303',
  success: '#206D24',
  info: '#2196F3',
  darkRed: '#B32700',
  lightwarning: '#E75225',
} as const;

export const darkNotificationColors = {
  error: '#F91313',
  warning: '#F0D053',
  success: '#78C57C',
  info: '#5FD4FF',
} as const;

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
