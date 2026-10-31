import { describe, expect, it } from 'vitest';
import { formatResultRow, getRunnerResults } from './MesheryResults';

// MesheryResults owns a paginated mui-datatables view of historical
// load-test results plus a chart drawer that re-renders MesheryChart and
// GrafanaCustomCharts on row selection. It pulls from RTK Query, a redux
// slice for the active context, and a custom router-driven query string.
describe.skip('MesheryResults', () => {
  it('skipped - mui-datatables + RTK Query + chart panel container', () => {});
});

const runnerResults = {
  StartTime: '2026-10-01T10:00:00Z',
  ActualQPS: 9.87,
  ActualDuration: 30000000000,
  NumThreads: 2,
  DurationHistogram: {
    Percentiles: [
      { Percentile: 50, Value: 0.0123 },
      { Percentile: 99.9, Value: 0.5 },
    ],
  },
};

describe('formatResultRow', () => {
  it('formats a camelCase result as sent by the server', () => {
    expect(formatResultRow({ name: 'run-1', mesh: 'istio', runnerResults })).toEqual({
      name: 'run-1',
      mesh: 'istio',
      test_start_time: '2026-10-01T10:00:00Z',
      qps: '9.9',
      duration: '30.0',
      threads: 2,
      p50: '0.012',
      p99_9: '0.500',
    });
  });

  it('still formats a snake_case result', () => {
    const row = formatResultRow({ name: 'run-1', runner_results: runnerResults });
    expect(row.qps).toBe('9.9');
  });

  it('defaults percentiles to 0 when the histogram has none', () => {
    const row = formatResultRow({
      name: 'run-1',
      runnerResults: { ...runnerResults, DurationHistogram: {} },
    });
    expect([row.p50, row.p75, row.p90, row.p99, row.p99_9]).toEqual([0, 0, 0, 0, 0]);
  });
});

describe('getRunnerResults', () => {
  it('prefers runnerResults over runner_results', () => {
    expect(getRunnerResults({ runnerResults: 'camel', runner_results: 'snake' })).toBe('camel');
    expect(getRunnerResults({ runner_results: 'snake' })).toBe('snake');
    expect(getRunnerResults(undefined)).toBeUndefined();
  });
});
