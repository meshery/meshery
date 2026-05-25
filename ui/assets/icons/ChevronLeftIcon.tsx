import React from 'react';
import type { IconProps } from './types';

/**
 * ChevronLeftIcon — typed SVG replacement for `@mui/icons-material/ChevronLeft`.
 *
 * Material Icons single-stroke chevron pointing left, viewBox 0 0 24 24.
 * Note: this is distinct from the existing `LeftArrowIcon` (thicker
 * double-stroke style) and from the unrelated `chevron_right.tsx` glyph.
 */
export const ChevronLeftIcon: React.FC<IconProps> = ({
  width = 24,
  height = 24,
  fill = 'currentColor',
  ...props
}) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={width}
    height={height}
    viewBox="0 0 24 24"
    fill={fill}
    {...props}
  >
    <path d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z" />
  </svg>
);

export default ChevronLeftIcon;
