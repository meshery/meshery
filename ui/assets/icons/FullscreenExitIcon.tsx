import React from 'react';
import type { IconProps } from './types';

/**
 * FullscreenExitIcon — typed SVG replacement for `@mui/icons-material/FullscreenExit`.
 *
 * Material Icons exit-fullscreen glyph, viewBox 0 0 24 24.
 */
export const FullscreenExitIcon: React.FC<IconProps> = ({
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
    <path d="M5 16h3v3h2v-5H5v2zm3-8H5v2h5V5H8v3zm6 11h2v-3h3v-2h-5v5zm2-11V5h-2v5h5V8h-3z" />
  </svg>
);

export default FullscreenExitIcon;
