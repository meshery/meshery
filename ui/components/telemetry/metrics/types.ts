// Types for the Prometheus-native telemetry module. These mirror the backend
// wire contract (server/models/telemetry/prometheus). Product-agnostic types
// (GridPos, TimeWindow, ChartSeries) live in ../common/types.
import type { GridPos } from '../common/types';

export type { GridPos } from '../common/types';

// A saved Prometheus panel: a titled PromQL query with a visualization type.
// Persisted per connection in the connection's metadata.
export interface MetricPanel {
  id: string;
  title: string;
  expr: string;
  type: string; // timeseries | stat | gauge | bar
  unit?: string;
  gridPos: GridPos;
}

// Visualization types a user can pick when saving a panel.
export const PANEL_TYPES = ['timeseries', 'stat', 'gauge', 'bar'] as const;

// Metric metadata as returned by Prometheus /api/v1/metadata, normalized to the
// first entry per metric.
export interface MetricMetadata {
  type?: string;
  help?: string;
  unit?: string;
}
