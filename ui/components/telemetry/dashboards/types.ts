// Shared types for the telemetry dashboards UI. These mirror the wire contract
// served by the backend grafana telemetry handlers
// (server/models/telemetry/grafana).

export interface GridPos {
  x: number;
  y: number;
  w: number;
  h: number;
}

export interface Target {
  refId: string;
  datasourceUid?: string;
  expr: string;
  legendFormat?: string;
}

export interface Panel {
  id: number;
  title: string;
  type: string;
  gridPos: GridPos;
  unit?: string;
  targets: Target[];
}

export interface TemplateVar {
  name: string;
  label?: string;
  query?: string;
  type?: string;
  current?: string[];
  multi?: boolean;
  includeAll?: boolean;
}

export interface Datasource {
  uid: string;
  name: string;
  type: string;
  isDefault?: boolean;
}

export interface Board {
  uid: string;
  title: string;
  tags?: string[];
  panels: Panel[];
  templateVars?: TemplateVar[];
}

export interface BoardSummary {
  uid: string;
  title: string;
  url?: string;
  tags?: string[];
}

export interface PinnedBoard {
  uid: string;
  title: string;
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
