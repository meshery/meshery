import React from 'react';
import type { IconProps } from './types';

/**
 * FullscreenIcon — typed SVG replacement for `@mui/icons-material/Fullscreen`.
 *
 * Material Icons enter-fullscreen glyph, viewBox 0 0 24 24.
 */
export const FullscreenIcon: React.FC<IconProps> = ({
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
    <path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z" />
  </svg>
);

export default FullscreenIcon;
