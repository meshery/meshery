import { describe, expect, it } from 'vitest';
import { nextGridPos, normalizeGridPos } from './gridPosition';

describe('nextGridPos', () => {
  const slot = (x: number, y: number, w = 12, h = 8) => ({ gridPos: { x, y, w, h } });

  it('places new panels two per row', () => {
    expect(nextGridPos([])).toEqual(slot(0, 0).gridPos);
    expect(nextGridPos([slot(0, 0)])).toEqual(slot(12, 0).gridPos);
    expect(nextGridPos([slot(0, 0), slot(12, 0)])).toEqual(slot(0, 8).gridPos);
  });

  it('keeps both surviving panels visible after removing the second of three', () => {
    expect(nextGridPos([slot(0, 0), slot(0, 8)])).toEqual(slot(12, 0).gridPos);
  });

  it('clears an existing panel that spans both halves', () => {
    expect(nextGridPos([slot(6, 4, 12, 12)])).toEqual(slot(0, 16).gridPos);
  });

  it('uses the grid origin and default spans for negative coordinates and spans', () => {
    const negativeColumn = slot(-14, 0);
    expect(normalizeGridPos(negativeColumn.gridPos)).toEqual(slot(0, 0).gridPos);
    expect(nextGridPos([negativeColumn])).toEqual(slot(12, 0).gridPos);

    const saved = slot(-14, -1, -12, -8);
    expect(normalizeGridPos(saved.gridPos)).toEqual(slot(0, 0).gridPos);
    expect(nextGridPos([saved, slot(12, 0)])).toEqual(slot(0, 8).gridPos);
  });

  it('leaves the default slots free when a saved panel starts in an implicit column', () => {
    expect(nextGridPos([slot(24, 0)])).toEqual(slot(0, 0).gridPos);
    expect(nextGridPos([slot(24, 0), slot(0, 0)])).toEqual(slot(12, 0).gridPos);
  });
});
