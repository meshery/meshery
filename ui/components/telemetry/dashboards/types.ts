// Grafana-specific telemetry types served by the backend grafana telemetry
// handlers (server/models/telemetry/grafana). Product-agnostic types
// (GridPos, TimeWindow, ChartSeries) live in ../common/types.
import type { GridPos } from '../common/types';

export type { GridPos, TimeWindow, ChartSeries } from '../common/types';

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
