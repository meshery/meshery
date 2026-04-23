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

export const BrandColors = {
  twitter: '#1da1f2',
  slack: '#4A154B',
  github: '#24292e',
  docker: '#2496ed',
  youtube: '#ff0000',
  linkedin: '#0077b5',
} as const;

export const CHART_COLORS = [
  '#393A3D',
  '#626469',
  '#C1C1C1',
  '#E6E6E6',
  '#222222',
  '#E0E0E0',
  '#BDBDBD',
];

export const notificationColors = {
  error: '#F24C3D',
  warning: '#F3AA60',
  info: '#22B14C',
  success: '#22B14C',
};

export const darkNotificationColors = {
  error: '#F24C3D',
  warning: '#F3AA60',
  info: '#22B14C',
  success: '#22B14C',
};

export const RE_JSON_Theme = (themeType: string) => ({
  base00: themeType === 'dark' ? '#303030' : '#ffffff',
  base01: '#444c56',
  base02: themeType === 'dark' ? '#586069' : '#abb2bf',
  base03: '#6a737d',
  base04: '#477E96',
  base05: '#9ea7a6',
  base06: '#d8dee9',
  base07: themeType === 'dark' ? '#FFF3C5' : '#002B36',
  base08: '#2a5491',
  base09: '#d19a66',
  base0A: '#EBC017',
  base0B: '#237986',
  base0C: '#56b6c2',
  base0D: '#B1B6B8',
  base0E: '#e1e6cf',
  base0F: '#647881',
});

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
