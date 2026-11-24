import { darken, lighten, type Theme } from '@/theme';

export const isValidColumnName = (name: unknown): boolean => {
  if (typeof name !== 'string') {
    return false;
  }
  return name.trim() !== '';
};

/**
 * Donut/bar chart slice palette, derived from `theme.palette` so it adapts
 * to light/dark mode. Callers must thread the active `theme` through because
 * this `.ts` module cannot call `useTheme()` directly.
 */
export const getChartColors = (theme?: Theme | null): string[] => {
  const primary = theme?.palette?.primary?.main;
  const info = theme?.palette?.info?.main;
  if (!primary || !info) {
    return [];
  }
  return [
    darken(primary, 0.6),
    darken(primary, 0.35),
    darken(primary, 0.15),
    primary,
    lighten(primary, 0.2),
    lighten(info, 0.1),
    lighten(info, 0.35),
    lighten(info, 0.65),
  ];
};

export const dataToColors = (
  data?: unknown[][] | null,
  theme?: Theme | null,
): Record<string, string> => {
  if (!Array.isArray(data) || data.length === 0) {
    return {};
  }
  const palette = getChartColors(theme);
  if (palette.length === 0) {
    return {};
  }
  const colors: Record<string, string> = {};
  let colorIdx = 0;
  data.forEach((item) => {
    if (Array.isArray(item) && item.length > 0 && item[0] !== undefined && item[0] !== null) {
      colors[String(item[0])] = palette[colorIdx++ % palette.length];
    }
  });
  return colors;
};
