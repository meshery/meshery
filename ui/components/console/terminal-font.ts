/**
 * Choosing a font the terminal can actually lay out on.
 *
 * xterm paints on a fixed character grid whose cell width it measures from the
 * font. If the font is not truly fixed-pitch, narrow glyphs (`i`, `l`, `.`, `:`)
 * sit left-aligned in a cell sized for wide ones and words come out visibly
 * gapped mid-word: `maxpr ocs`, `Cor eDNS`, `l i nux`.
 *
 * A hardcoded CSS font stack cannot prevent that, because the CSS fallback chain
 * is not consulted the way one would hope. On Linux, Chrome resolves families
 * through fontconfig, and fontconfig never reports a miss — it returns its best
 * guess. Asking for an uninstalled `"JetBrains Mono"` yields *Noto Sans*, a
 * proportional face, and the rest of the stack is never reached. Leading a stack
 * with a family that happens to be absent therefore breaks the terminal on that
 * machine and nowhere else, which is a miserable bug to chase.
 *
 * So probe instead of guessing: ask the platform to render `i` and `W` in each
 * candidate and keep the first family that gives them the same advance. The
 * generic `monospace` keyword anchors the list and is fixed-pitch by definition,
 * so the probe always terminates on something usable.
 */

/**
 * Candidates in order of preference, ending in the generic keyword. Everything
 * before `monospace` is an upgrade over the system default when it is genuinely
 * installed, and silently skipped when it is not.
 */
const CANDIDATES = [
  '"JetBrains Mono"',
  '"Cascadia Mono"',
  '"SF Mono"',
  'Menlo',
  'Consolas',
  '"Roboto Mono"',
  '"DejaVu Sans Mono"',
  '"Liberation Mono"',
  '"Noto Sans Mono"',
  'monospace',
] as const;

/** The widest and narrowest glyphs a fixed-pitch face must render identically. */
const WIDE = 'W';
const NARROW = 'i';

/** Sub-pixel slack, so a face is not rejected for a rounding difference. */
const TOLERANCE = 0.05;

let resolved: string | undefined;

const isFixedPitch = (context: CanvasRenderingContext2D, family: string): boolean => {
  // A font shorthand the browser cannot parse is ignored, silently leaving the
  // previous value in place — which would make the *previous* family's verdict
  // stand in for this one. Poison it first so a rejected assignment cannot pass.
  context.font = '';
  context.font = `16px ${family}`;
  if (!context.font.includes('16px')) return false;

  const wide = context.measureText(WIDE).width;
  const narrow = context.measureText(NARROW).width;
  if (!wide || !narrow) return false;

  return Math.abs(wide - narrow) <= TOLERANCE;
};

/**
 * The first candidate the platform renders fixed-pitch. Computed once and cached:
 * the answer cannot change while the page is open.
 *
 * Returns the generic keyword when no canvas is available (server render, or a
 * browser that refuses a 2D context), which is the safe answer anyway.
 */
export const resolveTerminalFontFamily = (): string => {
  if (resolved) return resolved;

  const context =
    typeof document === 'undefined' ? null : document.createElement('canvas').getContext('2d');
  if (!context) return 'monospace';

  resolved = CANDIDATES.find((family) => isFixedPitch(context, family)) ?? 'monospace';
  return resolved;
};
