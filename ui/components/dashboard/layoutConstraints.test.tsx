import { describe, expect, it } from 'vitest';
import { applyMinSizeConstraints } from './layoutConstraints';

const cols = { lg: 12, md: 10 };

describe('applyMinSizeConstraints', () => {
  it('returns an empty object when given empty layouts', () => {
    expect(applyMinSizeConstraints({}, {}, cols)).toEqual({});
    expect(applyMinSizeConstraints(undefined as never, {}, cols)).toEqual({});
  });

  it('clamps item w/h to be at least the default item size when smaller', () => {
    const layouts = {
      lg: [{ i: 'OVERVIEW', x: 0, y: 0, w: 1, h: 1 }],
    };
    const defaults = {
      lg: [{ i: 'OVERVIEW', x: 0, y: 0, w: 4, h: 2 }],
    };

    const result = applyMinSizeConstraints(layouts, defaults, cols);
    expect(result.lg[0].w).toBe(4);
    expect(result.lg[0].h).toBe(2);
    expect(result.lg[0].minW).toBe(4);
    expect(result.lg[0].minH).toBe(2);
  });

  it('respects existing larger sizes', () => {
    const layouts = {
      lg: [{ i: 'OVERVIEW', x: 0, y: 0, w: 8, h: 5 }],
    };
    const defaults = {
      lg: [{ i: 'OVERVIEW', x: 0, y: 0, w: 4, h: 2 }],
    };

    const result = applyMinSizeConstraints(layouts, defaults, cols);
    expect(result.lg[0].w).toBe(8);
    expect(result.lg[0].h).toBe(5);
  });

  it('caps minW at the breakpoint column count', () => {
    const layouts = {
      lg: [{ i: 'WIDE', x: 0, y: 0, w: 1, h: 1 }],
    };
    const defaults = {
      lg: [{ i: 'WIDE', x: 0, y: 0, w: 50, h: 2 }],
    };

    const result = applyMinSizeConstraints(layouts, defaults, cols);
    expect(result.lg[0].minW).toBe(12);
    expect(result.lg[0].w).toBe(12);
  });

  it('falls back to widget sizing when no default item matches', () => {
    const layouts = {
      lg: [{ i: 'NEW_WIDGET', x: 0, y: 0, w: 1, h: 1 }],
    };
    const widgetSizing = {
      NEW_WIDGET: { w: 3, h: 2 },
    };

    const result = applyMinSizeConstraints(layouts, {}, cols, widgetSizing);
    expect(result.lg[0].minW).toBe(3);
    expect(result.lg[0].minH).toBe(2);
  });

  it('defaults to 1x1 minimum when no defaults or sizing are provided', () => {
    const layouts = {
      lg: [{ i: 'UNKNOWN', x: 0, y: 0, w: 0, h: 0 }],
    };
    const result = applyMinSizeConstraints(layouts, {}, cols);
    expect(result.lg[0].minW).toBe(1);
    expect(result.lg[0].minH).toBe(1);
    expect(result.lg[0].w).toBe(1);
    expect(result.lg[0].h).toBe(1);
  });

  it('defaults breakpoint cols to 12 when not present in the cols config', () => {
    const layouts = {
      xl: [{ i: 'WIDE', x: 0, y: 0, w: 1, h: 1 }],
    };
    const defaults = {
      xl: [{ i: 'WIDE', x: 0, y: 0, w: 100, h: 1 }],
    };

    const result = applyMinSizeConstraints(layouts, defaults, cols);
    expect(result.xl[0].minW).toBe(12);
  });

  it('preserves additional layout item properties', () => {
    const layouts = {
      lg: [{ i: 'OVERVIEW', x: 4, y: 3, w: 6, h: 4, moved: true, static: true }],
    };
    const result = applyMinSizeConstraints(layouts, {}, cols);
    expect(result.lg[0]).toMatchObject({
      i: 'OVERVIEW',
      x: 4,
      y: 3,
      moved: true,
      static: true,
    });
  });
});
