import { describe, expect, it } from 'vitest';
import {
  linearXAxe,
  logXAxe,
  linearYAxe,
  logYAxe,
  getMetadata,
  makeTitle,
  fortioResultToJsChartData,
  myRound,
  pad,
  formatDate,
  makeChart,
  makeOverlayChartTitle,
  makeOverlayChart,
  makeMultiChart,
} from '../chartjs-formatter';

const makeResult = (overrides: any = {}) => ({
  Labels: 'My Test -_- https://example.com',
  StartTime: '2024-01-15T10:30:00Z',
  RequestedQPS: 100,
  ActualQPS: 99.5,
  NumThreads: 4,
  RequestedDuration: '30s',
  ActualDuration: 30_000_000_000, // 30s in ns
  RetCodes: { 200: 100 },
  DurationHistogram: {
    Count: 100,
    Min: 0.001,
    Avg: 0.005,
    Max: 0.02,
    Percentiles: [
      { Percentile: 50, Value: 0.004 },
      { Percentile: 75, Value: 0.006 },
      { Percentile: 90, Value: 0.01 },
      { Percentile: 99, Value: 0.018 },
      { Percentile: 99.9, Value: 0.02 },
    ],
    Data: [
      { Start: 0.001, End: 0.005, Count: 50, Percent: 50 },
      { Start: 0.005, End: 0.01, Count: 30, Percent: 80 },
      { Start: 0.01, End: 0.02, Count: 20, Percent: 100 },
    ],
  },
  URL: 'https://example.com',
  ...overrides,
});

describe('axis definitions', () => {
  it('linearXAxe has type "linear"', () => {
    expect(linearXAxe.type).toBe('linear');
    expect(linearXAxe.scaleLabel.labelString).toBe('Response time in ms');
  });

  it('logXAxe has type "logarithmic" and a numeric tick callback', () => {
    expect(logXAxe.type).toBe('logarithmic');
    expect(typeof (logXAxe as any).ticks.callback).toBe('function');
    expect((logXAxe as any).ticks.callback(1000)).toBe((1000).toLocaleString());
  });

  it('linearYAxe is linear, begins at zero, and is labelled "Count"', () => {
    expect(linearYAxe.type).toBe('linear');
    expect(linearYAxe.ticks.beginAtZero).toBe(true);
    expect(linearYAxe.scaleLabel.labelString).toBe('Count');
    expect(linearYAxe.id).toBe('H');
  });

  it('logYAxe is logarithmic and stringifies tick values', () => {
    expect(logYAxe.type).toBe('logarithmic');
    expect(typeof (logYAxe as any).ticks.callback).toBe('function');
    expect((logYAxe as any).ticks.callback(10)).toBe('10');
  });
});

describe('myRound', () => {
  it('rounds to 2 decimal places by default', () => {
    expect(myRound(1.2345)).toBe(1.23);
    expect(myRound(1.235)).toBe(1.24);
  });

  it('rounds to the requested number of digits', () => {
    expect(myRound(1.234567, 3)).toBe(1.235);
    expect(myRound(1.234567, 0)).toBe(1);
  });

  it('handles zero and negative numbers', () => {
    expect(myRound(0)).toBe(0);
    expect(myRound(-1.234)).toBe(-1.23);
  });

  it('handles very small numbers', () => {
    expect(myRound(0.0001, 4)).toBe(0.0001);
  });
});

describe('pad', () => {
  it('left-pads single-digit numbers with a leading zero', () => {
    expect(pad(0)).toBe('00');
    expect(pad(5)).toBe('05');
    expect(pad(9)).toBe('09');
  });

  it('returns two-digit numbers unchanged as strings/numbers', () => {
    expect(pad(10)).toBe(10);
    expect(pad(59)).toBe(59);
    expect(pad(100)).toBe(100);
  });
});

describe('formatDate', () => {
  it('formats ISO dates as "YYYY-MM-DD HH:MM:SS"', () => {
    const formatted = formatDate('2024-01-05T03:04:05Z');
    // We don't assert exact local-tz output (test env may vary), only structure
    expect(formatted).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/);
  });

  it('formats a date with two-digit padded segments', () => {
    const formatted = formatDate('2024-12-31T23:59:59Z');
    expect(formatted).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/);
  });
});

describe('getMetadata', () => {
  it('uses Labels when rawdata is not supplied', () => {
    const md = getMetadata(null, makeResult());
    expect(md.title.display.value).toBe('My Test');
    expect(md.url.display.value).toBe('https://example.com');
    expect(md.qps.display.value).toContain('Achieved 99.5');
    expect(md.qps.display.value).toContain('Requested 100');
    expect(md.numberOfConnections.display.value).toBe(4);
  });

  it('uses rawdata when supplied (over Labels)', () => {
    const rawdata = [
      {
        name: 'Raw Title',
        runner_results: { URL: 'https://raw.example.com' },
      },
    ];
    const md = getMetadata(rawdata, makeResult());
    expect(md.title.display.value).toBe('Raw Title');
    expect(md.url.display.value).toBe('https://raw.example.com');
  });

  it('uses schema camelCase rawdata URL when supplied', () => {
    const rawdata = [
      {
        name: 'Raw Title',
        runnerResults: { URL: 'https://raw.example.com' },
      },
    ];
    const md = getMetadata(rawdata, makeResult());
    expect(md.title.display.value).toBe('Raw Title');
    expect(md.url.display.value).toBe('https://raw.example.com');
  });

  it('falls back to "No Title" / "No URL" when both are missing', () => {
    const res = makeResult({ Labels: '' });
    const md = getMetadata(null, res);
    expect(md.title.display.value).toBe('No Title');
    expect(md.url.display.value).toBe('No URL');
  });

  it('reports "No Errors" when all responses are 200', () => {
    const md = getMetadata(null, makeResult());
    expect(md.errors.display.value).toBe('No Errors');
  });

  it('reports a percentage when some requests failed', () => {
    const res = makeResult({
      RetCodes: { 200: 75 },
      DurationHistogram: { ...makeResult().DurationHistogram, Count: 100 },
    });
    const md = getMetadata(null, res);
    expect(md.errors.display.value).toBe('25% errors');
  });

  it('reports "100% errors!" when no successful responses exist', () => {
    const res = makeResult({
      RetCodes: {},
      DurationHistogram: { ...makeResult().DurationHistogram, Count: 100 },
    });
    const md = getMetadata(null, res);
    expect(md.errors.display.value).toBe('100% errors!');
  });

  it('uses SERVING status code as a success indicator (grpc)', () => {
    const res = makeResult({ RetCodes: { SERVING: 100 } });
    const md = getMetadata(null, res);
    expect(md.errors.display.value).toBe('No Errors');
  });

  it('formats percentiles correctly', () => {
    const md = getMetadata(null, makeResult());
    const percentiles = md.percentiles.display.value;
    expect(Array.isArray(percentiles)).toBe(true);
    expect(percentiles).toHaveLength(5);
    expect(percentiles[0].display.key).toBe('p50');
    expect(percentiles[0].display.value).toBe('4 ms');
  });

  it('renders min/max/avg in ms with 3 decimal rounding', () => {
    const md = getMetadata(null, makeResult());
    expect(md.minimum.display.value).toBe('1 ms');
    expect(md.average.display.value).toBe('5 ms');
    expect(md.maximum.display.value).toBe('20 ms');
  });

  it('exposes kubernetes node info when present', () => {
    const res = makeResult({
      kubernetes: {
        server_version: 'v1.28',
        nodes: [
          {
            hostname: 'node-1',
            allocatableCpu: '4',
            allocatableMemory: '16Gi',
            architecture: 'amd64',
            operatingSystem: 'linux',
            kubeletVersion: 'v1.28.0',
            containerRuntimeVersion: 'containerd://1.7',
          },
        ],
      },
    });
    const md = getMetadata(null, res);
    expect(md.kubernetes.display.value[0].display.value).toBe('v1.28');
    const nodes = md.kubernetes.display.value[1].display.value;
    expect(Array.isArray(nodes)).toBe(true);
    expect(nodes[0].display.key).toBe('Node 1');
  });
});

describe('makeTitle', () => {
  it('includes title, URL, and start time when Labels include URL info', () => {
    const result = makeTitle(null, makeResult());
    expect(result).toEqual(expect.arrayContaining([expect.stringMatching(/^Title:/)]));
    expect(result).toEqual(expect.arrayContaining([expect.stringMatching(/^URL:/)]));
    expect(result).toEqual(expect.arrayContaining([expect.stringMatching(/^Start Time:/)]));
  });

  it('falls back to grpc format when URL is missing', () => {
    const res = makeResult({ URL: undefined, Destination: 'grpc.example.com' });
    const result = makeTitle(null, res);
    expect(result.find((s: string) => s.startsWith('Destination:'))).toBeDefined();
  });

  it('omits Title/URL when Labels is empty', () => {
    const res = makeResult({ Labels: '' });
    const result = makeTitle(null, res);
    expect(result.find((s: string) => s.startsWith('Title:'))).toBeUndefined();
    expect(result.find((s: string) => s.startsWith('Destination:'))).toBeUndefined();
  });

  it('includes No Error when all responses succeeded', () => {
    const result = makeTitle(null, makeResult());
    expect(result).toContain('Errors: No Error');
  });

  it('reports percent errors when only some requests succeeded', () => {
    const res = makeResult({
      RetCodes: { 200: 50 },
      DurationHistogram: { ...makeResult().DurationHistogram, Count: 100 },
    });
    const result = makeTitle(null, res);
    expect(result).toContain('Errors: 50% errors');
  });

  it('reports 100% errors when no successful responses', () => {
    const res = makeResult({
      RetCodes: null,
      DurationHistogram: { ...makeResult().DurationHistogram, Count: 100 },
    });
    const result = makeTitle(null, res);
    expect(result).toContain('Errors: 100% errors!');
  });

  it('uses rawdata title and URL when supplied', () => {
    const rawdata = [{ name: 'RawName', runner_results: { URL: 'https://raw' } }];
    const result = makeTitle(rawdata, makeResult());
    expect(result).toContain('Title: RawName');
    expect(result).toContain('URL: https://raw');
  });

  it('uses schema camelCase rawdata URL when supplied', () => {
    const rawdata = [{ name: 'RawName', runnerResults: { URL: 'https://raw' } }];
    const result = makeTitle(rawdata, makeResult());
    expect(result).toContain('Title: RawName');
    expect(result).toContain('URL: https://raw');
  });

  it('appends kubernetes node info when present', () => {
    const res = makeResult({
      kubernetes: {
        server_version: 'v1.28',
        nodes: [
          {
            hostname: 'node-1',
            allocatable_cpu: '4',
            allocatable_memory: '16Gi',
            architecture: 'amd64',
            os_image: 'linux',
            kubelet_version: 'v1.28.0',
            container_runtime_version: 'containerd://1.7',
          },
        ],
      },
    });
    const result = makeTitle(null, res);
    expect(result.some((s: string) => s.includes('Kubernetes server version'))).toBe(true);
    expect(result.some((s: string) => s.includes('Node 1'))).toBe(true);
  });
});

describe('fortioResultToJsChartData', () => {
  it('returns an object with title, metadata, dataP, dataH, and percentiles', () => {
    const out = fortioResultToJsChartData(null, makeResult());
    expect(out).toHaveProperty('title');
    expect(out).toHaveProperty('metadata');
    expect(out).toHaveProperty('dataP');
    expect(out).toHaveProperty('dataH');
    expect(out).toHaveProperty('percentiles');
    expect(Array.isArray(out.dataP)).toBe(true);
    expect(Array.isArray(out.dataH)).toBe(true);
  });

  it('dataP first point is at the origin', () => {
    const out = fortioResultToJsChartData(null, makeResult());
    expect(out.dataP[0]).toEqual({ x: 0, y: 0 });
  });

  it('dataH contains all the histogram counts', () => {
    const out = fortioResultToJsChartData(null, makeResult());
    // 3 bins ⇒ 3 pairs of (start, end) points = at least 6 entries
    expect(out.dataH.length).toBeGreaterThanOrEqual(6);
    const ys = out.dataH.map((d: any) => d.y);
    expect(ys).toContain(50);
    expect(ys).toContain(30);
    expect(ys).toContain(20);
  });
});

describe('makeChart', () => {
  it('produces a chart config with two datasets', () => {
    const data = fortioResultToJsChartData(null, makeResult());
    const chart = makeChart(data);
    expect(chart.data.datasets).toHaveLength(2);
    expect(chart.data.datasets[0].label).toBe('Cumulative %');
    expect(chart.data.datasets[1].label).toBe('Histogram: Count');
  });

  it('passes percentiles through', () => {
    const data = fortioResultToJsChartData(null, makeResult());
    const chart = makeChart(data);
    expect(chart.percentiles).toBe(data.percentiles);
  });

  it('sets the title text from the data', () => {
    const data = fortioResultToJsChartData(null, makeResult());
    const chart = makeChart(data);
    expect(chart.options.title.text).toBe(data.title);
  });

  it('sets responsive and maintainAspectRatio options', () => {
    const data = fortioResultToJsChartData(null, makeResult());
    const chart = makeChart(data);
    expect(chart.options.responsive).toBe(true);
    expect(chart.options.maintainAspectRatio).toBe(false);
  });
});

describe('makeOverlayChartTitle', () => {
  it('combines two titles with "A:" and "B:" prefixes', () => {
    const result = makeOverlayChartTitle(['Alpha', 'second-line', 'third'], ['Bravo', 'b2', 'b3']);
    expect(result).toEqual(['A: Alpha', 'second-line', '', 'B: Bravo', 'b2']);
  });
});

describe('makeOverlayChart', () => {
  it('produces a chart config with four datasets', () => {
    const a = fortioResultToJsChartData(null, makeResult());
    const b = fortioResultToJsChartData(null, makeResult());
    const chart = makeOverlayChart(a, b);
    expect(chart.data.datasets).toHaveLength(4);
    expect(chart.data.datasets[0].label).toBe('A: Cumulative %');
    expect(chart.data.datasets[1].label).toBe('B: Cumulative %');
    expect(chart.data.datasets[2].label).toBe('A: Histogram: Count');
    expect(chart.data.datasets[3].label).toBe('B: Histogram: Count');
  });

  it('includes the combined title from both inputs', () => {
    const a = fortioResultToJsChartData(null, makeResult());
    const b = fortioResultToJsChartData(null, makeResult());
    const chart = makeOverlayChart(a, b);
    expect(chart.options.title.text[0]).toMatch(/^A: /);
    expect(chart.options.title.text[3]).toMatch(/^B: /);
  });
});

describe('makeMultiChart', () => {
  it('produces a chart with the expected dataset labels', () => {
    const result = makeMultiChart(null, [makeResult(), makeResult()]);
    const labels = result.data.datasets.map((d: any) => d.label);
    expect(labels).toEqual(['Min', 'Median', 'Avg', 'p75', 'p90', 'p99', 'p99.9', 'Max', 'QPS']);
  });

  it('truncates dataset arrays to the number of results', () => {
    const results = [makeResult(), makeResult(), makeResult()];
    const chart = makeMultiChart(null, results);
    chart.data.datasets.forEach((d: any) => {
      expect(d.data.length).toBe(3);
    });
    expect(chart.data.labels.length).toBe(3);
  });

  it('uses ms-converted values for min/max/avg in series 0/2/7', () => {
    const chart = makeMultiChart(null, [makeResult()]);
    expect(chart.data.datasets[0].data[0]).toBeCloseTo(1, 5); // min 0.001s -> 1ms
    expect(chart.data.datasets[2].data[0]).toBeCloseTo(5, 5); // avg 0.005s -> 5ms
    expect(chart.data.datasets[7].data[0]).toBeCloseTo(20, 5); // max 0.02s -> 20ms
  });

  it('handles results without percentiles gracefully', () => {
    const res = makeResult();
    res.DurationHistogram.Percentiles = undefined;
    expect(() => makeMultiChart(null, [res])).not.toThrow();
  });

  it('handles labels without the " -_- " separator', () => {
    const res = makeResult({ Labels: 'simple-label' });
    const chart = makeMultiChart(null, [res]);
    // labels should still exist
    expect(chart.data.labels.length).toBe(1);
  });

  it('uses the QPS yAxisID for the QPS dataset', () => {
    const chart = makeMultiChart(null, [makeResult()]);
    const qps = chart.data.datasets.find((d: any) => d.label === 'QPS');
    expect(qps.yAxisID).toBe('qps');
  });
});
