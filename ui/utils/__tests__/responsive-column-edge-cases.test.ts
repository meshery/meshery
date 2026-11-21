import { describe, expect, it } from 'vitest';
import { getResponsiveColumnVisibility, updateVisibleColumns } from '../responsive-column';

describe('responsive column edge cases', () => {
  it('treats prototype-named breakpoints as unsupported', () => {
    expect(updateVisibleColumns([['c', 'constructor']], 0)).toEqual({ c: true });
    expect(updateVisibleColumns([['c', '__proto__']], 0)).toEqual({ c: true });
  });

  it('preserves prototype-named column keys as own properties', () => {
    const result = updateVisibleColumns(
      [
        ['__proto__', 'xs'],
        ['constructor', 'na'],
      ],
      800,
    );

    expect(Object.prototype.hasOwnProperty.call(result, '__proto__')).toBe(true);
    expect(result['__proto__']).toBe(true);
    expect(Object.prototype.hasOwnProperty.call(result, 'constructor')).toBe(true);
    expect(result.constructor).toBe(false);
  });

  it('skips malformed tuples and non-string breakpoints without throwing', () => {
    const entries = [
      null,
      [],
      ['missingSize'],
      ['extra', 'xs', 'value'],
      ['invalid', { toString: null, valueOf: null }],
      ['notString', 12345],
      ['nullSize', null],
      ['valid', 's'],
    ];
    expect(updateVisibleColumns(entries, 700)).toEqual({ valid: true });
  });

  it('rejects invalid column arrays without coercing objects to strings', () => {
    const names = ['valid', { toString: null, valueOf: null }, 123, null, undefined];
    expect(getResponsiveColumnVisibility(names, [['valid', 'xs']], 800)).toEqual({ valid: true });
    expect(getResponsiveColumnVisibility(['valid'], null, 800)).toEqual({ valid: undefined });
  });

  it('treats padded breakpoints as exact-match unsupported values', () => {
    expect(getResponsiveColumnVisibility(['c'], [['c', 'md ']], 0)).toEqual({ c: true });
  });
});
