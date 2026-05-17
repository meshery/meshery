import React from 'react';
import type { IconProps } from './types';

/**
 * ChevronRightIcon — typed SVG replacement for `@mui/icons-material/ChevronRight`.
 *
 * Material Icons single-stroke chevron pointing right, viewBox 0 0 24 24.
 * Note: this is distinct from the existing `RightArrowIcon` (thicker
 * double-stroke) and from the unrelated `chevron_right.tsx` glyph.
 */
export const ChevronRightIcon: React.FC<IconProps> = ({
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
    <path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z" />
  </svg>
);

export default ChevronRightIcon;
