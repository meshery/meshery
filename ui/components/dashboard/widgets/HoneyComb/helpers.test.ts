import { describe, expect, it } from 'vitest';
import { getColumnSize, getGridColumnsCount, getRowSize } from './helpers';

describe('getGridColumnsCount', () => {
  it('returns 1 when either argument is non-positive', () => {
    expect(getGridColumnsCount(0, 1000)).toBe(1);
    expect(getGridColumnsCount(10, 0)).toBe(1);
    expect(getGridColumnsCount(-5, 100)).toBe(1);
    expect(getGridColumnsCount(10, -100)).toBe(1);
  });

  it('returns at least 1 column even for narrow containers', () => {
    expect(getGridColumnsCount(50, 10)).toBe(1);
  });

  it('computes columns based on hexagon width (sqrt(3) * side)', () => {
    // For side=10 and width=200 -> floor(200/(sqrt(3)*10)) ~ 11
    const result = getGridColumnsCount(10, 200);
    expect(result).toBe(Math.floor(200 / (Math.sqrt(3) * 10)));
    expect(result).toBeGreaterThan(0);
  });
});

describe('getRowSize', () => {
  it('returns side/2 for positive sides above the floor', () => {
    expect(getRowSize(20)).toBe(10);
    expect(getRowSize(4)).toBe(2);
  });

  it('floors row size at 1 for small/zero/negative sides', () => {
    // Math.max(1, hexagonSide / 2)
    expect(getRowSize(1)).toBe(1);
    expect(getRowSize(0)).toBe(1);
    expect(getRowSize(-10)).toBe(1);
  });
});

describe('getColumnSize', () => {
  it('returns sqrt(3) * side / 4 for positive sides', () => {
    expect(getColumnSize(8)).toBeCloseTo((Math.sqrt(3) * 8) / 4);
  });

  it('returns at least 1 to avoid zero-width columns', () => {
    expect(getColumnSize(0)).toBe(1);
    expect(getColumnSize(-50)).toBe(1);
  });
});
