import React from 'react';
import type { IconProps } from './types';

/**
 * ArrowDropDownIcon — typed SVG replacement for `@mui/icons-material/ArrowDropDown`.
 *
 * Material Icons triangular-caret glyph, viewBox 0 0 24 24.
 */
export const ArrowDropDownIcon: React.FC<IconProps> = ({
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
    <path d="M7 10l5 5 5-5z" />
  </svg>
);

export default ArrowDropDownIcon;
