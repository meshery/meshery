import type { GridPos, MetricPanel } from './types';

export function nextGridPos(panels: Pick<MetricPanel, 'gridPos'>[]): GridPos {
  const width = 12;
  let leftBottom = 0;
  let rightBottom = 0;

  for (const { gridPos } of panels) {
    // Match the placement defaults used by MetricsGrid for saved panels.
    const x = gridPos?.x ?? 0;
    const bottom = (gridPos?.y ?? 0) + (gridPos?.h || 8);
    if (x < width) leftBottom = Math.max(leftBottom, bottom);
    if (x + (gridPos?.w || width) > width) rightBottom = Math.max(rightBottom, bottom);
  }

  return {
    x: leftBottom <= rightBottom ? 0 : width,
    y: Math.min(leftBottom, rightBottom),
    w: width,
    h: 8,
  };
}
