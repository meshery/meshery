import { describe, expect, it } from 'vitest';
import { nextGridPos } from './gridPosition';

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
});
