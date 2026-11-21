import { describe, expect, it } from 'vitest';
import { updateVisibleColumns, getResponsiveColumnVisibility } from '../responsive-column';

describe('updateVisibleColumns', () => {
  it('marks "na" columns as never visible', () => {
    expect(updateVisibleColumns([['hidden', 'na']], 10000)).toEqual({ hidden: false });
  });

  it('marks "xs" columns as always visible', () => {
    expect(updateVisibleColumns([['shown', 'xs']], 0)).toEqual({ shown: true });
    expect(updateVisibleColumns([['shown', 'xs']], 10000)).toEqual({ shown: true });
  });

  it('respects the s breakpoint', () => {
    expect(updateVisibleColumns([['c', 's']], 689)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 's']], 690)).toEqual({ c: true });
  });

  it('respects the m breakpoint', () => {
    expect(updateVisibleColumns([['c', 'm']], 774)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 'm']], 775)).toEqual({ c: true });
  });

  it('respects the l breakpoint', () => {
    expect(updateVisibleColumns([['c', 'l']], 914)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 'l']], 915)).toEqual({ c: true });
  });

  it('respects the xl breakpoint', () => {
    expect(updateVisibleColumns([['c', 'xl']], 1139)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 'xl']], 1140)).toEqual({ c: true });
  });

  it('treats unknown screen sizes as minimum width zero', () => {
    expect(updateVisibleColumns([['c', 'unknown']], 0)).toEqual({ c: true });
  });

  it.each(['foo', 'md '])(
    'treats unsupported breakpoint "%s" as minimum width zero',
    (screenSize) => {
      expect(updateVisibleColumns([['c', screenSize]], 0)).toEqual({ c: true });
    },
  );

  it('preserves fallback visibility for unknown and padded breakpoints', () => {
    expect(
      updateVisibleColumns(
        [
          ['unknown', 'foo'],
          ['padded', 'md '],
          ['paddedKnown', 's '],
        ],
        0,
      ),
    ).toEqual({ unknown: true, padded: true, paddedKnown: true });
  });

  it('handles multiple columns', () => {
    expect(
      updateVisibleColumns(
        [
          ['a', 'xs'],
          ['b', 's'],
          ['c', 'l'],
          ['d', 'na'],
        ],
        800,
      ),
    ).toEqual({ a: true, b: true, c: false, d: false });
  });

  it('returns an empty object for an empty array', () => {
    expect(updateVisibleColumns([], 500)).toEqual({});
  });

  it('defensively handles nullish and non-array colViews', () => {
    expect(updateVisibleColumns(undefined, 800)).toEqual({});
    expect(updateVisibleColumns(null, 800)).toEqual({});
    expect(updateVisibleColumns('invalid' as unknown as unknown[], 800)).toEqual({});
    expect(updateVisibleColumns(123 as unknown as unknown[], 800)).toEqual({});
  });

  it('skips malformed column tuples and handles non-numeric width', () => {
    const malformed = [
      null,
      undefined,
      [],
      ['missingSize'],
      ['name', 'xs', 'extra'],
      [123, 'xs'],
      ['valid', 's'],
    ];
    expect(updateVisibleColumns(malformed as unknown as unknown[], 700)).toEqual({ valid: true });
    expect(updateVisibleColumns([['valid', 's']], NaN)).toEqual({ valid: false });
    expect(updateVisibleColumns([['valid', 's']], undefined)).toEqual({ valid: false });
  });

  it('safely handles non-string screen sizes without TypeError', () => {
    const nonStringSizes = [
      ['col1', { toString: null, valueOf: null }],
      ['col2', 12345],
      ['col3', null],
    ];
    const res = updateVisibleColumns(nonStringSizes as unknown as unknown[], 800);
    expect(res).toEqual({ col1: true, col2: true, col3: true });
  });

  it('ignores blank screen sizes', () => {
    expect(
      updateVisibleColumns(
        [
          ['blank', ''],
          ['whitespace', '   '],
          ['valid', 'xs'],
        ],
        800,
      ),
    ).toEqual({ valid: true });
  });

  it.each([NaN, Infinity, -Infinity, -5, -1, '800', null, undefined])(
    'uses zero as the fallback for invalid or negative width: %s',
    (width) => {
      expect(updateVisibleColumns([['small', 's']], width)).toEqual({ small: false });
      expect(updateVisibleColumns([['always', 'xs']], width)).toEqual({ always: true });
    },
  );
  it('trims padded breakpoints and treats unknown padded breakpoints as zero', () => {
    expect(updateVisibleColumns([['known', ' s ']], 690)).toEqual({ known: true });
    expect(updateVisibleColumns([['unknown', ' md ']], 0)).toEqual({ unknown: true });
  });

  it('handles prototype names safely', () => {
    const result = updateVisibleColumns(
      [
        ['__proto__', 'xs'],
        ['constructor', 'na'],
        ['toString', 's'],
      ],
      800,
    );

    expect(Object.prototype.hasOwnProperty.call(result, '__proto__')).toBe(true);
    expect(result['__proto__']).toBe(true);
    expect(result.constructor).toBe(false);
    expect(result.toString).toBe(true);
  });
});

describe('getResponsiveColumnVisibility', () => {
  it('returns visibility keyed by the supplied column list', () => {
    const colViews = [
      ['name', 'xs'],
      ['size', 'l'],
      ['hidden', 'na'],
    ];
    expect(getResponsiveColumnVisibility(['name', 'size', 'hidden'], colViews, 800)).toEqual({
      name: true,
      size: false,
      hidden: false,
    });
  });

  it('yields undefined for columns missing from colViews', () => {
    const result = getResponsiveColumnVisibility(['name', 'missing'], [['name', 'xs']], 800);
    expect(result.name).toBe(true);
    expect(result.missing).toBeUndefined();
  });

  it('returns an empty object when no column names are supplied', () => {
    expect(getResponsiveColumnVisibility([], [['name', 'xs']], 800)).toEqual({});
  });

  it('defensively handles nullish and non-array columnNames or colViews', () => {
    expect(getResponsiveColumnVisibility(undefined, [['name', 'xs']], 800)).toEqual({});
    expect(getResponsiveColumnVisibility(null, [['name', 'xs']], 800)).toEqual({});
    expect(getResponsiveColumnVisibility(['name'], undefined, 800)).toEqual({});
    expect(getResponsiveColumnVisibility(['name'], null, 800)).toEqual({});
    expect(getResponsiveColumnVisibility(['name'], 'invalid' as any, 800)).toEqual({});
  });

  it('skips non-string column names without throwing TypeError', () => {
    const invalidNames = ['valid', { toString: null, valueOf: null }, 123, null, undefined];
    const result = getResponsiveColumnVisibility(
      invalidNames as unknown as string[],
      [['valid', 'xs']],
      800,
    );
    expect(result).toEqual({ valid: true });
  });

  it.each([null, undefined, 'invalid', {}])(
    'returns an empty object for invalid columnNames: %s',
    (columnNames) => {
      expect(getResponsiveColumnVisibility(columnNames, [['name', 'xs']], 800)).toEqual({});
    },
  );

  it('ignores invalid entries in columnNames', () => {
    expect(
      getResponsiveColumnVisibility(
        ['name', null, 123, '', 'hidden'],
        [
          ['name', 'xs'],
          ['hidden', 'na'],
        ],
        800,
      ),
    ).toEqual({ name: true, hidden: false });
  });

  it('handles invalid colViews safely', () => {
    expect(getResponsiveColumnVisibility(['name', 'size'], null, 800)).toEqual({
      name: undefined,
      size: undefined,
    });
  });

  it('does not expose inherited prototype properties for missing column names', () => {
    expect(getResponsiveColumnVisibility(['constructor', 'toString'], [], 800)).toEqual({
      constructor: undefined,
      toString: undefined,
    });
  });
});
