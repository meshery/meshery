import { describe, expect, it, beforeEach, afterEach, vi } from 'vitest';

// `utils.tsx` (where JsonParse lives) drags in next/router via _app.tsx and
// the relay environment. Mock the few symbols k8s-utils actually needs so we
// avoid evaluating the whole UI shell during testing.
vi.mock('../utils', () => ({
  JsonParse: (item: unknown, safe = true) => {
    if (typeof item === 'string') {
      try {
        return JSON.parse(item || '{}');
      } catch (e) {
        if (safe) return {};
        throw e;
      }
    }
    return item;
  },
}));

import {
  CLUSTER_ACTION_GRACE_PERIOD,
  compareUnits,
  formatDuration,
  getPercentStr,
  getReadyReplicas,
  getResourceStr,
  getStatus,
  getTotalReplicas,
  normalizeUnit,
  parseCpu,
  parseDiskSpace,
  parseRam,
  resourceParsers,
  timeAgo,
  TO_GB,
  TO_ONE_CPU,
  TO_ONE_M_CPU,
  unparseCpu,
  unparseRam,
} from '../k8s-utils';

describe('constants', () => {
  it('exposes the cluster-action grace period in ms', () => {
    expect(CLUSTER_ACTION_GRACE_PERIOD).toBe(5000);
  });

  it('exposes the byte/cpu conversion constants', () => {
    expect(TO_GB).toBe(1024 ** 3);
    expect(TO_ONE_M_CPU).toBe(1_000_000);
    expect(TO_ONE_CPU).toBe(1_000_000_000);
  });
});

describe('timeAgo / formatDuration', () => {
  beforeEach(() => {
    delete process.env.UNDER_TEST;
  });

  it('formats a recent duration with the brief format', () => {
    const result = timeAgo(new Date(Date.now() - 60_000), { format: 'brief' });
    // brief uses largest=1 and round=true; expect a unit somewhere
    expect(result).toMatch(/minute|second/);
  });

  it('formats with the mini format', () => {
    const result = timeAgo(new Date(Date.now() - 60_000), { format: 'mini' });
    // mini compresses unit labels - expect single-letter unit suffix
    expect(result).toMatch(/^\d+(s|m)$/);
  });

  it('uses a deterministic reference date when UNDER_TEST=true', () => {
    process.env.UNDER_TEST = 'true';
    const result = timeAgo('2020-01-01T00:00:00Z', { format: 'brief' });
    expect(result).toMatch(/month/i);
    delete process.env.UNDER_TEST;
  });

  it('formatDuration handles brief output for ms input', () => {
    expect(formatDuration(120_000, { format: 'brief' })).toMatch(/minute/);
  });

  it('formatDuration handles mini output for ms input', () => {
    expect(formatDuration(120_000, { format: 'mini' })).toMatch(/^\d+m$/);
  });

  it('formatDuration defaults to brief when format is omitted', () => {
    expect(formatDuration(120_000, {})).toMatch(/minute/);
  });
});

describe('getPercentStr', () => {
  it('returns null when total is 0', () => {
    expect(getPercentStr(5, 0)).toBeNull();
  });

  it('renders whole percentages without decimals', () => {
    expect(getPercentStr(5, 10)).toBe('50 %');
    expect(getPercentStr(10, 10)).toBe('100 %');
  });

  it('renders fractional percentages with one decimal', () => {
    expect(getPercentStr(1, 3)).toMatch(/^33\.3 %$/);
  });
});

describe('getReadyReplicas / getTotalReplicas', () => {
  it('reads ready replicas from status.readyReplicas first', () => {
    expect(getReadyReplicas({ status: { readyReplicas: 3, numberReady: 9 } })).toBe(3);
  });

  it('falls back to status.numberReady', () => {
    expect(getReadyReplicas({ status: { numberReady: 2 } })).toBe(2);
  });

  it('defaults to 0 when neither field is present', () => {
    expect(getReadyReplicas({ status: {} })).toBe(0);
  });

  it('reads total replicas from spec.replicas first', () => {
    expect(getTotalReplicas({ spec: { replicas: 5 }, status: { currentNumberScheduled: 9 } })).toBe(
      5,
    );
  });

  it('falls back to status.currentNumberScheduled', () => {
    expect(getTotalReplicas({ spec: {}, status: { currentNumberScheduled: 4 } })).toBe(4);
  });

  it('defaults to 0 when neither field is present', () => {
    expect(getTotalReplicas({ spec: {}, status: {} })).toBe(0);
  });
});

describe('compareUnits', () => {
  it('compares numeric prefixes regardless of whitespace and case', () => {
    expect(compareUnits('100Mi', '100 MI')).toBe(true);
    expect(compareUnits('1Gi', '1 gi')).toBe(true);
  });

  it('returns false when numeric prefixes differ', () => {
    expect(compareUnits('100Mi', '200Mi')).toBe(false);
  });
});

describe('normalizeUnit', () => {
  it('formats whole CPU as cores', () => {
    expect(normalizeUnit('cpu', '1')).toBe('1 core');
    expect(normalizeUnit('cpu', '2')).toBe('2 cores');
  });

  it('converts CPU milliunits to fractional cores', () => {
    expect(normalizeUnit('cpu', '500m')).toBe('0.5 cores');
  });

  it('strips a prefix on the type before switching', () => {
    expect(normalizeUnit('limits.cpu', '1')).toBe('1 core');
    expect(normalizeUnit('requests.memory', '1024')).toBe('1.02 KB');
  });

  it('handles common memory binary suffixes', () => {
    expect(normalizeUnit('memory', '1Ki')).toBe('1.02 KB');
    expect(normalizeUnit('memory', '1Mi')).toBe('1.05 MB');
    expect(normalizeUnit('memory', '1Gi')).toBe('1.07 GB');
    expect(normalizeUnit('memory', '1Ti')).toBe('1.1 TB');
  });

  it('handles decimal memory suffixes', () => {
    expect(normalizeUnit('memory', '1k')).toBe('1 KB');
    expect(normalizeUnit('memory', '1M')).toBe('1 MB');
    expect(normalizeUnit('memory', '1G')).toBe('1 GB');
    expect(normalizeUnit('memory', '1T')).toBe('1 TB');
  });

  it('handles sub-byte memory suffixes', () => {
    expect(normalizeUnit('memory', '1000m')).toBe('1 Bytes');
  });

  it('returns 0 Bytes for empty/zero memory', () => {
    expect(normalizeUnit('memory', '0')).toBe('0 Bytes');
  });

  it('passes through quantities for unknown resource types', () => {
    expect(normalizeUnit('unknown', '42')).toBe('42');
  });
});

describe('parseDiskSpace / parseRam (parseUnitsOfBytes)', () => {
  it('returns 0 for nullish input', () => {
    expect(parseDiskSpace(undefined as unknown as string)).toBe(0);
    expect(parseRam(undefined as unknown as string)).toBe(0);
  });

  it('parses a bare number', () => {
    expect(parseRam('500')).toBe(500);
  });

  it('parses decimal-suffix bytes', () => {
    expect(parseRam('1K')).toBe(1000);
    expect(parseRam('2M')).toBe(2_000_000);
    expect(parseRam('1G')).toBe(1_000_000_000);
  });

  it('parses binary-suffix bytes', () => {
    expect(parseRam('1Ki')).toBe(1024);
    expect(parseRam('1Mi')).toBe(1024 * 1024);
    expect(parseRam('1Gi')).toBe(1024 ** 3);
  });

  it('parses exponent notation', () => {
    expect(parseRam('5e3')).toBe(5000);
  });
});

describe('unparseRam', () => {
  it('returns Bi for sub-1024 values', () => {
    expect(unparseRam(512)).toEqual({ value: 512, unit: 'Bi' });
  });

  it('walks up through Ki/Mi/Gi/Ti/Pi/Ei', () => {
    expect(unparseRam(1024).unit).toBe('Ki');
    expect(unparseRam(1024 ** 2).unit).toBe('Mi');
    expect(unparseRam(1024 ** 3).unit).toBe('Gi');
    expect(unparseRam(1024 ** 4).unit).toBe('Ti');
    expect(unparseRam(1024 ** 5).unit).toBe('Pi');
    expect(unparseRam(1024 ** 6).unit).toBe('Ei');
  });

  it('rounds the value to one decimal', () => {
    expect(unparseRam(1536).value).toBe(1.5);
  });
});

describe('parseCpu', () => {
  it('returns 0 for nullish input', () => {
    expect(parseCpu('')).toBe(0);
    expect(parseCpu(undefined as unknown as string)).toBe(0);
  });

  it('treats suffix-less values as cpu * 1e9 (n nano-cpu)', () => {
    expect(parseCpu('1')).toBe(1_000_000_000);
  });

  it('parses suffixed CPU values', () => {
    expect(parseCpu('500n')).toBe(500);
    expect(parseCpu('500u')).toBe(500 * 1000);
    expect(parseCpu('500m')).toBe(500 * 1_000_000);
  });
});

describe('unparseCpu', () => {
  it('converts to millicpu rounded to 2 decimals', () => {
    expect(unparseCpu('500000000')).toEqual({ value: 500, unit: 'm' });
    expect(unparseCpu('250000000')).toEqual({ value: 250, unit: 'm' });
  });
});

describe('resourceParsers / getResourceStr', () => {
  it('exposes parsers for cpu and memory', () => {
    expect(resourceParsers.cpu('500m')).toBe(500 * 1_000_000);
    expect(resourceParsers.memory('1Ki')).toBe(1024);
  });

  it('formats cpu via unparseCpu', () => {
    expect(getResourceStr('500000000', 'cpu')).toBe('500m');
  });

  it('formats memory via unparseRam', () => {
    expect(getResourceStr(1024, 'memory')).toBe('1Ki');
  });
});

describe('getStatus', () => {
  it('returns false for non-strings or empty input', () => {
    expect(getStatus(null)).toBe(false);
    expect(getStatus(undefined)).toBe(false);
    expect(getStatus(42 as unknown as string)).toBe(false);
  });

  it('returns false when the status JSON is missing meaningful fields', () => {
    expect(getStatus(JSON.stringify({}))).toBe(false);
  });

  it('returns the phase when present', () => {
    expect(getStatus(JSON.stringify({ phase: 'Running' }))).toBe('Running');
  });

  it('returns Ready when the most recent condition is a True Ready', () => {
    const status = JSON.stringify({
      conditions: [
        { type: 'Available', status: 'True' },
        { type: 'Ready', status: 'True' },
      ],
    });
    expect(getStatus(status)).toBe('Ready');
  });

  it('returns false when most recent condition is not Ready', () => {
    const status = JSON.stringify({
      conditions: [
        { type: 'Ready', status: 'True' },
        { type: 'Available', status: 'False' },
      ],
    });
    expect(getStatus(status)).toBe(false);
  });

  it('returns false when Ready condition status is False', () => {
    const status = JSON.stringify({
      conditions: [{ type: 'Ready', status: 'False' }],
    });
    expect(getStatus(status)).toBe(false);
  });

  it('returns false for malformed JSON (JsonParse returns {} safely)', () => {
    expect(getStatus('not json')).toBe(false);
  });

  it('falls through to false when JSON parses to a non-object truthy value', () => {
    expect(getStatus(JSON.stringify(null))).toBe(false);
    expect(getStatus(JSON.stringify({ conditions: [] }))).toBe(false);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });
});
