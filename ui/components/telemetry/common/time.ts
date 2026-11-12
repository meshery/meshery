import type { ChartSeries, TimeWindow } from './types';

export interface RangePreset {
  label: string;
  // window duration in seconds
  durationSec: number;
}

// Selectable relative time ranges, all ending at "now".
export const RANGE_PRESETS: RangePreset[] = [
  { label: 'Last 5 minutes', durationSec: 5 * 60 },
  { label: 'Last 15 minutes', durationSec: 15 * 60 },
  { label: 'Last 30 minutes', durationSec: 30 * 60 },
  { label: 'Last 1 hour', durationSec: 60 * 60 },
  { label: 'Last 3 hours', durationSec: 3 * 60 * 60 },
  { label: 'Last 6 hours', durationSec: 6 * 60 * 60 },
  { label: 'Last 12 hours', durationSec: 12 * 60 * 60 },
  { label: 'Last 24 hours', durationSec: 24 * 60 * 60 },
  { label: 'Last 7 days', durationSec: 7 * 24 * 60 * 60 },
];

export const DEFAULT_RANGE = RANGE_PRESETS[3]; // 1 hour

// Auto-refresh interval options, in milliseconds (0 = off).
export const REFRESH_INTERVALS: Array<{ label: string; ms: number }> = [
  { label: 'Off', ms: 0 },
  { label: '5s', ms: 5_000 },
  { label: '10s', ms: 10_000 },
  { label: '30s', ms: 30_000 },
  { label: '1m', ms: 60_000 },
  { label: '5m', ms: 300_000 },
];

// Target number of points across a panel's width; step is derived from this so
// charts stay readable and queries stay cheap regardless of range.
const TARGET_POINTS = 300;
const MIN_STEP_SEC = 15;

// resolveWindow turns a relative duration into an absolute [start, end] window
// ending now, with a sensible step.
export function resolveWindow(durationSec: number, nowMs = Date.now()): TimeWindow {
  const end = Math.floor(nowMs / 1000);
  const start = end - durationSec;
  const step = Math.max(MIN_STEP_SEC, Math.round(durationSec / TARGET_POINTS));
  return { start, end, step };
}

// renderLegend resolves a Grafana legendFormat template (e.g. "{{pod}}")
// against a Prometheus series' metric labels, falling back to a compact
// label-set representation.
export function renderLegend(
  legendFormat: string | undefined,
  metric: Record<string, string>,
): string {
  // Grafana uses the sentinel "__auto" to mean "derive the series name from the
  // metric labels" — treat it the same as an unset legendFormat.
  if (legendFormat && legendFormat.trim() !== '' && legendFormat.trim() !== '__auto') {
    return legendFormat.replace(/\{\{\s*([\w.]+)\s*\}\}/g, (_m, key) => metric[key] ?? '');
  }
  const name = metric.__name__ ?? '';
  const labels = Object.entries(metric)
    .filter(([k]) => k !== '__name__')
    .map(([k, v]) => `${k}="${v}"`)
    .join(', ');
  if (name && labels) return `${name}{${labels}}`;
  return name || (labels ? `{${labels}}` : 'value');
}

// parsePromMatrix converts a Prometheus query_range response into chart series.
// It tolerates both matrix (range) and vector (instant) result types.
export function parsePromMatrix(resp: any, legendFormat?: string): ChartSeries[] {
  const result = resp?.data?.result;
  if (!Array.isArray(result)) return [];
  const resultType = resp?.data?.resultType;

  return result.map((s: any) => {
    const metric: Record<string, string> = s.metric ?? {};
    let points: Array<[number, number]> = [];
    if (resultType === 'vector' && Array.isArray(s.value)) {
      points = [[Number(s.value[0]), Number(s.value[1])]];
    } else if (Array.isArray(s.values)) {
      points = s.values.map((v: [number, string]) => [Number(v[0]), Number(v[1])]);
    }
    return { name: renderLegend(legendFormat, metric), points };
  });
}

// formatValue renders a numeric metric value with light unit awareness for
// stat panels.
export function formatValue(value: number, unit?: string): string {
  if (!isFinite(value)) return '—';
  const abs = Math.abs(value);
  let formatted: string;
  if (abs >= 1_000_000_000) formatted = `${(value / 1_000_000_000).toFixed(2)}B`;
  else if (abs >= 1_000_000) formatted = `${(value / 1_000_000).toFixed(2)}M`;
  else if (abs >= 1_000) formatted = `${(value / 1_000).toFixed(2)}K`;
  else if (abs < 1 && abs > 0) formatted = value.toFixed(4);
  else formatted = value.toFixed(2);

  switch (unit) {
    case 'percent':
      return `${formatted}%`;
    case 'percentunit':
      return `${(value * 100).toFixed(2)}%`;
    case 'bytes':
      return `${formatted} B`;
    case 's':
      return `${formatted} s`;
    default:
      return formatted;
  }
}
