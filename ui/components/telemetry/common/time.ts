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

const SHORT_SUFFIXES = ['', 'K', 'M', 'B'];
const IEC_PREFIXES = ['', 'Ki', 'Mi', 'Gi', 'Ti', 'Pi', 'Ei'];
const SI_PREFIXES = ['', 'k', 'M', 'G', 'T', 'P', 'E'];

// Time units in seconds, largest first.
const TIME_UNITS: Array<[number, string]> = [
  [365 * 24 * 60 * 60, 'year'],
  [7 * 24 * 60 * 60, 'week'],
  [24 * 60 * 60, 'day'],
  [60 * 60, 'hour'],
  [60, 'min'],
  [1, 's'],
];

// scale divides value by base until it fits below base, choosing the matching
// entry of suffixes. It compares the rounded value so 999.999 becomes 1.00K
// rather than 1000.00.
function scale(value: number, base: number, suffixes: string[]): [string, string] {
  let scaled = value;
  let i = 0;
  while (i < suffixes.length - 1 && Math.abs(Number(scaled.toFixed(2))) >= base) {
    scaled /= base;
    i++;
  }
  return [scaled.toFixed(2), suffixes[i]];
}

function formatShort(value: number): string {
  const abs = Math.abs(value);
  if (abs < 1 && abs > 0) return value.toFixed(4);
  const [num, suffix] = scale(value, 1000, SHORT_SUFFIXES);
  return `${num}${suffix}`;
}

function formatPrefixed(value: number, base: number, prefixes: string[], unit: string): string {
  const [num, prefix] = scale(value, base, prefixes);
  return `${num} ${prefix}${unit}`;
}

function formatSeconds(seconds: number): string {
  const abs = Math.abs(seconds);
  if (abs > 0 && abs < 1) return `${(seconds * 1000).toFixed(2)} ms`;
  for (const [size, label] of TIME_UNITS) {
    if (Number((abs / size).toFixed(2)) >= 1) return `${(seconds / size).toFixed(2)} ${label}`;
  }
  return `${seconds.toFixed(2)} s`;
}

// formatValue renders a numeric metric value for chart axes, tooltips and stat
// panels, following Grafana's unit ids (fieldConfig.defaults.unit).
export function formatValue(value: number, unit?: string): string {
  if (!isFinite(value)) return '—';

  if (unit?.startsWith('suffix:')) return `${formatShort(value)}${unit.slice('suffix:'.length)}`;
  if (unit?.startsWith('prefix:')) return `${unit.slice('prefix:'.length)}${formatShort(value)}`;

  switch (unit) {
    case undefined:
    case '':
    case 'none':
    case 'short':
      return formatShort(value);
    case 'percent':
      return `${formatShort(value)}%`;
    case 'percentunit':
      return `${(value * 100).toFixed(2)}%`;
    case 'bytes':
      return formatPrefixed(value, 1024, IEC_PREFIXES, 'B');
    case 'decbytes':
      return formatPrefixed(value, 1000, SI_PREFIXES, 'B');
    case 'binBps':
      return formatPrefixed(value, 1024, IEC_PREFIXES, 'B/s');
    case 'Bps':
      return formatPrefixed(value, 1000, SI_PREFIXES, 'B/s');
    case 's':
      return formatSeconds(value);
    case 'ms':
      return formatSeconds(value / 1000);
    case 'reqps':
      return `${formatShort(value)} req/s`;
    case 'ops':
      return `${formatShort(value)} ops/s`;
    default:
      return `${formatShort(value)} ${unit}`;
  }
}
