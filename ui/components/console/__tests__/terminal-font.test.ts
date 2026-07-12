import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

/**
 * Widths a fake platform reports per family. `i` and `W` differ for a
 * proportional face and match for a fixed-pitch one.
 */
type Advances = Record<string, { i: number; W: number }>;

/**
 * Stands in for a browser + fontconfig. The `substitute` map models the trap
 * this module exists for: an uninstalled family silently resolving to some other
 * face rather than falling through to the next candidate.
 */
const installPlatform = (advances: Advances, substitute: Record<string, string> = {}) => {
  let current = '';

  const context = {
    set font(value: string) {
      current = value;
    },
    get font() {
      return current;
    },
    measureText(text: string) {
      const family = current.replace(/^\d+px\s*/, '');
      const resolvedFamily = substitute[family] ?? family;
      const advance = advances[resolvedFamily];
      if (!advance) return { width: 0 };
      return { width: text === 'W' ? advance.W : advance.i };
    },
  };

  vi.spyOn(document, 'createElement').mockReturnValue({
    getContext: () => context,
  } as unknown as HTMLCanvasElement);
};

/** The module caches its answer, so each case needs a fresh import. */
const loadResolver = async () => {
  const fresh = await import('../terminal-font');
  return fresh.resolveTerminalFontFamily;
};

beforeEach(() => {
  vi.resetModules();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('resolveTerminalFontFamily', () => {
  it('rejects a family the platform substitutes with a proportional face', async () => {
    // Exactly the observed bug: "JetBrains Mono" is not installed, fontconfig
    // hands back Noto Sans, and every glyph advance differs.
    installPlatform(
      {
        'Noto Sans': { i: 4, W: 12 },
        monospace: { i: 8, W: 8 },
      },
      {
        '"JetBrains Mono"': 'Noto Sans',
        '"Cascadia Mono"': 'Noto Sans',
        '"SF Mono"': 'Noto Sans',
        Menlo: 'Noto Sans',
        Consolas: 'Noto Sans',
        '"Roboto Mono"': 'Noto Sans',
        '"DejaVu Sans Mono"': 'Noto Sans',
        '"Liberation Mono"': 'Noto Sans',
        '"Noto Sans Mono"': 'Noto Sans',
      },
    );

    expect(await (await loadResolver())()).toBe('monospace');
  });

  it('keeps the first genuinely fixed-pitch family', async () => {
    installPlatform({
      '"JetBrains Mono"': { i: 9, W: 9 },
      monospace: { i: 8, W: 8 },
    });

    expect(await (await loadResolver())()).toBe('"JetBrains Mono"');
  });

  it('skips substituted families and settles on the first real monospace', async () => {
    installPlatform(
      {
        'Noto Sans': { i: 4, W: 12 },
        'DejaVu Sans Mono': { i: 8, W: 8 },
        monospace: { i: 8, W: 8 },
      },
      {
        '"JetBrains Mono"': 'Noto Sans',
        '"Cascadia Mono"': 'Noto Sans',
        '"SF Mono"': 'Noto Sans',
        Menlo: 'Noto Sans',
        Consolas: 'Noto Sans',
        '"Roboto Mono"': 'Noto Sans',
        '"DejaVu Sans Mono"': 'DejaVu Sans Mono',
      },
    );

    expect(await (await loadResolver())()).toBe('"DejaVu Sans Mono"');
  });

  it('tolerates sub-pixel rounding rather than rejecting a fixed-pitch face', async () => {
    installPlatform({
      '"JetBrains Mono"': { i: 8, W: 8.02 },
      monospace: { i: 8, W: 8 },
    });

    expect(await (await loadResolver())()).toBe('"JetBrains Mono"');
  });

  it('rejects a family that reports zero-width glyphs', async () => {
    installPlatform({ monospace: { i: 8, W: 8 } });

    // Every named candidate is unknown to the fake platform, so measureText
    // reports 0 and only the generic keyword survives.
    expect(await (await loadResolver())()).toBe('monospace');
  });

  it('falls back to the generic keyword without a canvas', async () => {
    vi.spyOn(document, 'createElement').mockReturnValue({
      getContext: () => null,
    } as unknown as HTMLCanvasElement);

    expect(await (await loadResolver())()).toBe('monospace');
  });
});
