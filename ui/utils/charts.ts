import { darken, lighten, type Theme } from '@/theme';

export const isValidColumnName = (name) => {
  return name !== '' && name !== ' ' && name != undefined && name != null;
};

/**
 * Donut/bar chart slice palette, derived from `theme.palette` so it adapts
 * to light/dark mode. Callers must thread the active `theme` through because
 * this `.ts` module cannot call `useTheme()` directly.
 */
export const getChartColors = (theme: Theme): string[] => {
  const primary = theme.palette.primary.main;
  const info = theme.palette.info.main;
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

export const dataToColors = (data: unknown[][], theme: Theme): Record<string, string> => {
  const palette = getChartColors(theme);
  const colors: Record<string, string> = {};
  data.forEach((item, idx) => {
    colors[String(item[0])] = palette[idx % palette.length];
  });
  return colors;
};
