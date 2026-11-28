import type { GridPos, MetricPanel } from './types';

export const GRID_COLUMNS = 24;
const DEFAULT_WIDTH = GRID_COLUMNS / 2;
const DEFAULT_HEIGHT = 8;

export function normalizeGridPos(gridPos?: GridPos): GridPos {
  return {
    x: Math.max(0, gridPos?.x ?? 0),
    y: Math.max(0, gridPos?.y ?? 0),
    w: gridPos?.w && gridPos.w > 0 ? gridPos.w : DEFAULT_WIDTH,
    h: gridPos?.h && gridPos.h > 0 ? gridPos.h : DEFAULT_HEIGHT,
  };
}

export function nextGridPos(panels: Pick<MetricPanel, 'gridPos'>[]): GridPos {
  let leftBottom = 0;
  let rightBottom = 0;

  for (const { gridPos } of panels) {
    const { x, y, w, h } = normalizeGridPos(gridPos);
    const bottom = y + h;
    if (x < DEFAULT_WIDTH) leftBottom = Math.max(leftBottom, bottom);
    if (x < GRID_COLUMNS && x + w > DEFAULT_WIDTH) {
      rightBottom = Math.max(rightBottom, bottom);
    }
  }

  return {
    x: leftBottom <= rightBottom ? 0 : DEFAULT_WIDTH,
    y: Math.min(leftBottom, rightBottom),
    w: DEFAULT_WIDTH,
    h: DEFAULT_HEIGHT,
  };
}
