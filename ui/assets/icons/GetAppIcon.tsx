import React from 'react';
import type { IconProps } from './types';

/**
 * GetAppIcon — typed SVG replacement for `@mui/icons-material/GetApp`.
 *
 * Material Icons download arrow glyph, viewBox 0 0 24 24.
 */
export const GetAppIcon: React.FC<IconProps> = ({
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
    <path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z" />
  </svg>
);

export default GetAppIcon;
