// Shared, product-agnostic telemetry types used by both the Grafana (dashboards)
// and Prometheus (metrics) telemetry modules.

// Grafana's 24-column panel layout coordinates.
export interface GridPos {
  x: number;
  y: number;
  w: number;
  h: number;
}

// A resolved time window for queries: an absolute [start, end] in epoch seconds
// plus a step (resolution) in seconds.
export interface TimeWindow {
  start: number;
  end: number;
  step: number;
}

// A single rendered series for a chart.
export interface ChartSeries {
  name: string;
  // points are [epochSeconds, value]; value may be NaN for gaps.
  points: Array<[number, number]>;
}

// The minimal panel shape the presentational Panel component renders. Both the
// Grafana Panel and the Prometheus MetricPanel are structurally compatible.
export interface PanelMeta {
  title?: string;
  type: string;
  unit?: string;
}
