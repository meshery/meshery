import { describe, expect, it } from 'vitest';
import { dataToColors, getChartColors, isValidColumnName } from '../charts';

// Minimal theme stub matching the shape `darken` and `lighten` expect.
// `@mui/material`'s `darken` / `lighten` accept hex/rgb strings, so any valid
// palette entry will produce real string outputs at runtime.
const makeTheme = () =>
  ({
    palette: {
      primary: { main: '#1976d2' },
      info: { main: '#0288d1' },
    },
  }) as never;

describe('isValidColumnName', () => {
  it('returns true for non-empty strings', () => {
    expect(isValidColumnName('name')).toBe(true);
    expect(isValidColumnName('hello world')).toBe(true);
  });

  it('returns false for empty string and single-space', () => {
    expect(isValidColumnName('')).toBe(false);
    expect(isValidColumnName(' ')).toBe(false);
  });

  it('returns false for undefined and null', () => {
    expect(isValidColumnName(undefined)).toBe(false);
    expect(isValidColumnName(null)).toBe(false);
  });

  it('treats arbitrary whitespace beyond a single space as valid', () => {
    // The implementation only checks against '' and ' '. Multi-character
    // whitespace strings are considered valid by current rules.
    expect(isValidColumnName('  ')).toBe(true);
  });
});

describe('getChartColors', () => {
  it('returns an 8-color palette derived from theme primary/info', () => {
    const colors = getChartColors(makeTheme());
    expect(Array.isArray(colors)).toBe(true);
    expect(colors).toHaveLength(8);
    colors.forEach((c) => {
      expect(typeof c).toBe('string');
      expect(c.length).toBeGreaterThan(0);
    });
  });

  it('the middle slot (index 3) is exactly the primary color', () => {
    const theme = makeTheme();
    const colors = getChartColors(theme);
    expect(colors[3]).toBe(theme.palette.primary.main);
  });
});

describe('dataToColors', () => {
  it('maps each row label to a palette color', () => {
    const theme = makeTheme();
    const data = [
      ['a', 1],
      ['b', 2],
      ['c', 3],
    ];
    const mapping = dataToColors(data, theme);
    expect(Object.keys(mapping)).toEqual(['a', 'b', 'c']);
    const palette = getChartColors(theme);
    expect(mapping.a).toBe(palette[0]);
    expect(mapping.b).toBe(palette[1]);
    expect(mapping.c).toBe(palette[2]);
  });

  it('wraps palette colors when there are more rows than palette entries', () => {
    const theme = makeTheme();
    const palette = getChartColors(theme);
    const data = Array.from({ length: palette.length + 2 }, (_, i) => [`row-${i}`, i]);
    const mapping = dataToColors(data, theme);

    expect(mapping[`row-${palette.length}`]).toBe(palette[0]);
    expect(mapping[`row-${palette.length + 1}`]).toBe(palette[1]);
  });

  it('returns an empty object for empty data', () => {
    expect(dataToColors([], makeTheme())).toEqual({});
  });

  it('handles duplicate column labels by overwriting earlier entries', () => {
    const theme = makeTheme();
    const mapping = dataToColors(
      [
        ['a', 1],
        ['a', 2],
      ],
      theme,
    );
    const palette = getChartColors(theme);
    // The forEach over `columns` will set `a` twice -- the latest assignment wins
    expect(mapping.a).toBe(palette[1]);
  });
});
