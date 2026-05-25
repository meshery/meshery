import { describe, expect, it } from 'vitest';
import { updateVisibleColumns, getResponsiveColumnVisibility } from '../responsive-column';

describe('updateVisibleColumns', () => {
  it('marks "na" columns as never visible', () => {
    const cols = [['hidden', 'na']];
    expect(updateVisibleColumns(cols, 10000)).toEqual({ hidden: false });
  });

  it('marks "xs" columns as always visible (even when width=0)', () => {
    const cols = [['shown', 'xs']];
    expect(updateVisibleColumns(cols, 0)).toEqual({ shown: true });
    expect(updateVisibleColumns(cols, 10000)).toEqual({ shown: true });
  });

  it('shows "s" columns only above 690px', () => {
    expect(updateVisibleColumns([['c', 's']], 689)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 's']], 690)).toEqual({ c: true });
    expect(updateVisibleColumns([['c', 's']], 800)).toEqual({ c: true });
  });

  it('shows "m" columns only above 775px', () => {
    expect(updateVisibleColumns([['c', 'm']], 774)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 'm']], 775)).toEqual({ c: true });
  });

  it('shows "l" columns only above 915px', () => {
    expect(updateVisibleColumns([['c', 'l']], 914)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 'l']], 915)).toEqual({ c: true });
  });

  it('shows "xl" columns only above 1140px', () => {
    expect(updateVisibleColumns([['c', 'xl']], 1139)).toEqual({ c: false });
    expect(updateVisibleColumns([['c', 'xl']], 1140)).toEqual({ c: true });
  });

  it('treats unknown screen size as min-width 0 (always visible)', () => {
    expect(updateVisibleColumns([['c', 'unknown']], 0)).toEqual({ c: true });
  });

  it('handles multiple columns in one call', () => {
    const cols = [
      ['a', 'xs'],
      ['b', 's'],
      ['c', 'l'],
      ['d', 'na'],
    ];
    expect(updateVisibleColumns(cols, 800)).toEqual({
      a: true,
      b: true,
      c: false,
      d: false,
    });
  });

  it('returns an empty object when given no columns', () => {
    expect(updateVisibleColumns([], 500)).toEqual({});
  });
});

describe('getResponsiveColumnVisibility', () => {
  it('returns visibility keyed by the supplied column list', () => {
    const colViews = [
      ['name', 'xs'],
      ['size', 'l'],
      ['hidden', 'na'],
    ];
    const result = getResponsiveColumnVisibility(['name', 'size', 'hidden'], colViews, 800);
    expect(result).toEqual({ name: true, size: false, hidden: false });
  });

  it('yields undefined for columns missing from colViews', () => {
    const colViews = [['name', 'xs']];
    const result = getResponsiveColumnVisibility(['name', 'missing'], colViews, 800);
    expect(result.name).toBe(true);
    expect(result.missing).toBeUndefined();
  });

  it('returns an empty object when no column names are supplied', () => {
    expect(getResponsiveColumnVisibility([], [['name', 'xs']], 800)).toEqual({});
  });
});
