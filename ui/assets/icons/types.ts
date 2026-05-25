import type { SVGProps } from 'react';

/**
 * Shared icon prop contract for typed SVG icons in `ui/assets/icons/`.
 *
 * Each icon accepts standard SVG element props plus convenience overrides
 * for size and color. `fill` defaults to `currentColor` so that the icon
 * inherits color from CSS (`color`) by default — matching MUI icon behavior
 * and Sistent icon conventions.
 */
export type IconProps = Omit<SVGProps<SVGSVGElement>, 'width' | 'height' | 'fill'> & {
  width?: number | string;
  height?: number | string;
  fill?: string;
};
